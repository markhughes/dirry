package swa

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/h2non/filetype"
)

func Swa2Mp3FromFile(filePath string) ([]byte, string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, "", err
	}
	return Swa2Mp3FromBytes(fileBytes)
}

func Swa2Mp3FromBytes(swaBytes []byte) ([]byte, string, error) {

	// is it already an mp3?
	kind, _ := filetype.Match(swaBytes)
	if kind.MIME.Type == "mp3" {
		return swaBytes, "mp3", nil
	}

	// AC 44
	ac44 := []byte{0xAC, 0x44}

	// 49 44 33
	newHeader := []byte{0x49, 0x44, 0x33}

	// Locate AC 44 in the input
	ac44Index := bytes.Index(swaBytes, ac44)

	if ac44Index != -1 {
		// Remove everything before and including AC44
		updatedSwaBytes := swaBytes[ac44Index+len(ac44):]

		// Insert 494433 at the start
		result := append(newHeader, updatedSwaBytes...)
		return result, "mp3", nil
	}

	// TODO: I don't understand this container format just yet.. it's an mp3 in some weird container
	//		 ffmpeg seems to be good at finding and extracting?
	return nil, "", errors.New("unknown swa format")
}
