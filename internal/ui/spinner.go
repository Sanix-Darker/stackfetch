package ui

import (
	"fmt"
	"time"
)

// i got inspired by this project :https://github.com/sindresorhus/cli-spinners
var frames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

type Spinner struct {
	stop chan struct{}
}

func NewSpinner() *Spinner { return &Spinner{stop: make(chan struct{})} }

func (s *Spinner) Start(msg string) {
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\r%s %c", msg, frames[i%len(frames)])
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop() { close(s.stop) }
