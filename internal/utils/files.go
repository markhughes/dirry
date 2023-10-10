package utils

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/markhughes/dirry/internal/consts"
)

func SaveChunkToFile(chunkType string, offset int, index int, shockwaveFilePath string, data string, prefix string) error {
	outputFolder := filepath.Join(consts.PathDump, filepath.Base(shockwaveFilePath), "resources", chunkType)
	os.MkdirAll(outputFolder, os.ModePerm)

	outputFile := filepath.Join(outputFolder, prefix+strconv.Itoa(index)+"_"+strconv.Itoa((offset))+".json")
	err := os.WriteFile(outputFile, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func SaveChunkToFileBetter(chunkType string, offset int, index int, shockwaveFilePath string, data string, pkg string, prefix string) error {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, filepath.Base(shockwaveFilePath), "resources", chunkType)
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", filepath.Base(shockwaveFilePath), "resources", chunkType)
	}

	os.MkdirAll(outputFolder, os.ModePerm)

	outputFile := filepath.Join(outputFolder, prefix+strconv.Itoa(index)+"_"+strconv.Itoa((offset))+".json")
	err := os.WriteFile(outputFile, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func SaveAfterburnerBinToFile(chunkType string, offset int, index int, shockwaveFilePath string, data string, prefix string) error {
	outputFolder := filepath.Join(consts.PathDump, filepath.Base(shockwaveFilePath), "decompressed", "afterburner")
	os.MkdirAll(outputFolder, os.ModePerm)

	outputFile := filepath.Join(outputFolder, chunkType+".bin")
	err := os.WriteFile(outputFile, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}
