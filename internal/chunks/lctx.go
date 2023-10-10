package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type ScriptFile struct {
	Unknown1   int32
	Index      int32
	Flags      int16
	NextUnused int16
}
type LctxChunk struct {
	Reader *binary_reader.BinaryReader

	Unknown1 uint32
	Unknown2 uint32

	NScripts  uint32
	NScripts2 uint32

	Offset      uint16
	ScriptIndex uint16
	ScriptFiles []ScriptFile

	TableId int32

	FirstUnused int16
}

func (chunk *LctxChunk) Read(endian binary.ByteOrder) error {
	var err error

	chunk.Unknown1, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk unknown1: %s", err)
	}

	chunk.Unknown2, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk unknown2: %s", err)
	}

	chunk.NScripts, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk nscripts: %s", err)
	}

	chunk.NScripts2, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk nscripts2: %s", err)
	}

	chunk.Offset, err = chunk.Reader.ReadUInt16(endian) // items offset
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk offset: %s", err)
	}

	chunk.Reader.ReadUInt16(endian) // entry size
	chunk.Reader.ReadUInt32(endian) // unk1
	chunk.Reader.ReadUInt32(endian) // file type
	chunk.Reader.ReadUInt32(endian) // unk2

	chunk.TableId, err = chunk.Reader.ReadInt32(endian) // name table id
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk table id: %s", err)
	}
	chunk.Reader.ReadInt16(endian)                          // valid count
	chunk.Reader.ReadUInt16(endian)                         // flags
	chunk.FirstUnused, err = chunk.Reader.ReadInt16(endian) // first unused
	if err != nil {
		return fmt.Errorf("failed to read lctx chunk first unused: %s", err)
	}

	if chunk.NScripts != chunk.NScripts2 {
		return fmt.Errorf("script count mismatch: %d != %d", chunk.NScripts, chunk.NScripts2)
	}

	chunk.Reader.Seek(int64(chunk.Offset), 0)

	for i := 0; i < int(chunk.NScripts); i++ {
		scriptFile := ScriptFile{}

		// Could this contain a bit field?
		scriptFile.Unknown1, err = chunk.Reader.ReadInt32(endian)
		if err != nil {
			return fmt.Errorf("failed to read unknown: %s", err)
		}

		scriptFile.Index, err = chunk.Reader.ReadInt32(endian)
		if err != nil {
			return fmt.Errorf("failed to read script file: %s", err)
		}

		scriptFile.Flags, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return fmt.Errorf("failed to read unknown: %s", err)
		}

		scriptFile.NextUnused, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return fmt.Errorf("failed to read unknown: %s", err)
		}

		chunk.ScriptFiles = append(chunk.ScriptFiles, scriptFile)

	}

	return nil

}

func ReadLctxChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*LctxChunk, error) {

	var err error
	chunk := &LctxChunk{
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

func (c *LctxChunk) ToJSON() (string, error) {

	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
