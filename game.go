// +build linux

package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hpcloud/tail"
	"github.com/kr/pty"
	"github.com/mattn/go-mastodon"
)

func logError(err error, message string) {
	if err != nil {
		log.Println(message, err)
	}
}

func logErrorD(f func() error, message string) {
	logError(f(), message)
}

type dataBuffer struct {
	queue     []string
	threshold int
	recent    [minLinesBeforeDuplicate]string
	running   bool
}

func dwarfFortress(ctx context.Context, client *mastodon.Client, ch chan<- string) {
	addch := make(chan string, 100)
	debugch := make(chan struct{})

	buffer := &dataBuffer{
		queue: make([]string, 0, maxQueuedLines),
	}

	go watchLog(ctx, addch)
	go watchDebug(ctx, debugch)

	pullExistingStatuses(ctx, buffer, client)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		runGame(ctx, buffer, ch, addch, debugch)
	}
}

func pullExistingStatuses(ctx context.Context, buffer *dataBuffer, client *mastodon.Client) {
	if minLinesBeforeDuplicate == 0 {
		return
	}

	account, err := client.GetAccountCurrentUser(ctx)
	if err != nil {
		panic(err)
	}

	statuses, err := client.GetAccountStatuses(ctx, account.ID, &mastodon.Pagination{
		Limit: minLinesBeforeDuplicate,
	})
	if err != nil {
		panic(err)
	}

	i := minLinesBeforeDuplicate - 1
	for _, s := range statuses {
		if !strings.HasPrefix(s.Content, "<p>") {
			continue
		}
		if j := strings.Index(s.Content, "</p>"); j != -1 {
			buffer.recent[i] = html.UnescapeString(s.Content[len("<p>"):j])
			log.Println("Loaded recent toot:", buffer.recent[i])
			i--
		}
	}
}

func runGame(ctx context.Context, buffer *dataBuffer, ch chan<- string, addch <-chan string, debugch <-chan struct{}) {
	cmd := exec.CommandContext(ctx, "/df_linux/dfhack")
	cmd.Env = append(cmd.Env, "TERM=xterm-256color")
	cmd.Dir = "/df_linux"

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer logErrorD(f.Close, "Closing Pseudo-TTY:")

	exited := make(chan error, 1)
	go func() {
		exited <- cmd.Wait()
	}()

	defer logErrorD(clearSaves, "Removing saved worlds:")

	signal := func(s os.Signal) {
		logError(cmd.Process.Signal(s), "Sending "+s.String()+":")
	}

	buffer.running = true
	first := true

	for {
		var nextLine string
		out := ch
		if len(buffer.queue) == 0 {
			out = nil
		} else {
			nextLine = buffer.queue[0]
		}

		select {
		case <-time.After(time.Minute * 30):
			if buffer.running {
				log.Println("30 minutes without any log output. Assuming DF is stuck. Resetting.")
				signal(syscall.SIGTERM)
				select {
				case <-time.After(time.Minute):
					log.Println("Process has not exited. Killing.")
					signal(syscall.SIGKILL)
				case <-exited:
					return
				}
			}

		case out <- nextLine:
			buffer.queue = buffer.queue[1:]
			checkQueueLength(buffer, signal)

		case line := <-addch:
			if isDuplicate(buffer, line) {
				continue
			}

			if first || len(buffer.queue) == 0 {
				log.Println("First toot queued:", line)
				first = false
			}

			buffer.queue = append(buffer.queue, line)
			checkQueueLength(buffer, signal)

		case err := <-exited:
			log.Println("Dwarf Fortress process exited:", err)
			return

		case <-debugch:
			signal(syscall.SIGKILL)
			if err := os.Remove("/df_linux/df-ai-debug.log"); err != nil {
				log.Println("Removing debug log:", err)
			}
			return
		}
	}
}

func isDuplicate(buffer *dataBuffer, line string) bool {
	if minLinesBeforeDuplicate == 0 {
		return false
	}

	firstLine := line[:strings.IndexByte(line, '\n')]
	for _, dup := range buffer.recent {
		if dup == firstLine {
			return true
		}
	}

	copy(buffer.recent[:], buffer.recent[1:])
	buffer.recent[len(buffer.recent)-1] = firstLine
	return false
}

func checkQueueLength(buffer *dataBuffer, signal func(os.Signal)) {
	if l := len(buffer.queue); l%100 == 0 && (buffer.threshold >= l+100 || buffer.threshold <= l-100) {
		buffer.threshold = l
		log.Println(l, "lines are buffered.")
	}

	if len(buffer.queue) < minQueuedLines {
		if !buffer.running {
			log.Println("Unsuspending Dwarf Fortress")
			buffer.running = true
		}
		signal(syscall.SIGCONT)
	} else if len(buffer.queue) > maxQueuedLines {
		if buffer.running {
			log.Println("Suspending Dwarf Fortress")
			buffer.running = false
		}
		signal(syscall.SIGSTOP)
	}
}

func watchDebug(ctx context.Context, ch chan<- struct{}) {
	for {
		// Wait at least a minute between kills to make sure there's time to clean up.
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
		}

		f, err := tail.TailFile("/df_linux/df-ai-debug.log", tail.Config{
			Follow: true,
		})
		if err != nil {
			panic(err)
		}

		select {
		case <-ctx.Done():
			return
		case first := <-f.Lines:
			log.Println("df-ai crashed! debug log follows:")
			go func() {
				f.Kill(f.StopAtEOF())
				f.Cleanup()
			}()
			fmt.Println(first.Text)
			for line := range f.Lines {
				fmt.Println(line.Text)
			}

			ch <- struct{}{}
		}
	}
}

func watchLog(ctx context.Context, ch chan<- string) {
	f, err := tail.TailFile("/df_linux/gamelog.txt", tail.Config{
		ReOpen: true,
		Follow: true,
	})
	if err != nil {
		panic(err)
	}
	defer f.Cleanup()

	go func() {
		<-ctx.Done()
		f.Kill(ctx.Err())
	}()

	for line := range f.Lines {
		if line.Err != nil {
			log.Println(line.Err)
			continue
		}

		text := mapCP437(strings.TrimSpace(line.Text))

		// Chatter messages are formatted as "Urist McName, Occupation: It was inevitable."
		if i, j := strings.Index(text, ", "), strings.Index(text, ": "); i > 0 && j > i {
			ch <- text[j+2:] + "\n\n\u2014 " + text[:j]
		}
	}
}