// +build linux

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigch := make(chan os.Signal, 1)
		defer signal.Stop(sigch)

		signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigch:
			log.Println(sig, "- cleaning up")
		case <-ctx.Done():
		}

		cancel()
	}()

	buffer := &dataBuffer{
		recent: make([]string, minLinesBeforeDuplicate),
		fuzzy:  make([][]string, fuzzyDuplicateWindow),
		queue:  make([]string, 0, maxQueuedLines),
	}

	client := initClient()
	pullExistingStatuses(ctx, buffer, client)

	ch := make(chan string)
	go dwarfFortress(ctx, buffer, ch)

	initialDelay := nextDelay()
	log.Println("Waiting", initialDelay, "before making first toot.")
	time.Sleep(initialDelay)

	for {
		var line string
		select {
		case <-ctx.Done():
			return
		case line = <-ch:
		default:
			log.Println("Warning: no toots are ready")
			select {
			case <-ctx.Done():
				return
			case line = <-ch:
			}
		}

		makeToot(ctx, client, line)
		time.Sleep(nextDelay())
	}
}

func nextDelay() time.Duration {
	return tootInterval - time.Duration(time.Now().UnixNano())%tootInterval
}
