// +build !example

package main

import (
	"context"
	"html"
	"log"
	"strings"

	"github.com/mattn/go-mastodon"
)

const isExampleMode = false

func initClient() *mastodon.Client {
	return mastodon.NewClient(config)
}

func pullExistingStatuses(ctx context.Context, buffer *dataBuffer, client *mastodon.Client) {
	if minLinesBeforeDuplicate == 0 && fuzzyDuplicateWindow == 0 {
		return
	}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		panic(err)
	}

	i := minLinesBeforeDuplicate - 1
	j := fuzzyDuplicateWindow - 1

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
			if k := strings.Index(s.Content, "</p>"); k != -1 {
				line := html.UnescapeString(s.Content[len("<p>"):k])
				log.Println("Loaded recent toot:", line)

				if i >= 0 {
					buffer.recent[i] = line
					i--
				}

				if j >= 0 {
					buffer.fuzzy[j] = strings.Fields(line)
					j--
				}

				if i < 0 && j < 0 {
					return
				}
			}
		}
	}
}

func makeToot(ctx context.Context, client *mastodon.Client, message string) {
	for attempts := 0; attempts < 5; attempts++ {
		_, err := client.PostStatus(ctx, &mastodon.Toot{
			Status: message,
		})
		if err == nil {
			return
		}

		log.Println("Failed to make toot:", err)
		log.Println("Attempt", attempts+1, "of 5.")
	}
	log.Println("Giving up on toot:", message)
}
