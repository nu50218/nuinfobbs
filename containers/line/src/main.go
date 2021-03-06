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
	LineChannelSecret      string `env:"LINE_CHANNEL_SECRET"`
	LineChannelAccessToken string `env:"LINE_CHANNEL_ACCESS_TOKEN"`
}

func main() {
	conf := config{}
	if err := env.Parse(&conf); err != nil {
		log.Println(err)
		return
	}

	bot, err := linebot.New(conf.LineChannelSecret, conf.LineChannelAccessToken)
	if err != nil {
		log.Println(err)
		return
	}

	store, err := jobutils.NewFireStore(conf.GCPProjectID)
	if err != nil {
		log.Println(err)
		return
	}
	defer store.Close()

	jobs, err := store.GetWaitingJobsByTag(conf.Tag)
	if err != nil {
		log.Println(err)
		return
	}
	if len(jobs) == 0 {
		return
	}

	// 集計は翌日中に終わるらしいので、2日前で取得
	followers, err := bot.GetNumberFollowers(time.Now().Add(-48 * time.Hour).In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("20060102")).Do()
	if err != nil {
		log.Println(err)
		return
	}
	if followers.Status != "ready" {
		log.Println("followers.Status is not ready")
		return
	}

	messageConsumption, err := bot.GetMessageConsumption().Do()
	if err != nil {
		log.Println(err)
		return
	}
	if messageConsumption.TotalUsage+followers.Followers > monthlyMessageLimit {
		return
	}

	// いくつかまとめて配信することでリミットを超えないようにしたい
	// 月間の投稿数は最大80個程度と仮定
	// 1000 < 予想月間投稿数 ＝ 180/len(jobs)*友達の数ならやめる
	if monthlyMessageLimit*int64(len(jobs)) < int64(80)*followers.Followers {
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

	for i := range jobs {
		job := jobs[i]
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
		log.Println(err)
		return
	}

}
