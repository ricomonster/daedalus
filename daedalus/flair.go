package daedalus

import (
	"fmt"
	"math/rand"
	"time"
)

type label struct {
	icon  string
	label string
}

var labels = []label{
	{">>>", "sifting through changes..."},
	{"...", "reading the diff entrails..."},
	{"---", "taking inventory..."},
	{"<?>", "inspecting the damage..."},
	{"<~>", "rummaging through your mess..."},
	{"[#]", "cataloguing the chaos..."},
}

func PrintChangedFiles(files []string) {
	l := labels[rand.Intn(len(labels))]

	fmt.Printf("  %s  \033[2m%s\033[0m\n", l.icon, l.label)
	for _, f := range files {
		fmt.Printf("  \033[2m↳\033[0m %s\n", f)
		time.Sleep(40 * time.Millisecond)
	}
}

func WithSpinner(label string, fn func() error) error {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	done := make(chan error, 1)
	go func() { done <- fn() }()

	fmt.Printf("%s  %s  0.0s\n", frames[0], label)

	start := time.Now()
	i := 0
	for {
		select {
		case err := <-done:
			fmt.Printf("\033[1A\r\033[K")
			return err
		default:
			fmt.Printf("\033[1A\r%s  %s  %.1fs\n", frames[i%len(frames)], label, time.Since(start).Seconds())
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
