package xmed

import (
	"encoding/binary"
	"encoding/json"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/pfr"
)

func CreateFontBinary(reader *binary_reader.BinaryReader, endian binary.ByteOrder) (data []byte, meta []byte, err error) {

	var fontMeta = pfr.PfrFont{}

	// read all of reader into bytes
	originalData, err := reader.ReadAllBytes()
	if err != nil {
		return nil, meta, err
	}

	// remove first 4 bytes
	data = originalData[4:]

	// add new 4 bytes for signature PFR0 - for some reason shockwave added PFR1?
	data = append([]byte("PFR0"), data...)

	fontMeta.Parse(originalData)

	// convert meta to json encoded bytes
	metaBytes, err := json.Marshal(fontMeta)
	if err != nil {
		return nil, meta, err
	}

	return data, (metaBytes), nil

}
