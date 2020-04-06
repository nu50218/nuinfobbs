package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/caarlos0/env"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

// env
type config struct {
	GCPProjectID string   `env:"GCP_PROJECT_ID"`
	TargetURL    string   `env:"TARGET_URL"`
	DefaultDone  bool     `env:"DEFAULT_DONE"`
	JobTags      []string `env:"JOB_TAGS" envSeparetor:","`
	Topic        string   `env:"GCP_PUBSUB_TOPIC"`
}

func main() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		log.Fatalln(err)
	}

	store, err := jobutils.NewFireStore(conf.GCPProjectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	pubsubClient, err := pubsub.NewClient(context.Background(), conf.GCPProjectID)
	if err != nil {
		log.Fatalln(err)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	posts, err := crawl(conf.TargetURL)
	if err != nil {
		log.Fatalln(err)
	}

	submittedNewPost := false

	for _, post := range posts {
		for _, tag := range conf.JobTags {
			job := jobutils.NewJob(post.Number, post.Title, post.URL, tag, conf.DefaultDone)
			ok, err := store.SubmitJobIfNotExist(job)
			if err != nil {
				log.Fatalln(err)
			}
			submittedNewPost = submittedNewPost || ok
		}
	}

	if submittedNewPost {
		if _, err := pubsubClient.Topic(conf.Topic).Publish(
			context.Background(),
			&pubsub.Message{
				Data: []byte("unused_data"),
			},
		).Get(context.Background()); err != nil {
			log.Fatal(err)
		}
	}
}
