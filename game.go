// +build linux

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hpcloud/tail"
	"github.com/kr/pty"
)

func dwarfFortress(ctx context.Context, ch chan<- string) {
	addch := make(chan string, 100)
	debugch := make(chan struct{})
	buffer := make([]string, 0, maxQueuedLines)
	lastThreshold := 0
	var recent [minLinesBeforeDuplicate]string
	go watchLog(ctx, addch)

	runGame := func() {
		cmd := exec.CommandContext(ctx, "/df_linux/dfhack")
		cmd.Env = append(cmd.Env, "TERM=xterm-256color")
		cmd.Dir = "/df_linux"

		f, err := pty.Start(cmd)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()

		exited := make(chan error, 1)
		go func() {
			exited <- cmd.Wait()
		}()

		defer func() {
			if err := clearSaves(); err != nil {
				log.Println("Removing saved worlds:", err)
			}
		}()

		wasRunning := true

		for {
			var nextLine string
			out := ch
			if len(buffer) == 0 {
				out = nil
			} else {
				nextLine = buffer[0]
			}

			select {
			case <-time.After(time.Minute * 30):
				if wasRunning {
					log.Println("30 minutes without any log output. Assuming DF is stuck. Resetting.")
					if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
						log.Println("Sending SIGTERM:", err)
					}
					select {
					case <-time.After(time.Minute):
						log.Println("Process has not exited. Killing.")
						if err := cmd.Process.Signal(syscall.SIGKILL); err != nil {
							log.Println("Sending SIGKILL:", err)
						}
					case <-exited:
						return
					}
				}
			case out <- nextLine:
				buffer = buffer[1:]
				if len(buffer)%100 == 0 && lastThreshold >= len(buffer)+100 {
					lastThreshold = len(buffer)
					log.Println(len(buffer), "lines are buffered.")
				}
				if len(buffer) < minQueuedLines {
					if !wasRunning {
						log.Println("Unsuspending Dwarf Fortress")
						wasRunning = true
					}
					if err := cmd.Process.Signal(syscall.SIGCONT); err != nil {
						log.Println("Sending SIGCONT:", err)
					}
				}
			case line := <-addch:
				if minLinesBeforeDuplicate != 0 {
					firstLine := line[:strings.IndexByte(line, '\n')]
					var duplicate bool
					for _, dup := range recent {
						if dup == firstLine {
							duplicate = true
							break
						}
					}
					if duplicate {
						continue
					}
					copy(recent[:], recent[1:])
					recent[len(recent)-1] = firstLine
				}
				buffer = append(buffer, line)
				if len(buffer)%100 == 0 && lastThreshold <= len(buffer)-100 {
					lastThreshold = len(buffer)
					log.Println(len(buffer), "lines are buffered.")
				}
				if len(buffer) > maxQueuedLines {
					if wasRunning {
						log.Println("Suspending Dwarf Fortress")
						wasRunning = false
					}
					if err := cmd.Process.Signal(syscall.SIGSTOP); err != nil {
						log.Println("Sending SIGSTOP:", err)
					}
				}
			case err := <-exited:
				log.Println("Dwarf Fortress process exited:", err)
				return
			case <-debugch:
				if err := cmd.Process.Signal(syscall.SIGKILL); err != nil {
					log.Println("Sending SIGKILL:", err)
				}
				if err := os.Remove("/df_linux/df-ai-debug.log"); err != nil {
					log.Println("Removing debug log:", err)
				}
				return
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		runGame()
	}

}

func watchDebug(ctx context.Context, ch chan<- struct{}) {
	for {
		// Wait at least a minute between kills to make sure there's time to clean up.
		time.Sleep(time.Minute)

		f, err := tail.TailFile("/df_linux/df-ai-debug.log", tail.Config{
			Follow: true,
		})
		if err != nil {
			panic(err)
		}

		first := <-f.Lines
		log.Println("df-ai crashed! debug log follows:")
		go func() {
			// ignore error
			_ = f.StopAtEOF()
			f.Cleanup()
		}()
		fmt.Println(first.Text)
		for line := range f.Lines {
			fmt.Println(line.Text)
		}

		ch <- struct{}{}
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
