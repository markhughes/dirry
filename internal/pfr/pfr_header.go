package pfr

import (
	"encoding/binary"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type PFRHeader struct {
	PfrHeaderSig            string
	PfrVersion              uint16
	PfrHeaderSig2           string
	PfrHeaderSize           uint16
	LogFontDirSize          uint16
	LogFontDirOffset        uint16
	LogFontMaxSize          uint16
	LogFontSectionSize      uint32 // 24 bits
	LogFontSectionOffset    uint32 // 24 bits
	PhysFontMaxSize         uint16
	PhysFontSectionSize     uint32 // 24 bits
	PhysFontSectionOffset   uint32 // 24 bits
	GpsMaxSize              uint16
	GpsSectionSize          uint32 // 24 bits
	GpsSectionOffset        uint32 // 24 bits
	MaxBlueValues           uint8
	MaxXorus                uint8
	MaxYorus                uint8
	PhysFontMaxSizeHighByte uint8
	Zeros                   uint8  // 6 bits
	PfrInvertBitmap         uint8  // 1 bit
	PfrBlackPixel           uint8  // 1 bit
	BctMaxSize              uint32 // 24 bits
	BctSetMaxSize           uint32 // 24 bits
	PftBctSetMaxSize        uint32 // 24 bits
	NPhysFonts              uint16
	MaxStemSnapVsize        uint8
	MaxStemSnapHsize        uint8
	MaxChars                uint16
}

func (chunk *PfrFont) parseHeader(reader *binary_reader.BinaryReader) (err error) {
	chunk.Header.PfrHeaderSig, err = reader.ReadString(4, binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PfrVersion, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PfrHeaderSig2, err = reader.ReadString(2, binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PfrHeaderSize, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.LogFontDirSize, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.LogFontDirOffset, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.LogFontMaxSize, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.LogFontSectionSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.LogFontSectionOffset, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PhysFontMaxSize, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PhysFontSectionSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PhysFontSectionOffset, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.GpsMaxSize, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.GpsSectionSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.GpsSectionOffset, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.MaxBlueValues, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.MaxXorus, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.MaxYorus, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.PhysFontMaxSizeHighByte, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.Zeros, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.PfrInvertBitmap, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.PfrBlackPixel, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.BctMaxSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.BctSetMaxSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.PftBctSetMaxSize, err = reader.ReadUInt32(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.NPhysFonts, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Header.MaxStemSnapVsize, err = reader.ReadUInt8()
	if err != nil {
		return err
	}

	chunk.Header.MaxStemSnapHsize, err = reader.ReadUInt8()

	if err != nil {
		return err
	}

	chunk.Header.MaxChars, err = reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	return nil
}
