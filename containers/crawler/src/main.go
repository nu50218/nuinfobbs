package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nu50218/nuinfobbs/library/env"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func crawl(targetURL string) ([]*jobutils.Post, error) {
	posts := []*jobutils.Post{}

	doc, err := goquery.NewDocument(targetURL)
	if err != nil {
		return nil, err
	}

	selection := doc.Find("table#ichiran > tbody > .ichiran_odd,.ichiran_even")
	selection.EachWithBreak(func(_ int, tr *goquery.Selection) bool {
		post := &jobutils.Post{}
		tr.Find("td").EachWithBreak(func(i int, td *goquery.Selection) bool {
			switch i {
			case 0:
				if post.Number, err = strconv.Atoi(td.Text()); err != nil {
					return false
				}
			case 6:
				a := td.Find("a")
				if a.Length() != 1 {
					err = errors.New("There was not only one <a/> in <td/>")
					return false
				}
				var exists bool
				if post.URL, exists = a.Attr("href"); !exists {
					err = errors.New("There was not href in <a/>")
					return false
				}
				post.Title = a.Text()
			}
			return true
		})
		if err != nil {
			return false
		}
		posts = append(posts, post)
		return true
	})
	return posts, err
}

func do() {
	projectID := env.AssertAndGetEnv("GCP_PROJECT_ID")
	targetURL := env.AssertAndGetEnv("TARGET_URL")
	defalutDone, err := strconv.ParseBool(env.AssertAndGetEnv("DEFAULT_DONE"))
	if err != nil {
		log.Fatalln(err)
	}
	tags := strings.Split(env.AssertAndGetEnv("JOB_TAGS"), ",")

	firestoreClient, err := jobutils.NewFirestoreClient(projectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()

	posts, err := crawl(targetURL)
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		for _, tag := range tags {
			job := jobutils.NewJob(post.Number, post.Title, post.URL, tag, defalutDone)
			if err := jobutils.SubmitJobIfNotExist(firestoreClient, job); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	do()
}

func main() {
	log.Println("crawler", "started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
