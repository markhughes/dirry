package pfr

import (
	"github.com/markhughes/dirry/internal/binary_reader"
)

type LogicalFont struct {
	LogFontSize   uint16
	LogFontOffset uint32 // 24 bits
}

type LogicalFontDirectory struct {
	NLogFonts uint16
	Fonts     []LogicalFont
}

type PfrFont struct {
	Header               PFRHeader
	LogicalFontDirectory LogicalFontDirectory
	LogicalFontSection   LogicalFontSection
}

func (chunk *PfrFont) Parse(data []byte) error {
	// TODO: find or work on a seperate pfr convert utility
	// https://web.archive.org/web/20040720143852/http://www.bitstream.com:80/categories/developer/truedoc/pfrspec1.2.pdf
	reader, err := binary_reader.NewBinaryReader(data, -1)
	if err != nil {
		return err
	}

	err = chunk.parseHeader(reader)
	if err != nil {
		return err
	}

	err = chunk.parseLfd(chunk.Header.LogFontDirOffset, reader)
	if err != nil {
		return err
	}

	err = chunk.parseLfs(chunk.Header.LogFontSectionOffset, reader)
	if err != nil {
		return err
	}

	return nil
}
