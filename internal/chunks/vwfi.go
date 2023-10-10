package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type VwfiChunk struct {
	Reader *binary_reader.BinaryReader

	Unk1  uint32
	Unk2  uint32
	Flags uint32

	ScriptId uint32

	StringsCount uint16
	Strings      []string

	AllowOutdatedLingo      bool
	RemapPalettesWhenNeeded bool

	Script string

	ChangedBy     string
	CreatedBy     string
	OrigDirectory string

	Preload uint16
}

func (chunk *VwfiChunk) Read(endian binary.ByteOrder) error {
	return fmt.Errorf("not implemented")
}

func ReadVwfiChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*VwfiChunk, error) {
	var err error
	chunk := &VwfiChunk{
		Reader: r,
	}

	chunk.Reader.HexDump(true)

	r.Seek(0, 0)
	err = chunk.Read(binary.BigEndian)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

func (c *VwfiChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
