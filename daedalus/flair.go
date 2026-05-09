package daedalus

import (
	"fmt"
	"time"
)

func WithSpinner(label string, fn func() error) error {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan error, 1)
	go func() { done <- fn() }()

	start := time.Now()
	i := 0
	for {
		select {
		case err := <-done:
			fmt.Printf("\r\033[K") // clear line
			return err
		default:
			fmt.Printf("\r%s  %s  %.1fs", frames[i%len(frames)], label, time.Since(start).Seconds())
			i++
			time.Sleep(80 * time.Millisecond)
		}
	}
}

func WithInkStroke(label string, fn func() error) error {
	frames := []string{"▱▱▱▱▱", "▰▱▱▱▱", "▰▰▱▱▱", "▰▰▰▱▱", "▰▰▰▰▱", "▰▰▰▰▰"}
	for _, f := range frames {
		fmt.Printf("\r%s  Committing...", f)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\r\033[K")

	return fn()
}
