

```
go run main.go
=> ONE PIECEの作者は **尾田栄一郎** です。 
```

以下のf1の部分をf1~f4に変更して実行可能
```
func run(ctx context.Context) error {
	return f1(ctx)
}
```

geminiのAPIキーは作成済みで、環境変数に設定済み
```
echo $GEMINI_API_KEY
```
