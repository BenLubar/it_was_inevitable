// +build example
// [Remove the previous line to enable this file.]

package main

import (
	"time"

	"github.com/mattn/go-mastodon"
)

const (
	// Pause DF-AI when the queue length gets above this threshold.
	maxQueuedLines = 1000

	// Unpause DF-AI when the queue length gets below this threshold.
	minQueuedLines = 500

	// Minimum number of lines before a duplicate is allowed.
	minLinesBeforeDuplicate = 10

	// Make a toot this often.
	tootInterval = 5 * time.Minute
)

var config = &mastodon.Config{
	Server:       "https://mastodon.example",
	ClientID:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	ClientSecret: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	AccessToken:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
}
