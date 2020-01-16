package main

import (
	"log"
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

func main() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		log.Fatalln(err)
	}

	bot, err := linebot.New(conf.LineChannelSecret, conf.LineChannelAccessToken)
	if err != nil {
		log.Fatalln(err)
	}

	store, err := jobutils.NewFireStore(conf.GCPProjectID)
	if err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	jobs, err := store.GetWaitingJobsByTag(conf.Tag)
	if err != nil {
		log.Fatalln(err)
	}
	if len(jobs) == 0 {
		return
	}

	followers, err := bot.GetNumberFollowers(time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("20060102")).Do()
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
	// 月間の投稿数は50個と仮定
	// 1000 < 予想月間投稿数 ＝ 50/len(jobs)*友達の数ならやめる
	if monthlyMessageLimit*int64(len(jobs)) < int64(50)*followers.Followers {
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
