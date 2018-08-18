// +build example

package main

import (
	"context"
	"log"
)

func initClient() struct{} {
	return struct{}{}
}

func pullExistingStatuses(ctx context.Context, buffer *dataBuffer, client struct{}) {
	log.Println("Running in example mode. Not connecting to Mastodon.")
}

func makeToot(ctx context.Context, client struct{}, message string) {
	log.Println("TOOT:", message)
}
