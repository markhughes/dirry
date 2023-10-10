package pfr

import (
	"encoding/binary"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type LogicalFontRecord struct {
	fontMatrix                  [4]int32
	zero                        uint8
	extraItemsPresent           uint8
	twoByteBoldThicknessValue   uint8
	boldFlag                    uint8
	twoByteStrokeThicknessValue uint8
	strokeFlag                  uint8
	lineJoinType                uint8
	strokeThickness             int16
	miterLimit                  int32
	boldThickness               int16
	nExtraItems                 uint8
	extraItemSize               uint8
	extraItemType               uint8
	extraItemData               []uint8
	physFontSize                uint16
	physFontOffset              uint32
	physFontSizeIncrement       uint8
}

type LogicalFontSection struct {
	LogFonts []LogicalFontRecord
}

func (chunk *PfrFont) parseLfs(offset uint32, reader *binary_reader.BinaryReader) (err error) {
	reader.Seek(int64(offset), 0)

	section := LogicalFontSection{}
	section.LogFonts = make([]LogicalFontRecord, 0)

	for i := 0; i < int(chunk.LogicalFontDirectory.NLogFonts); i++ {
		record := LogicalFontRecord{}
		record.fontMatrix[0], err = reader.ReadInt32(binary.BigEndian)
		if err != nil {
			return err
		}
		record.fontMatrix[1], err = reader.ReadInt32(binary.BigEndian)
		if err != nil {
			return err
		}

		record.fontMatrix[2], err = reader.ReadInt32(binary.BigEndian)
		if err != nil {
			return err
		}

		record.fontMatrix[3], err = reader.ReadInt32(binary.BigEndian)
		if err != nil {
			return err
		}

		record.zero, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.extraItemsPresent, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.twoByteBoldThicknessValue, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.boldFlag, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.twoByteStrokeThicknessValue, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.strokeFlag, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.lineJoinType, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.strokeThickness, err = reader.ReadInt16(binary.BigEndian)
		if err != nil {
			return err
		}

		record.miterLimit, err = reader.ReadInt32(binary.BigEndian)
		if err != nil {
			return err
		}

		record.boldThickness, err = reader.ReadInt16(binary.BigEndian)
		if err != nil {
			return err
		}

		record.nExtraItems, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.extraItemSize, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.extraItemType, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		record.extraItemData, err = reader.ReadBytes(int(record.extraItemSize))
		if err != nil {
			return err
		}

		record.physFontSize, err = reader.ReadUInt16(binary.BigEndian)
		if err != nil {
			return err
		}

		record.physFontOffset, err = reader.ReadUInt32(binary.BigEndian)
		if err != nil {
			return err
		}

		record.physFontSizeIncrement, err = reader.ReadUInt8()
		if err != nil {
			return err
		}

		section.LogFonts = append(section.LogFonts, record)
	}

	chunk.LogicalFontSection = section
	return nil
}
