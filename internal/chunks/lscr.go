package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type LscrChunk struct {
	Reader *binary_reader.BinaryReader

	TotalLength  uint32
	TotalLength2 uint32
	HeaderLength uint16
	ScriptNumber uint16
	Unknown2     int16
	ParentNumber int16

	ScriptFlags          uint32
	Unk42                int16
	CastID               int32
	FactoryNameID        int16
	HandlerVectorsCount  uint16
	HandlerVectorsOffset uint32
	HandlerVectorsSize   uint32
	PropertiesCount      uint16
	PropertiesOffset     uint32
	GlobalsCount         uint16
	GlobalsOffset        uint32
	HandlersCount        uint16
	HandlersOffset       uint32
	LiteralsCount        uint16
	LiteralsOffset       uint32
	LiteralsDataCount    uint32
	LiteralsDataOffset   uint32
}

func (chunk *LscrChunk) Read(endian binary.ByteOrder) error {
	var err error

	chunk.TotalLength, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.TotalLength2, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.HeaderLength, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.ScriptNumber, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.Unknown2, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.ParentNumber, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.Reader.Seek(38, 0)

	chunk.ScriptFlags, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.Unk42, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.CastID, err = chunk.Reader.ReadInt32(endian)
	if err != nil {
		return err
	}

	chunk.FactoryNameID, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.HandlerVectorsCount, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.HandlerVectorsOffset, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.HandlerVectorsSize, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.PropertiesCount, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.PropertiesOffset, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.GlobalsCount, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.GlobalsOffset, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.HandlersCount, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.HandlersOffset, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.LiteralsCount, err = chunk.Reader.ReadUInt16(endian)
	if err != nil {
		return err
	}

	chunk.LiteralsOffset, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err

	}

	chunk.LiteralsDataCount, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	chunk.LiteralsDataOffset, err = chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}

	return nil

}

func ReadLscrChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*LscrChunk, error) {
	var err error
	chunk := &LscrChunk{
		Reader: r,
	}

	chunk.Reader.HexDump(true)

	r.Seek(0, 0)
	err = chunk.Read(endian)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

func (c *LscrChunk) Print() {
	fmt.Println()
	fmt.Printf("-- Lscr*\n")
	fmt.Printf("--\n")
	fmt.Println()

}

func (c *LscrChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
