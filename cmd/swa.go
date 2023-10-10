//go:build !js

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/swa"
	"github.com/spf13/cobra"
)

var swaCmd = &cobra.Command{
	Use:   "swa",
	Short: "(unstable) Convert an SWA to a MP3 - it's not really working",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		PreRunHandler()
		if len(args) < 1 {
			return cmd.Help()
		}

		var filePath = args[0]
		var outputProjectName = filepath.Base(filePath)
		var outputFileName = outputProjectName
		// if it has a file extension, remove it
		if filepath.Ext(outputFileName) != "" {
			outputFileName = outputFileName[:len(outputFileName)-len(filepath.Ext(outputFileName))]
		}
		outputFileName = outputFileName + ".mp3"

		fmt.Printf("Converting %v\n", filePath)

		bytes, fileExtension, err := swa.Swa2Mp3FromFile(filePath)
		if err != nil {
			return err
		} else {
			outputFolder := filepath.Join(consts.PathDump, filepath.Base(filePath), "resources", fileExtension)
			os.MkdirAll(outputFolder, os.ModePerm)

			outputFile := filepath.Join(outputFolder, outputFileName)

			err = os.WriteFile(outputFile, bytes, 0644)
			if err != nil {
				return err
			}

		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(swaCmd)
}
