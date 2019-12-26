package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nu50218/nuinfobbs/library/jobutils"
	"golang.org/x/sync/errgroup"
)

const monthlyMessageLimit = 1000

type config struct {
	GCPProjectID           string `env:"GCP_PROJECT_ID"`
	Tag                    string `env:"TAG"`
	LineChannelSecret      string `env:"TWITTER_CHANNEL_SECRET"`
	LineChannelAccessToken string `env:"TWITTER_CHANNEL_ACCESS_TOKEN"`
}

func do() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		log.Fatalln(err)
	}

	bot, err := linebot.New(conf.LineChannelSecret, conf.LineChannelAccessToken)

	store, err := jobutils.NewFireStore(conf.GCPProjectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	jobs, err := store.GetWaitingJobsByTag(conf.Tag)

	followers, err := bot.GetNumberFollowers(time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02")).Do()
	if err != nil {
		log.Fatalln(err)
	}

	messageConsumption, err := bot.GetMessageConsumption().Do()
	if err != nil {
		log.Fatalln(err)
	}
	if messageConsumption.TotalUsage+followers.Followers > monthlyMessageLimit {
		return
	}

	// いくつかまとめて配信することでリミットを超えないようにしたい
	if int64(len(jobs))*int64(followers.Followers) > monthlyMessageLimit {
		return
	}

	posts := []*jobutils.Post{}

	for _, job := range jobs {
		posts = append(posts, job.Post)
	}

	if err := push(bot, posts); err != nil {
		log.Fatalln(err)
	}

	var eg errgroup.Group

	for _, job := range jobs {
		eg.Go(func() error {
			if err := store.MakeJobDone(job); err != nil {
				var succeeded bool
				for i := 0; i < 10; i++ {
					if err := store.MakeJobDone(job); err == nil {
						succeeded = true
						break
					}
				}
				if !succeeded {
					return err
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatalln(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	do()
}

func main() {
	log.Println("", "started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
