package chunks

import (
	"encoding/binary"
	"encoding/json"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type SordEntry struct {
	CastLibId int16
	MemberId  int16
}

type SordChunk struct {
	Reader  *binary_reader.BinaryReader
	Entries []SordEntry
}

func (chunk *SordChunk) Read(endian binary.ByteOrder) error {
	chunk.Reader.ReadInt32(endian)
	chunk.Reader.ReadInt32(endian)

	length, err := chunk.Reader.ReadInt32(endian)
	if err != nil {
		return err
	}

	chunk.Entries = make([]SordEntry, 0)

	for i := 0; i < int(length); i++ {
		entry := SordEntry{}

		// TODO: in scumm, only if _version >= kFileVer500
		// ottherwise, dont read just set to our default 1
		entry.CastLibId, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

		entry.MemberId, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

		chunk.Entries = append(chunk.Entries, entry)
	}

	return nil

}

func ReadSordChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*SordChunk, error) {
	var err error
	chunk := &SordChunk{
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

func (c *SordChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
