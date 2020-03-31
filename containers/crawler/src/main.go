package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

// env
type config struct {
	GCPProjectID string   `env:"GCP_PROJECT_ID"`
	TargetURL    string   `env:"TARGET_URL"`
	DefaultDone  bool     `env:"DEFAULT_DONE"`
	JobTags      []string `env:"JOB_TAGS" envSeparetor:","`
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

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	posts, err := crawl(conf.TargetURL)
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		for _, tag := range conf.JobTags {
			job := jobutils.NewJob(post.Number, post.Title, post.URL, tag, conf.DefaultDone)
			if err := store.SubmitJobIfNotExist(job); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
