package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"remote_procedure_call/handler"
    "remote_procedure_call/protocol"
)

const socketPath = "/tmp/socket_file"

func StartServer() {
	if err := cleanUpSocketFile(); err != nil {
		logWithTimestamp(fmt.Sprintf("Error cleaning up socket file: %v", err))
		return
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		logWithTimestamp(fmt.Sprintf("Error creating UNIX socket: %v", err))
		return
	}
	defer listener.Close()

	logWithTimestamp(fmt.Sprintf("Server started at %s", socketPath))

	for {
		conn, err := listener.Accept()
		if err != nil {
			logWithTimestamp(fmt.Sprintf("Error accepting connection: %v", err))
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	logWithTimestamp("Connection established")

	for {
		req, err := readRequest(conn)
		if err != nil {
			if err == io.EOF {
				logWithTimestamp("Client closed the connection.")
			} else {
				logWithTimestamp(fmt.Sprintf("Error reading request: %v", err))
			}
			break
		}

		response := processRequest(req)
		if err := sendResponse(conn, response); err != nil {
			logWithTimestamp(fmt.Sprintf("Error sending response: %v", err))
			break
		}
	}
}

func cleanUpSocketFile() error {
	if _, err := os.Stat(socketPath); err == nil {
		return os.Remove(socketPath)
	}
	return nil
}

// クライアントからのリクエストを読み取り、構造体に変換
func readRequest(conn net.Conn) (*protocol.Request, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	var req protocol.Request
	// バイト配列（JSON文字列）をRequest構造体に変換
	if err := json.Unmarshal(buffer[:n], &req); err != nil {
		return nil, fmt.Errorf("invalid request format: %v", err)
	}

	logWithTimestamp(fmt.Sprintf("Received request: %+v", req))
	return &req, nil // 構造体ポインタを返す
}

// クライアントリクエストのメソッドを実行し、結果を生成
func processRequest(req *protocol.Request) *protocol.Response {
	// handlerにメソッド名とパラメータを渡して処理を依頼
	results, resultType, err := handler.ExecuteRPCMethod(req.Method, req.Params)

	// エラーがあればエラー情報をレスポンスとして返す
	if err != nil {
		return &protocol.Response{
			ID:    req.ID,
			Error: fmt.Sprintf("Error executing method: %v", err),
		}
	}

	// 成功時のレスポンスを返す
	return &protocol.Response{
		Results:    results,
		ResultType: resultType,
		ID:         req.ID,
	}
}

func sendResponse(conn net.Conn, response *protocol.Response) error {
	// 構造体をJSON形式のバイト配列にエンコード
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to serialize response: %v", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to send response: %v", err)
	}

	logWithTimestamp(fmt.Sprintf("Sent response: %+v", response))
	return nil
}

func logWithTimestamp(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}
