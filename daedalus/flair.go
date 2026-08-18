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

	fmt.Printf("%s  \033[2m%s\033[0m\n", l.icon, l.label)
	for _, f := range files {
		fmt.Printf("  \033[2m↳\033[0m %s\n", f)
		time.Sleep(40 * time.Millisecond)
	}
}

func WithSpinner(label string, start time.Time, fn func() error) error {
	return withAnimation(label, start, []string{
		"⠋", "⠙", "⠹", "⠸", "⠼",
		"⠴", "⠦", "⠧", "⠇", "⠏",
	}, 80*time.Millisecond, fn)
}

func WithInkStroke(label string, start time.Time, fn func() error) error {
	return withAnimation(label, start, []string{
		"▱▱▱▱▱", "▰▱▱▱▱", "▰▰▱▱▱",
		"▰▰▰▱▱", "▰▰▰▰▱", "▰▰▰▰▰",
	}, 100*time.Millisecond, fn)
}

func withAnimation(
	label string,
	start time.Time,
	frames []string,
	delay time.Duration,
	fn func() error,
) error {
	done := make(chan error, 1)
	go func() { done <- fn() }()

	for i := 0; ; i++ {
		select {
		case err := <-done:
			fmt.Printf("\r\033[K")
			return err
		default:
			fmt.Printf(
				"\r%s  %s  %.1fs",
				frames[i%len(frames)],
				label,
				time.Since(start).Seconds(),
			)
			time.Sleep(delay)
		}
	}
}
