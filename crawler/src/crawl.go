package main

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type post struct {
	number int
	title  string
	url    string
}

func crawl(targetURL string, db *sql.DB) {
	doc, err := goquery.NewDocument(targetURL)
	if err != nil {
		log.Println(err)
		return
	}

	posts := []*post{}
	selection := doc.Find("table#ichiran > tbody > .ichiran_odd,.ichiran_even")
	selection.Each(func(_ int, tr *goquery.Selection) {
		p := &post{}
		tr.Find("td").EachWithBreak(func(i int, td *goquery.Selection) bool {
			switch i {
			case 0:
				if p.number, err = strconv.Atoi(td.Text()); err != nil {
					return false
				}
			case 6:
				a := td.Find("a")
				if a.Length() != 1 {
					err = errors.New("There was not only one <a/> in <td/>")
					return false
				}
				var exists bool
				if p.url, exists = a.Attr("href"); !exists {
					err = errors.New("There was not href in <a/>")
					return false
				}
				p.title = a.Text()
			}
			return true
		})
		if err != nil {
			log.Println(err)
			return
		}
		posts = append(posts, p)
	})
	if err := insert(db, posts); err != nil {
		log.Println(err)
	}
}

func insert(db *sql.DB, posts []*post) error {
	if len(posts) == 0 {
		return nil
	}
	values := []string{}
	for _, p := range posts {
		value := "(" + regexp.QuoteMeta(strconv.Itoa(p.number)) + ",'" + regexp.QuoteMeta(p.title) + "','" + regexp.QuoteMeta(p.url) + "')"
		values = append(values, value)
	}
	_, err := db.Query("insert ignore into posts (number,title,url) values " + strings.Join(values, ",") + ";")
	return err
}
