// +build linux

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mattn/go-mastodon"
)

func main() {
	cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigch := make(chan os.Signal, 1)
		defer signal.Stop(sigch)

		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-sigch:
			log.Println("Caught", sig, "- cleaning up")
		case <-ctx.Done():
		}

		cancel()
	}()

	client := mastodon.NewClient(config)

	ch := make(chan string)
	go dwarfFortress(ctx, client, ch)

	initialDelay := tootInterval - time.Duration(time.Now().UnixNano())%tootInterval
	log.Println("Waiting", initialDelay, "before making first toot.")
	time.Sleep(initialDelay)

	for {
		makeToot(ctx, client, <-ch)
		time.Sleep(tootInterval)
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
