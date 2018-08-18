# Sentient Dwarf Fortress was inevitable.

The source code that powers [@it\_was\_inevitable@botsin.space](https://botsin.space/@it_was_inevitable) and [@it\_was\_inevitable\_slow@botsin.space](https://botsin.space/@it_was_inevitable_slow) on Mastodon.

See the entry on [BotWiki.org](https://botwiki.org/bot/it_was_inevitable/).

## To run locally:

1. If you haven't already, install Docker.
2. Create an application on your Mastodon instance on the bot account.
3. Have at least these scopes selected:
  - `read:accounts` and `read:statuses` (or just `read`)
  - `write:statuses` (or just `write`)
4. Copy [config.example.go](config.example.go) and make the following changes:
  - Remove the `// +build example` line.
  - Replace the placeholder data in the config object with the data from your
    Mastodon application. (The application interface lists the data in the same
    order it is defined in config.example.go.)
  - Optional: Modify the other constants in the file.
5. Run [./run.bash](run.bash)
