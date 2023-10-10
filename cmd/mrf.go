//go:build !js

package cmd

import (
	"github.com/markhughes/dirry/internal/mrf"
	"github.com/spf13/cobra"
)

var mrfCmd = &cobra.Command{
	Use:   "mrf <filePath>",
	Short: "(unstable) A Mac System 1-7 resource fork",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		PreRunHandler()
		mrf.Dump(args[0])

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mrfCmd)
}
