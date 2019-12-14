# nuinfobbs

名大情報学部掲示板（[@nuinfobbs](https://twitter.com/nuinfobbs)）を動かすやつ
学部SlackとかLINEbotとかでも配信したい

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

WIP

### Cloud Schedulerの設定

WIP
