package daedalus

import (
	"fmt"
	"time"
)

func WithSpinner(label string, fn func() error) error {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan error, 1)
	go func() { done <- fn() }()

	i := 0
	for {
		select {
		case err := <-done:
			fmt.Printf("\r\033[K") // clear line
			return err
		default:
			fmt.Printf("\r%s  %s", frames[i%len(frames)], label)
			i++
			time.Sleep(80 * time.Millisecond)
		}
	}
}
