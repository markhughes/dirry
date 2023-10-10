package chunks

import (
	"encoding/binary"
	"encoding/json"
	"strings"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type Offset struct {
	Frame  int16
	Offset int16
}

type VwlbChunk struct {
	Reader *binary_reader.BinaryReader

	Labels map[int16]string
}

func (chunk *VwlbChunk) Read(endian binary.ByteOrder) error {
	chunk.Labels = make(map[int16]string)

	numOffsets, err := chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	offsetMap := make([]Offset, numOffsets)

	for i := int16(0); i < numOffsets; i++ {
		offsetMap[i].Frame, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

		offsetMap[i].Offset, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

	}

	labelsLen, err := chunk.Reader.ReadInt32(endian)
	if err != nil {
		return err
	}

	labelsBytes, err := chunk.Reader.ReadBytes(int(labelsLen))
	if err != nil {
		return err
	}
	labels := string(labelsBytes)

	for i := 0; i < len(offsetMap); i++ {
		start := int(offsetMap[i].Offset)
		var end int

		if i == len(offsetMap)-1 {
			end = len(labels)
		} else {
			end = int(offsetMap[i+1].Offset)
		}

		chunk.Labels[offsetMap[i].Frame] = strings.TrimSpace(labels[start:end])
	}

	return nil
}

func ReadVwlbChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*VwlbChunk, error) {
	var err error
	chunk := &VwlbChunk{
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

func (c *VwlbChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
