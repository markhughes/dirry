//go:build !js

package cmd

import (
	"github.com/markhughes/dirry/internal/dump"
	"github.com/spf13/cobra"
)

var dump2Cmd = &cobra.Command{
	Use:   "dump <filePath>",
	Short: "Dumps the binary file, extracts resources, and any other files describing them.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		PreRunHandler()

		dump.Dump(args[0], "", 0)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dump2Cmd)
}
