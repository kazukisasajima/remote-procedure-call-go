package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"remote_procedure_call/protocol"
)

const socketPath = "/tmp/socket_file"
const basePath = "json/"

func StartClient() {
	fmt.Printf("Connecting to %s\n", socketPath)

	// サーバーへの接続を確立
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// ユーザー入力を受け取るためのスキャナー
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> Enter JSON filename: ")
		scanner.Scan()
		input := scanner.Text()

		// "exit"と入力された場合、クライアントを終了
		if input == "exit" {
			fmt.Println("Exiting client.")
			break
		}

		// 入力されたファイル名からリクエストを生成
		filepath := basePath + input + ".json"
		req, err := readFile(filepath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			continue
		}

		// リクエストをサーバーに送信
		if err := sendRequest(conn, req); err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}

		// サーバーからのレスポンスを読み取る
		if err := receiveResponse(conn); err != nil {
			fmt.Printf("Error receiving response: %v\n", err)
			continue
		}
	}
	fmt.Println("Closing socket")
}

// 指定されたJSONファイルを読み取り、リクエスト構造体に変換
func readFile(filepath string) (*protocol.Request, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var req protocol.Request
	if err := json.NewDecoder(file).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %v", err)
	}
	return &req, nil
}

// リクエストをサーバーに送信
func sendRequest(conn net.Conn, req *protocol.Request) error {
	// 構造体をJSON形式のバイト配列にエンコード
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to serialize request: %v", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}
	fmt.Println("Request sent to server.")
	return nil
}

// サーバーからのレスポンスを読み取り表示
func receiveResponse(conn net.Conn) error {
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	n, err := conn.Read(buffer)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("timeout while waiting for response")
		}
		return fmt.Errorf("error reading from server: %v", err)
	}

	var res protocol.Response
	// バイト配列（JSON文字列）をRequest構造体に変換
	if err := json.Unmarshal(buffer[:n], &res); err != nil {
		return fmt.Errorf("invalid response format: %v", err)
	}

	if res.Error != "" {
		fmt.Printf("Server responded with error: %s\n", res.Error)
	} else {
		fmt.Printf("Server response: %v (Type: %s)\n", res.Results, res.ResultType)
	}
	return nil
}
