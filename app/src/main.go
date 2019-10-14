package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"

	"github.com/dghubble/oauth1"

	_ "github.com/go-sql-driver/mysql"
)

func myGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalln(key, "not exists")
	}
	return value
}

func main() {
	password := myGetEnv("MYSQL_ROOT_PASSWORD")
	intervalString := myGetEnv("INTERVAL")
	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		log.Fatalln(err)
	}
	consumerKey := myGetEnv("TWITTER_CONSUMER_KEY")
	consumerSecret := myGetEnv("TWITTER_CONSUMER_SECRET")
	accessToken := myGetEnv("TWITTER_ACCESS_TOKEN")
	accessSecret := myGetEnv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	db, err := sql.Open("mysql", "root:"+password+"@tcp(db)/db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	for true {
		check(db, client)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
