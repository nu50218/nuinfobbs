# nuinfobbs

名古屋大学の情報学部の掲示板をクロールして、[Twitter](https://twitter.com/nuinfobbs)、[LINE](lin.ee/qqzeIFh)、Slackに更新情報を配信します。

`containers`のそれぞれが独立して動いています。

## 使い方

GCPのCloud RunとCloud Schedulerを使います。

### 準備

GCPで新規プロジェクトを作成して、コマンドラインで`gcloud`が動くようにします。

### コンテナのビルド

`containers`ディレクトリの中のディレクトリ（`crawler`など）それぞれに対して、`Dockerfile`が存在する階層で

`$ gcloud builds submit --config cloudbuild.yaml .`

を実行します。

そうすると、Cloud BuildでビルドされてContainer Registryにpushされます。

### Cloud Runにデプロイ

Cloud Runにデプロイします。

#### 環境変数の設定

WIP

### Cloud Schedulerの設定

Cloud Schedulerで定期的にCloud RunのURLを叩きます。
