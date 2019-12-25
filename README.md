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

Container Registryのイメージを、Cloud Runにデプロイします。最大1リクエストかつスケール数も最大１にしてください。

#### 環境変数の設定

`containers/*/src/main.go`の`config`という構造体を見れば必要な環境変数がわかります。

ひと目でわからなそうな環境変数だけ説明します。

##### crawler

| 環境変数 | 説明 |
| -- | -- |
| TARGET_URL | クロール対象の掲示板を開いたトップページのURL |
| DEFAULT_DONE | 既に投稿を配信済みとしてDBに投げるかを書きます。(true/false) |
| JOB_TAGS | 他のコンテナ向けのタグを`,`区切りで入力します。例）`twitter,slack,line` |

##### crawler以外

| 環境変数 | 説明 |
| -- | -- |
| TAG | タグを指定します。このタグに一致するポストでDoneがfalseのものを配信していきます。例）`twitter` |

### Cloud Schedulerの設定

Cloud Schedulerで定期的にCloud RunのURLを叩きます。
