package main

import (
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func push(client *linebot.Client, posts []*jobutils.Post) error {
	messages := make([]string, len(posts))
	for _, post := range posts {
		messages = append(messages, strings.Join([]string{post.Title, post.URL}, "\n"))
	}
	_, err := client.BroadcastMessage(linebot.NewTextMessage(strings.Join(messages, "\n"))).Do()
	return err
}
