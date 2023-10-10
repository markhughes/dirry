package libadf

import (
	"encoding/binary"
	"errors"
	"os"
)

type AdfEntry struct {
	EntryId int
	Offset  int
	Length  int
}

func UnpackAdfFromFile(path string) (map[int][]byte, error) {
	adfData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return UnpackAdf(adfData)
}

// incomplete
func UnpackAdf(adfData []byte) (map[int][]byte, error) {
	const ADF_MAGIC = 0x00051607
	const ADF_VERSION = 0x00020000

	if len(adfData) < 26 {
		return nil, errors.New("adfData too short")
	}

	magic := binary.BigEndian.Uint32(adfData[0:4])
	if magic != ADF_MAGIC {
		return nil, errors.New("AppleDouble magic number not found")
	}

	version := binary.BigEndian.Uint32(adfData[4:8])
	if version != ADF_VERSION {
		return nil, errors.New("not a Version 2 ADF")
	}

	filler := adfData[8:24]
	numEntries := binary.BigEndian.Uint16(adfData[24:26])

	entries := map[int][]byte{
		0: filler, // Entry #0 is invalid -- use it for the filler..
	}

	adfData = adfData[26:]

	for i := 0; i < int(numEntries); i++ {
		if len(adfData) < 12 {
			return nil, errors.New("adfData too short for entry")
		}

		entryId := int(binary.BigEndian.Uint32(adfData[0:4]))
		offset := int(binary.BigEndian.Uint32(adfData[4:8]))
		length := int(binary.BigEndian.Uint32(adfData[8:12]))

		adfData = adfData[12:]

		if offset+length > len(adfData) {
			return entries, errors.New("entry data extends beyond end of adfData")
		}

		entries[entryId] = adfData[offset : offset+length]
	}

	return entries, nil
}
