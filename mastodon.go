// +build !example

package main

import (
	"context"
	"html"
	"log"
	"strings"

	"github.com/mattn/go-mastodon"
)

func initClient() *mastodon.Client {
	return mastodon.NewClient(config)
}

func pullExistingStatuses(ctx context.Context, buffer *dataBuffer, client *mastodon.Client) {
	if minLinesBeforeDuplicate == 0 {
		return
	}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		panic(err)
	}

	i := minLinesBeforeDuplicate - 1

	pg := &mastodon.Pagination{}

	for pg.Limit == 0 {
		pg.Limit = int64(i + 1)

		statuses, err := client.GetAccountStatuses(ctx, account.ID, pg)
		if err != nil {
			panic(err)
		}

		for _, s := range statuses {
			if !strings.HasPrefix(s.Content, "<p>") {
				continue
			}
			if j := strings.Index(s.Content, "</p>"); j != -1 {
				buffer.recent[i] = html.UnescapeString(s.Content[len("<p>"):j])
				log.Println("Loaded recent toot:", buffer.recent[i])
				i--

				if i < 0 {
					return
				}
			}
		}
	}
}

func makeToot(ctx context.Context, client *mastodon.Client, message string) {
	_, err := client.PostStatus(ctx, &mastodon.Toot{
		Status: message,
	})
	if err != nil {
		log.Println(err)
	}
}
