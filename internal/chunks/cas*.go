package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type CasEntry struct {
	Index int32
}

type CasChunk struct {
	Reader  *binary_reader.BinaryReader
	Entries []CasEntry
}

func (chunk *CasChunk) Read(endian binary.ByteOrder) error {
	var err error
	entryCount := chunk.Reader.Length / 4

	chunk.Entries = make([]CasEntry, entryCount)

	for i := 0; i < int(entryCount); i++ {
		entry := &CasEntry{}

		entry.Index, err = chunk.Reader.ReadInt32(endian)
		if err != nil {
			return err
		}

		chunk.Entries[i] = *entry
	}

	return nil

}

func ReadCasChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*CasChunk, error) {
	var err error
	chunk := &CasChunk{
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

func (c *CasChunk) Print() {
	fmt.Printf("-- CAS*\n")
	fmt.Printf("%d entries total\n", len(c.Entries))

	fmt.Printf("Index | Value \n")
	fmt.Printf("--------------\n")
	for i, entry := range c.Entries {
		fmt.Printf(" %4d | %5d \n", i+1, entry.Index)
	}
}

func (c *CasChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
