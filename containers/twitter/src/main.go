package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/nu50218/nuinfobbs/library/env"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func tweet(client *twitter.Client, post *jobutils.Post) error {
	status := strings.Join([]string{"-表題-", post.Title, "", "-URL-", post.URL}, "\n")
	_, _, err := client.Statuses.Update(status, nil)
	return err
}

func do() {
	projectID := env.AssertAndGetEnv("GCP_PROJECT_ID")
	tag := env.AssertAndGetEnv("TAG")

	consumerKey := env.AssertAndGetEnv("TWITTER_CONSUMER_KEY")
	consumerSecret := env.AssertAndGetEnv("TWITTER_CONSUMER_SECRET")
	accessToken := env.AssertAndGetEnv("TWITTER_ACCESS_TOKEN")
	accessSecret := env.AssertAndGetEnv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	firestoreClient, err := jobutils.NewFirestoreClient(projectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()

	jobs, err := jobutils.GetWaitingJobsByTag(firestoreClient, tag)
	for _, job := range jobs {
		if err := tweet(client, job.Post); err != nil {
			log.Fatalln(err)
		}
		if err := jobutils.MakeJobDone(firestoreClient, job); err != nil {
			var succeeded bool
			for i := 0; i < 10; i++ {
				if err := jobutils.MakeJobDone(firestoreClient, job); err == nil {
					succeeded = true
					break
				}
			}
			if !succeeded {
				log.Fatalln(err)
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	do()
}

func main() {
	tag := env.AssertAndGetEnv("TAG")

	log.Println("twitter", tag, "started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
