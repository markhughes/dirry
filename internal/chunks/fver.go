package chunks

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

type FverChunk struct {
	StartPos        int64
	Chunk           *BinaryChunk
	Version         uint32
	IMapVersion     uint32
	DirectorVersion uint32
}

// This chunk is mostly best-guess at the moment, it doesn't seem useful?
func ReadFverChunk(r io.ReadSeeker, endian binary.ByteOrder) (*FverChunk, error) {
	var err error

	chunk := &FverChunk{}

	chunk.Chunk, err = FromBinary(r, "Fver", endian)
	if err != nil {
		return nil, err
	}

	chunk.StartPos, err = chunk.Chunk.r.Seek(0, io.SeekCurrent) // Get current file pointer position
	if err != nil {
		return chunk, err
	}

	chunk.Chunk.r.Seek(chunk.StartPos, io.SeekStart)
	chunk.Chunk.Data, err = chunk.Chunk.ReadBytes(int(chunk.Chunk.Length))
	if err != nil {
		utils.DebugMsg("fver", "Fver: error reading bytes: %s\n", err)
		return chunk, err
	}

	chunk.Chunk.r.Seek(chunk.StartPos, io.SeekStart)

	chunk.Version, err = chunk.Chunk.ReadVarInt()
	if err != nil {
		return chunk, err
	}
	utils.DebugMsg("fver", "Fver: version: %x\n", chunk.Version)

	chunk.IMapVersion, err = chunk.Chunk.ReadVarInt()
	if err != nil {
		return chunk, err
	}
	utils.DebugMsg("fver", "Fver: iMapVersion: %x\n", chunk.IMapVersion)

	chunk.DirectorVersion, err = chunk.Chunk.ReadVarInt()
	if err != nil {
		return chunk, err
	}
	utils.DebugMsg("fver", "Fver: directorVersion: %x\n", chunk.DirectorVersion)

	end, err := chunk.Chunk.r.Seek(0, io.SeekCurrent) // Get current file pointer position
	if err != nil {
		return chunk, err
	}
	if end-chunk.StartPos != int64(chunk.Chunk.Length) {
		utils.DebugMsg("fver", "Expected Fver of length %d but read %d bytes\n", chunk.Chunk.Length, end-chunk.StartPos)
		chunk.Chunk.r.Seek(chunk.StartPos+int64(chunk.Chunk.Length), io.SeekStart)
	}

	utils.DebugMsg("fver", "Fver chunk length = %d\n", chunk.Chunk.Length)

	return chunk, nil
}

func (c *FverChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
