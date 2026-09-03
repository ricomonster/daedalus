/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ricomonster/daedalus/config"
	"github.com/ricomonster/daedalus/daedalus"
	"github.com/ricomonster/daedalus/gemini"
	"github.com/ricomonster/daedalus/git"
	"github.com/spf13/cobra"
)

var spinning = []string{
	"Getting things ready",
	"Making progress",
	"Working on it",
	"Putting things together",
	"Taking a closer look",
	"Making some adjustments",
	"Processing your request",
	"Almost there",
	"Finishing up",
	"Wrapping things up",
}

var stroking = []string{
	"Preparing your changes",
	"Gathering everything together",
	"Taking one last look",
	"Putting things in order",
	"Saving your progress",
	"Getting everything ready",
	"Making it official",
	"Sealing the deal",
	"Wrapping things up",
	"Almost finished",
}

// stylusCmd represents the stylus command
var stylusCmd = &cobra.Command{
	Use:   "stylus",
	Short: "Generate a commit message from your staged git changes.",
	Long:  `Analyzes your staged git diff and uses an LLM to generate a commit message following the Conventional Commits spec. Stage your changes with git add, then run this command to get a ready-to-use commit message or full PR description.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		// Here you will define your flags and configuration settings.
		conf, err := config.New()
		if err != nil {
			fmt.Printf("failed to load config: %v\n", err)
			os.Exit(1)
		}

		// llm
		ge := gemini.New(conf)

		// services
		gi := git.New()

		// apps
		sa := daedalus.NewStylusService(gi, ge)

		changes, err := sa.GetChanges(cmd.Context())
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		daedalus.PrintChangedFiles(changes.Files)

		// Timeout after 2 minutes
		ctx, cancel := context.WithTimeout(cmd.Context(), 120*time.Second)
		defer cancel()

		spin := spinning[rand.Intn(len(spinning))]

		var commit string
		if err := daedalus.WithSpinner(spin, start, func() error {
			var e error
			commit, e = sa.GetCommitMessage(ctx, changes)
			return e
		}); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		stroke := stroking[rand.Intn(len(stroking))]

		// Commit the changes
		var commitOut []byte
		err = daedalus.WithInkStroke(fmt.Sprintf("%s...", stroke), start, func() error {
			var err error
			commitOut, err = sa.Commit(cmd.Context(), commit)

			return err
		})

		fmt.Printf("\n%s", bytes.TrimRight(commitOut, "\r\n"))

		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		// Check if we need to push
		push, _ := cmd.Flags().GetBool("push")
		if push {
			if err := gi.Push(); err != nil {
				fmt.Printf("error: %v\n", err)
				os.Exit(1)
			}

			elapsed := time.Since(start)
			fmt.Printf("\n✅  Committed and pushed! (%.1fs)\n", elapsed.Seconds())
			return
		}

		elapsed := time.Since(start)
		fmt.Printf("\n✅  Ready to push! (%.1fs)\n", elapsed.Seconds())
	},
}

func init() {
	rootCmd.AddCommand(stylusCmd)

	stylusCmd.Flags().Bool("push", false, "push current branch after committing")
}
