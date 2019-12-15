# nuinfobbs

名大情報学部掲示板（[@nuinfobbs](https://twitter.com/nuinfobbs)）を動かすやつです。

[LINEbot](lin.ee/qqzeIFh)もつくりました。

学部Slackにも配信しています。

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
