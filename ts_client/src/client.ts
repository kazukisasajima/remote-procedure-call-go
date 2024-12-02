import * as net from 'net';
import * as fs from 'fs';
import * as readline from 'readline';

// サーバーのソケットアドレス
const serverAddress = '/tmp/socket_file';
// JSONファイルのベースディレクトリ
const basePath = '../json/';

// リクエストとレスポンスの型定義
interface Request {
    method: string;
    params: (string | number)[];
    id: number;
}

interface Response {
    results?: string | number | boolean | string[];
    result_type?: string;
    id: number;
    error?: string;
}

// JSONファイルからリクエストを読み取る関数
function readFile(filepath: string): Promise<Request> {
    return new Promise((resolve, reject) => {
        fs.readFile(filepath, 'utf8', (err, data) => {
            if (err) {
                reject(`Error reading file: ${err.message}`);
            } else {
                try {
                    const req: Request = JSON.parse(data);
                    resolve(req);
                } catch (parseError) {
                    reject(`Invalid JSON format: ${(parseError as any).message}`);
                }
            }
        });
    });
}

// サーバーからのレスポンスを処理する関数
function handleServerResponse(data: Buffer): void {
    try {
        const response: Response = JSON.parse(data.toString());
        if (response.error) {
            console.error(`Error: ${response.error}`);
        } else {
            console.log(`Server response: ${response.results} (Type: ${response.result_type})`);
        }
    } catch (error) {
        console.error(`Invalid response format: ${error}`);
    }
}

// クライアントの開始
function startClient() {
    console.log(`Connecting to ${serverAddress}`);

    const client = net.createConnection(serverAddress, () => {
        console.log('Connected to server.');

        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });

        // プロンプト設定と初回表示
        process.stdout.write('> Enter JSON filename: ');

        // クライアント入力を受け付けるループ
        rl.on('line', async (input) => {
            if (input === 'exit') {
                console.log('Exiting client.');
                rl.close();
                client.end();
                return;
            }

            const filepath = `${basePath}${input}.json`;
            try {
                const req = await readFile(filepath);
                const data = JSON.stringify(req);
                client.write(data);
                console.log('Request sent to server.');
            } catch (err) {
                console.error(err);
                process.stdout.write('> Enter JSON filename: '); // プロンプトを再表示
            }
        });
    });

    // サーバーからのレスポンスを処理
    client.on('data', (data: Buffer) => {
        handleServerResponse(data);

        // ユーザー入力プロンプトを再表示
        process.stdout.write('> Enter JSON filename: ');
    });

    // エラー処理
    client.on('error', (err: Error) => {
        console.error('Connection error:', err.message);
    });

    // 通信終了時
    client.on('end', () => {
        console.log('Disconnected from server.');
    });
}

// 実行
startClient();
