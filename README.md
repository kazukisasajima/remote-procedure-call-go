# remote-procedure-call-go

## 概要
クライアントとサーバが異なるプログラミング言語で書かれていても、クライアントプログラムがサーバ上の機能を呼び出せるRPC(remote procedure call)のシステムを作成しました。<br>
このプロジェクトはコンピュータサイエンス学習サービス[Recursion](https://recursion.example.com)の課題でPythonを用いて作成したものをGoで作成しました。<br>
クライアント側はGoとTypeScriptの2つ作成しました。


## サーバーによるRPC関数の提供
- floor(double x): 10 進数 x を最も近い整数に切り捨て、その結果を整数で返す。
- nroot(int n, int x): 方程式 rn = x における、r の値を計算する。
- reverse(string s): 文字列 s を入力として受け取り、入力文字列を逆にした新しい文字列を返す。
- validAnagram(string str1, string str2): 2 つの文字列を入力として受け取り，2 つの入力文字列が互いにアナグラムであるかどうかを示すブール値を返す。
- sort(string[] strArr): 文字列の配列を入力として受け取り、その配列をソートして、ソート後の文字列の配列を返す。
- 存在しないメソッド名を入力するとクライアント側でエラーになります。

## 実行方法
### 1. ターミナルを2つ開きます。

- 1つ目のターミナルでサーバーを起動します
```sh
go run main.go server
```

- 2つ目のターミナルでクライアントを起動します
```sh
go run main.go client
```
もしくは
```sh
npx ts-node ts_client/src/client.ts
```

クライアント側のターミナルでメソッド名を入力してサーバー側の関数を呼び出す<br>
<img width="1200" alt="remote_procedure_call" src="https://github.com/user-attachments/assets/a1a42045-cafc-4e8e-853f-588e450e60d5">

### 注意
- 名前付きパイプ（UNIXソケット）を使用しているため、Windowsの標準的なコマンドプロンプトやPowerShellでは動作しません。
- このアプリケーションはLinuxまたはWSL（Windows Subsystem for Linux）環境で動作します。
