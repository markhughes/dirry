package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/utils"
)

type FmapChunk struct {
	Reader       *binary_reader.BinaryReader
	Fonts        []*Font
	NamesLength  uint32
	EntriesUsed  uint32
	EntriesTotal uint32
}

type Font struct {
	Platform uint16
	FontID   int16
	Name     string
}

func (chunk *FmapChunk) Read(endian binary.ByteOrder) error {
	var err error

	chunk.Fonts = make([]*Font, 0)

	mapLength, err := chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading map length: %s", err)
	}

	utils.DebugMsg("Fmap", fmt.Sprintf("Map length: %d", mapLength))

	chunk.NamesLength, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading names length: %s", err)
	}

	utils.DebugMsg("Fmap", fmt.Sprintf("Names length: %d", chunk.NamesLength))

	bodyStart := chunk.Reader.Pos()
	namesStart := bodyStart + int64(mapLength)

	chunk.Reader.ReadUInt32(endian) // ?
	chunk.Reader.ReadUInt32(endian) // ?

	chunk.EntriesUsed, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading entries used: %s", err)
	}

	chunk.EntriesTotal, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading entries total: %s", err)
	}

	_, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading unknown: %s", err)
	}

	_, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading unknown: %s", err)
	}

	_, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading unknown: %s", err)
	}

	for i := uint32(0); i < chunk.EntriesUsed; i++ {
		nameOffset, err := chunk.Reader.ReadUInt32(endian)
		if err != nil {
			return fmt.Errorf("error reading name offset: %s", err)
		}

		returnPos := chunk.Reader.Pos()

		chunk.Reader.Seek(namesStart+int64(nameOffset), io.SeekStart)

		nameLength, err := chunk.Reader.ReadUInt32(endian)
		if err != nil {
			return fmt.Errorf("error reading name length: %s", err)
		}

		nameBytes, err := chunk.Reader.ReadBytes(int(nameLength))
		if err != nil {
			return fmt.Errorf("error reading name bytes: %s", err)
		}

		name := string(nameBytes)

		chunk.Reader.Seek(returnPos, io.SeekStart)

		platform, err := chunk.Reader.ReadUInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading platform: %s", err)
		}

		id, err := chunk.Reader.ReadUInt16(endian)

		if err != nil {
			return fmt.Errorf("error reading font ID: %s", err)
		}

		// Map cast font ID to window manager font ID
		font := &Font{
			Platform: platform,
			FontID:   int16(id),
			Name:     name,
		}

		chunk.Fonts = append(chunk.Fonts, font)

	}

	return nil

}

func ReadFmapChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*FmapChunk, error) {
	var err error
	chunk := &FmapChunk{
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

func (c *FmapChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
