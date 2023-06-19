package main

import (
	"context"
	"strings"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func tweet(client *twitter.Client, post *jobutils.Post) error {
	_, err := client.CreateTweet(context.Background(), twitter.CreateTweetRequest{
		Text: strings.Join([]string{"-表題-", post.Title, "", "-URL-", post.URL}, "\n"),
	})
	return err
}
