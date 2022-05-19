# stockcode-scraping

(Golang練習用AWSLambda関数)

Yahoo!ファイナンスをスクレイピングして銘柄一覧を作るバッチ処理。

### 環境
* Docker
* AWS SAM CLI
* Go 1.18
### 実行方法
```bash
# ローカル
cd stockcode-scraping
sam build
cd ..
sam local start-lambda
aws lambda invoke --function-name StockScraping --endpoint http://127.0.0.1:3001/ output.txt
```