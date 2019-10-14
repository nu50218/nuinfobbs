package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

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
	targetURL := myGetEnv("TARGET_URL")
	intervalString := myGetEnv("INTERVAL")
	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("mysql", "root:"+password+"@tcp(db)/db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	
	if err != nil {
		log.Fatalln(err)
	}

	for true {
		crawl(targetURL, db)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
