// +build example
// [Change the tag in the previous line to enable this file.]

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
	minLinesBeforeDuplicate = 500

	// Maximum number of "fuzzy" (some words changed) duplicates allowed.
	maxFuzzyDuplicates = 5

	// Number of lines to remember for "fuzzy" duplicate checking.
	fuzzyDuplicateWindow = 10

	// Maximum number of words that can differ in a "fuzzy" duplicate.
	maxFuzzyDifferentWords = 2

	// Make a toot this often.
	tootInterval = 5 * time.Second
)

var config = &mastodon.Config{
	Server:       "https://mastodon.example",
	ClientID:     "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	ClientSecret: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	AccessToken:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
}
