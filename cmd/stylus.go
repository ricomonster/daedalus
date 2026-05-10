/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ricomonster/daedalus/application"
	"github.com/ricomonster/daedalus/config"
	"github.com/ricomonster/daedalus/daedalus"
	"github.com/ricomonster/daedalus/gemini"
	"github.com/ricomonster/daedalus/git"
)

var sa daedalus.StylusApplication

// stylusCmd represents the stylus command
var stylusCmd = &cobra.Command{
	Use:   "stylus",
	Short: "Generate a commit message from your staged git changes.",
	Long:  `Analyzes your staged git diff and uses an LLM to generate a commit message following the Conventional Commits spec. Stage your changes with git add, then run this command to get a ready-to-use commit message or full PR description.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		if sa == nil {
			fmt.Println("failed to initialize stylus application")
			os.Exit(1)
		}

		changes, err := sa.GetChanges(cmd.Context())
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		daedalus.PrintChangedFiles(changes.Files)

		// Timeout after 2 minutes
		ctx, cancel := context.WithTimeout(cmd.Context(), 120*time.Second)
		defer cancel()

		var commit string
		if err := daedalus.WithSpinner("Oracle is checking", func() error {
			var e error
			commit, e = sa.GetCommitMessage(ctx, changes)
			return e
		}); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		// Commit the changes
		if err := daedalus.WithInkStroke("Committing...", func() error {
			return sa.Commit(cmd.Context(), commit)
		}); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		elapsed := time.Since(start)
		fmt.Printf("\n✅  Ready to push! (%.1fs)\n", elapsed.Seconds())
	},
}

func init() {
	rootCmd.AddCommand(stylusCmd)

	// Here you will define your flags and configuration settings.
	conf, err := config.New()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	// llm
	ge := gemini.New(conf)

	// services
	gi, _ := git.New()

	// apps
	sa = application.NewStylusApplication(gi, ge)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stylusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stylusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
