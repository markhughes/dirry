package pfr

import (
	"encoding/binary"
	"io"

	"github.com/markhughes/dirry/internal/binary_reader"
)

func (chunk *PfrFont) parseLfd(offset uint16, reader *binary_reader.BinaryReader) (err error) {
	reader.Seek(int64(offset), io.SeekStart)

	// Read the nLogFonts value
	nLogFonts, err := reader.ReadUInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.LogicalFontDirectory = LogicalFontDirectory{
		NLogFonts: nLogFonts,
		Fonts:     make([]LogicalFont, nLogFonts),
	}

	for i := uint16(0); i < nLogFonts; i++ {
		chunk.LogicalFontDirectory.Fonts[i].LogFontSize, err = reader.ReadUInt16(binary.BigEndian)
		if err != nil {
			return err
		}

		// Read 24-bit LogFontOffset manually
		offsetBytes, err := reader.ReadBytes(3)
		if err != nil {
			return err
		}

		chunk.LogicalFontDirectory.Fonts[i].LogFontOffset = uint32(offsetBytes[0])<<16 | uint32(offsetBytes[1])<<8 | uint32(offsetBytes[2])
	}

	return nil
}
