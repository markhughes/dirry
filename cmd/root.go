//go:build !js

package cmd

import (
	"fmt"
	"strings"

	"github.com/markhughes/dirry/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dirry",
	Short: "Dirry is a tool for parsing director files",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return setup(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("No command specified.")
			fmt.Println(cmd.UsageString())
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

var verbose string
var logging bool

func setup(cmd *cobra.Command) error {
	verbose, _ = cmd.Flags().GetString("verbose")
	logging, _ = cmd.Flags().GetBool("logging")

	return nil
}

func PreRunHandler() {
	if logging {
		utils.EnableLogging = true
	}

	if verbose != "" {
		if verbose == "all" {
			utils.EnabledDebugAll = true
		} else {
			categories := strings.Split(verbose, ",")
			for _, category := range categories {
				utils.EnabledDebugCategoriesMap[category] = true
			}
		}
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("logging", "l", false, "Enable extra log files")
	rootCmd.PersistentFlags().StringP("verbose", "v", "", "Enable verbose output for categories (use 'all' for all categories)")
}
