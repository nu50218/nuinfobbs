package main

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func post(client *slack.Client, channelID string, post *jobutils.Post) error {
	message := strings.Join([]string{post.Title, post.URL}, "\n")
	_, _, err := client.PostMessage(channelID, slack.MsgOptionText(message, false))
	return err
}
