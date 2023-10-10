package chunks

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

type KeyRecord struct {
	ElementIndex int32
	CastIndex    int32
	ChunkType    string
	CastNumber   int32
}

func (r *KeyRecord) IsValid() bool {
	_, ok := ChunkTypes[r.ChunkType]
	return ok
}

type KeyChunk struct {
	Chunk         *BinaryChunk
	HeaderSize    int16
	RecordSize    int16
	RecordCount   int32
	ActiveRecords int32
	Records       []*KeyRecord
}

func ReadKeyChunkRaw(r *bytes.Reader, endian binary.ByteOrder) (*KeyChunk, error) {

	var err error
	chunk := &KeyChunk{}
	chunk.Chunk, err = FromBinaryAtPartHeadless(r, "KEY*", endian)
	if err != nil {
		return nil, err
	}

	err = processKeyChunk(chunk, endian)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

func ReadKeyChunk(r io.ReadSeeker, endian binary.ByteOrder, offset int64) (*KeyChunk, error) {
	_, err := r.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	chunk := &KeyChunk{}

	chunk.Chunk, err = FromBinaryAt(r, "KEY*", endian, offset)
	if err != nil {
		return nil, err
	}

	err = processKeyChunk(chunk, endian)
	if err != nil {
		return nil, fmt.Errorf("failed to process key chunk: %s", err)
	}
	return chunk, nil
}

func processKeyChunk(chunk *KeyChunk, endian binary.ByteOrder) error {
	var err error

	chunk.HeaderSize, err = chunk.Chunk.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.RecordSize, err = chunk.Chunk.ReadInt16(endian)
	if err != nil {
		return err
	}
	chunk.RecordCount, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return err
	}
	chunk.ActiveRecords, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return err
	}

	chunk.Records = make([]*KeyRecord, chunk.RecordCount)

	utils.DebugMsg("key", "KEY* records: %d\n", chunk.RecordCount)
	for i := 0; i < int(chunk.RecordCount); i++ {
		rec := &KeyRecord{}

		rec.ElementIndex, err = chunk.Chunk.ReadInt32(endian)
		if err != nil {
			return err
		}

		rec.CastIndex, err = chunk.Chunk.ReadInt32(endian)
		if err != nil {
			return err
		}

		chunkType, err := chunk.Chunk.ReadBytes(4)
		if err != nil {
			return err
		}

		if endian == binary.LittleEndian {
			rec.ChunkType = string(utils.ReverseBytes(chunkType[:]))
		} else {
			rec.ChunkType = string((chunkType[:]))

		}
		if rec.ElementIndex >= 1024 {
			rec.CastNumber = rec.ElementIndex - 1024
		} else {
			rec.CastNumber = -1
		}

		chunk.Records[i] = rec

		utils.DebugMsg("key", "KEY* record %d: elementindex - %d, castindex - %d, type - %s, castNo - %d\n", i, rec.ElementIndex, rec.CastIndex, rec.ChunkType, rec.CastNumber)
	}

	return nil

}

func (c *KeyChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
