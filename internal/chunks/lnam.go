package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type LnamChunk struct {
	Reader *binary_reader.BinaryReader

	FileSize  uint32
	FileSize2 uint32

	Unknown1 uint16

	NamesLength uint16
	Names       []string
}

func (chunk *LnamChunk) Read(endian binary.ByteOrder) error {
	var err error

	chunk.Reader.ReadUInt32(endian) // ?
	chunk.Reader.ReadUInt32(endian) // ?

	chunk.FileSize, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("failed to read lnam chunk file size: %s", err)
	}

	chunk.FileSize2, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("failed to read lnam chunk file size 2: %s", err)
	}

	chunk.Unknown1, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return fmt.Errorf("failed to read lnam chunk unknown 1: %s", err)
	}

	chunk.NamesLength, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return fmt.Errorf("failed to read lnam chunk names length: %s", err)
	}

	if chunk.FileSize != chunk.FileSize2 {
		return fmt.Errorf("file size mismatch: %d != %d", chunk.FileSize, chunk.FileSize2)
	}

	for i := 0; i < int(chunk.NamesLength); i++ {
		nameLen, err := chunk.Reader.ReadUInt8()
		if err != nil {
			return fmt.Errorf("failed to read name length: %s", err)
		}

		name, err := chunk.Reader.ReadString(int(nameLen), endian)
		if err != nil {
			return fmt.Errorf("failed to read name: %s", err)
		}

		chunk.Names = append(chunk.Names, name)
	}

	return nil

}
func ReadLnamChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*LnamChunk, error) {

	var err error
	chunk := &LnamChunk{
		Reader: r,
	}

	chunk.Reader.HexDump(true)

	// Always big endian?
	r.Seek(0, 0)
	err = chunk.Read(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

func (c *LnamChunk) ToJSON() (string, error) {

	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
