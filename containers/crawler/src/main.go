package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

func do() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := jobutils.NewFirestoreClient(conf.GCPProjectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()

	posts, err := crawl(conf.TargetURL)
	if err != nil {
		log.Fatalln(err)
	}

	for _, post := range posts {
		for _, tag := range conf.JobTags {
			job := jobutils.NewJob(post.Number, post.Title, post.URL, tag, conf.DefaultDone)
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
