package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

type config struct {
	GCPProjectID          string `env:"GCP_PROJECT_ID"`
	Tag                   string `env:"TAG"`
	TwitterConsumerKey    string `env:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret string `env:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken    string `env:"TWITTER_ACCESS_TOKEN"`
	TwitterAccessSecret   string `env:"TWITTER_ACCESS_SECRET"`
}

func do() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		log.Fatalln(err)
	}

	config := oauth1.NewConfig(conf.TwitterConsumerKey, conf.TwitterConsumerSecret)
	token := oauth1.NewToken(conf.TwitterAccessToken, conf.TwitterAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	store, err := jobutils.NewFireStore(conf.GCPProjectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	jobs, err := store.GetWaitingJobsByTag(conf.Tag)
	for _, job := range jobs {
		if err := tweet(client, job.Post); err != nil {
			log.Fatalln(err)
		}
		if err := store.MakeJobDone(job); err != nil {
			var succeeded bool
			for i := 0; i < 10; i++ {
				if err := store.MakeJobDone(job); err == nil {
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
	log.Println("twitter", "started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
