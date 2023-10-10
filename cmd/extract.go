//go:build !js

package cmd

import (
	"github.com/markhughes/dirry/internal/dzip"
	"github.com/spf13/cobra"
)

var zipCmd = &cobra.Command{
	Use:   "zip <filePath>",
	Short: "Creates a zip file of the raw director export.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		PreRunHandler()

		dzip.DZip(args[0], "")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
