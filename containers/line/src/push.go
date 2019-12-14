package main

import (
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func push(client *linebot.Client, post *jobutils.Post) error {
	message := strings.Join([]string{post.Title, post.URL}, "\n")
	_, err := client.BroadcastMessage(linebot.NewTextMessage(message)).Do()
	return err
}
