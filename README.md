# Sentient Dwarf Fortress was inevitable.

The source code that powers <a href="https://botsin.space/@it_was_inevitable" rel="me">@it\_was\_inevitable</a> and <a href="https://botsin.space/@it_was_inevitable_slow" rel="me">@it\_was\_inevitable\_slow</a> on botsin.space in the Fediverse.

See the entry on [BotWiki.org](https://botwiki.org/bot/it_was_inevitable/).

## To run locally:

1. If you haven't already, install Docker.
2. Create an application on your Mastodon instance on the bot account.
3. Have at least these scopes selected:
  - `read:accounts` and `read:statuses` (or just `read`)
  - `write:statuses` (or just `write`)
4. Copy [config.example.go](config.example.go) and make the following changes:
  - Change `example` in the `// +build example` line to some other word (for example, `yourtag`).
  - Replace the placeholder data in the config object with the data from your
    Mastodon application. (The application interface lists the data in the same
    order it is defined in config.example.go.)
  - Optional: Modify the other constants in the file.
5. Run [./run.bash yourtag](run.bash)

## To run locally without a Mastodon account:

1. If you haven't already, install Docker.
2. Run [./run.bash example](run.bash)
3. When you're done watching the example, push Ctrl+C.
