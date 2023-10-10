package chunks

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

type FgeiChunk struct {
	Chunk    *BinaryChunk
	Position int64
}

func ReadFGEIChunk(r io.ReadSeeker, endian binary.ByteOrder, abmp *ABMPChunk) (*FgeiChunk, error) {
	var err error

	chunk := &FgeiChunk{}

	chunk.Chunk, err = FromBinary(r, "FGEI", endian)
	if err != nil {

		return nil, err
	}

	chunk.Position, err = r.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	utils.DebugMsg("fgei", "FGEI chunk at [%d:%d]\n", chunk.Position, chunk.Chunk.Length)

	return chunk, nil

}

func (c *FgeiChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
