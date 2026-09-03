/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ricomonster/daedalus/daedalus"
	"github.com/spf13/cobra"
)

var caliperCmd = &cobra.Command{
	Use:   "caliper",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := daedalus.NewCaliperService()

		versionPath, err := cmd.Flags().GetString("path")
		if err != nil {
			return fmt.Errorf("read --path flag: %w", err)
		}

		source, err := app.DetectVersionSource(versionPath)
		if err != nil {
			return fmt.Errorf("Error getting version: %w", err)
		}

		// Get the version
		version, err := source.Read()
		if err != nil {
			return fmt.Errorf("Error reading version source: %w", err)
		}

		cmd.Println("Current version:", version)

		bump, err := getBumpType(cmd)
		if err != nil {
			return err
		}

		if bump != "" {
			next, err := app.Bump(bump, version)
			if err != nil {
				return err
			}

			cmd.Printf("Bumping version %s -> %s\n", version, next)
			if err := source.Write(next); err != nil {
				return err
			}

			cmd.Printf("Version bumped to %s\n", next)
		}
		return nil
	},
}

func getBumpType(cmd *cobra.Command) (daedalus.VersionBumpType, error) {
	opts := []struct {
		flag string
		bump daedalus.VersionBumpType
	}{
		{"major", daedalus.VersionBumpTypeMajor},
		{"minor", daedalus.VersionBumpTypeMinor},
		{"patch", daedalus.VersionBumpTypePatch},
	}

	var selected daedalus.VersionBumpType
	for _, opt := range opts {
		enabled, err := cmd.Flags().GetBool(opt.flag)
		if err != nil {
			return "", fmt.Errorf("read --%s: %w", opt.flag, err)
		}

		if !enabled {
			continue
		}

		if selected != "" {
			return "", fmt.Errorf("version bump flags are mutually exclusive")
		}

		selected = opt.bump
	}

	return selected, nil
}

func init() {
	rootCmd.AddCommand(caliperCmd)

	caliperCmd.Flags().String("path", "", "file or project directory")

	caliperCmd.Flags().Bool("major", false, "bump major version")
	caliperCmd.Flags().Bool("minor", false, "bump minor version")
	caliperCmd.Flags().Bool("patch", false, "bump patch version")
	caliperCmd.Flags().Bool("auto", false, "determine bump automatically")
}
