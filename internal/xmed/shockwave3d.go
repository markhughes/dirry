package xmed

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

func CreateShockwave3DBinary(reader io.Reader, endian binary.ByteOrder) ([]byte, error) {
	fourcc, err := utils.ReadString(reader, 4, false)
	if err != nil {
		return nil, fmt.Errorf("could not read value1: %v", err)
	}

	chunkLength, err := utils.ReadUInt32(reader, endian) // not sure what this is
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
	utils.DebugMsg("xmed/shockwave3d", "fourcc: %s\n", fourcc)
	utils.DebugMsg("xmed/shockwave3d", "chunkLength: %d\n", chunkLength)
	utils.DebugMsg("xmed/shockwave3d", "unknown2: %d\n", unknown2)
	utils.DebugMsg("xmed/shockwave3d", "dataLength: %d\n", dataLength)

	data := make([]byte, dataLength)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, fmt.Errorf("could not read data: %v", err)
	}

	return data, nil
}
