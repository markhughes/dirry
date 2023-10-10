//go:build !js

package cmd

import (
	"github.com/markhughes/dirry/internal/adf"
	"github.com/spf13/cobra"
)

var adfCmd = &cobra.Command{
	Use:   "adf <filePath>",
	Short: "(unstable) A Mac System 1-7 apple double file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		PreRunHandler()
		adf.Dump(args[0])

		return nil
	},
}

func init() {
	rootCmd.AddCommand(adfCmd)
}
