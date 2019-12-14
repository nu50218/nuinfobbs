package main

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/nu50218/nuinfobbs/library/jobutils"
)

func tweet(client *twitter.Client, post *jobutils.Post) error {
	status := strings.Join([]string{"-表題-", post.Title, "", "-URL-", post.URL}, "\n")
	_, _, err := client.Statuses.Update(status, nil)
	return err
}
