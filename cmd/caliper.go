/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ricomonster/daedalus/daedalus"
	"github.com/spf13/cobra"
)

var versionPath string

var caliperCmd = &cobra.Command{
	Use:   "caliper",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		app := daedalus.NewCaliperService()
		source, err := app.DetectVersionSource(versionPath)
		if err != nil {
			cmd.PrintErrln("Error getting version:", err)
			return
		}

		// Get the version
		version, err := source.Read()
		if err != nil {
			cmd.PrintErrln("Error reading version source:", err)
			return
		}

		cmd.Println("Current version:", version)
	},
}

func init() {
	rootCmd.AddCommand(caliperCmd)

	caliperCmd.Flags().StringVar(&versionPath, "path", "", "file or project directory")
}
