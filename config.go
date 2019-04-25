package main

import (
	"flag"
	"os"
	"time"

	"github.com/mattn/go-mastodon"
)

var (
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

func parseFlags() {
	flag.IntVar(&maxQueuedLines, "max-queued-lines", maxQueuedLines, "Pause DF-AI when the queue length gets above this threshold.")
	flag.IntVar(&minQueuedLines, "min-queued-lines", minQueuedLines, "Unpause DF-AI when the queue length gets below this threshold.")
	flag.IntVar(&minLinesBeforeDuplicate, "min-lines-before-duplicate", minLinesBeforeDuplicate, "Minimum number of lines before a duplicate is allowed.")
	flag.IntVar(&maxFuzzyDuplicates, "max-fuzzy-duplicates", maxFuzzyDuplicates, "Maximum number of “fuzzy” (some words changed) duplicates allowed.")
	flag.IntVar(&fuzzyDuplicateWindow, "fuzzy-duplicate-window", fuzzyDuplicateWindow, "Number of lines to remember for “fuzzy” duplicate checking.")
	flag.IntVar(&maxFuzzyDifferentWords, "max-fuzzy-different-words", maxFuzzyDifferentWords, "Maximum number of words that can differ in a “fuzzy” duplicate.")
	flag.DurationVar(&tootInterval, "toot-interval", tootInterval, "Make a toot this often.")

	const (
		placeholderServer = "https://mastodon.example"
		placeholderToken  = "[missing]"
	)

	flag.StringVar(&config.Server, "server", placeholderServer, "Mastodon server name (required)")
	flag.StringVar(&config.ClientID, "client-id", placeholderToken, "OAuth2 Client ID (required)")
	flag.StringVar(&config.ClientSecret, "client-secret", placeholderToken, "OAuth2 Client Secret (required)")
	flag.StringVar(&config.AccessToken, "access-token", placeholderToken, "OAuth2 Access Token (required)")

	flag.Parse()

	if !isExampleMode && (config.Server == placeholderServer || config.ClientID == placeholderToken || config.ClientSecret == placeholderToken || config.AccessToken == placeholderToken) {
		_, _ = os.Stderr.WriteString("server, client-id, client-secret, and access-token are required.\n\n")
		flag.Usage()
		os.Exit(2)
	}
}
