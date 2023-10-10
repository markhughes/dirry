package xmed

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

func CreateFlashBinary(reader io.Reader, endian binary.ByteOrder) ([]byte, error) {
	unknown1, err := utils.ReadUInt32(reader, endian) // not sure what this is
	if err != nil {
		return nil, fmt.Errorf("could not read value1: %v", err)
	}

	unknown2, err := utils.ReadUInt32(reader, endian) // not sure what this is
	if err != nil {
		return nil, fmt.Errorf("could not read value2: %v", err)
	}

	dataLength, err := utils.ReadUInt32(reader, endian) // not sure what this is
	if err != nil {
		return nil, fmt.Errorf("could not read value3: %v", err)
	}

	utils.DebugMsg("xmed/flash", "value1: %d\n", unknown1)
	utils.DebugMsg("xmed/flash", "value2: %d\n", unknown2)
	utils.DebugMsg("xmed/flash", "value3: %d\n", dataLength)

	data := make([]byte, dataLength)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, fmt.Errorf("could not read data: %v", err)
	}
	return data, nil

}
