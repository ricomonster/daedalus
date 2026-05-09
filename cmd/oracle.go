/*
Copyright © 2026 ricomonster <https://github.com/ricomonster/daedalus>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ricomonster/daedalus/application"
	"github.com/ricomonster/daedalus/config"
	"github.com/ricomonster/daedalus/daedalus"
	"github.com/ricomonster/daedalus/gemini"
)

// Flags
var (
	list bool
	llm  string
	key  string
)

var oracleApp daedalus.OracleApplication

// oracleCmd represents the oracle command
var oracleCmd = &cobra.Command{
	Use:   "oracle",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var oracleConfigCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		if len(key) == 0 {
			// Retrieve the key of whatever llm is provided
			os.Exit(0)
		}

		// Set the key
		err := oracleApp.SetLLMKey(context.Background(), daedalus.LLM(llm), key)
		if err != nil {
			fmt.Printf("failed to save key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("key now stored")
	},
}

func init() {
	rootCmd.AddCommand(oracleCmd)
	oracleCmd.Flags().BoolVar(&list, "list", false, "list of the supported llm's")

	// Here you will define your flags and configuration settings.
	c, err := config.New()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	// llm
	gemini := gemini.New(c)

	// apps
	oracleApp = application.NewOracleApplication(c, gemini)

	// oracle config
	oracleCmd.AddCommand(oracleConfigCmd)
	oracleConfigCmd.Flags().StringVar(&llm, "llm", "", "name of the llm. run (daedalus oracle --list for supported llm's)")
	oracleConfigCmd.Flags().StringVar(&key, "key", "", "api key from the llm service")
	// oracleConfigCmd.MarkFlagRequired("llm")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// oracleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// oracleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
