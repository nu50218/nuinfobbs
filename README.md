# nuinfobbs

名大情報学部掲示板（[@nuinfobbs](https://twitter.com/nuinfobbs)）を動かすやつ

## 説明

イメージ3つが動いていて、**crawler**が掲示板をクロールして**db**に突っ込んでいって**app**が`tweeted=0`なものを監視してツイートする。

## 使い方

### `env_file`を書く

#### .env_app

`MYSQL_ROOT_PASSWORD`は同じにします。

```sh
MYSQL_ROOT_PASSWORD= # SQLのrootのパスワード
INTERVAL= # DBを監視する間隔（秒）
# twitter
TWITTER_CONSUMER_KEY=
TWITTER_CONSUMER_SECRET=
TWITTER_ACCESS_TOKEN=
TWITTER_ACCESS_SECRET=
```

### .env_crawler

```sh
MYSQL_ROOT_PASSWORD= # SQLのrootのパスワード
TARGET_URL= # 対象掲示板のURL
INTERVAL= # 掲示板アクセスの間隔（秒）
```

### .env_db

```sh
MYSQL_ROOT_PASSWORD= # SQLのrootのパスワード
```

### docker-compose up する

`$ docker-compose up`
