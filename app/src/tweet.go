package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func check(db *sql.DB, client *twitter.Client) {
	rows, err := db.Query("select number,title,url from posts where tweeted = 0 order by number asc;")
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var number int
		var title, url string
		if err := rows.Scan(&number, &title, &url); err != nil {
			log.Println(err)
			continue
		}

		if err := tweet(title, url, client); err != nil {
			log.Println("failed to tweet:", number)
			apiError, isAPIError := err.(twitter.APIError)
			if !isAPIError {
				log.Println(err)
				continue
			}
			for _, detail := range apiError.Errors {
				log.Println(detail.Code, detail.Message)
				if detail.Code == 187 {
					if _, err := db.Query("update posts set tweeted = 1 where number = " + strconv.Itoa(number)); err != nil {
						log.Println(err)
						continue
					}
				}
			}
			continue
		}

		if _, err := db.Query("update posts set tweeted = 1 where number = " + strconv.Itoa(number)); err != nil {
			ok := false
			for index := 0; index < 10; index++ {
				if _, err := db.Query("update posts set tweeted = 1 where number = " + strconv.Itoa(number)); err != nil {
					time.Sleep(1 * time.Second)
					continue
				}
				ok = true
				break
			}
			if !ok {
				log.Println("failed to update to tweeted:", number)
				break
			}
		}
		log.Println("tweeted:", number)
		time.Sleep(1 * time.Second)
	}
}

func tweet(title, url string, client *twitter.Client) error {
	status := strings.Join([]string{"-表題-", title, "", "-URL-", url}, "\n")
	_, _, err := client.Statuses.Update(status, nil)
	return err
}
