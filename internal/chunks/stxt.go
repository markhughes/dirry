package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/palettes"
	"github.com/markhughes/dirry/internal/utils"
)

type Formatting struct {
	Offset       int32
	Height       int16
	Ascent       int16
	FontId       int16
	FontName     string
	FontPlatform uint16
	Bold         bool
	Italic       bool
	Underline    bool
	Padding      int
	FontSize     int16
	Colour       palettes.Pixel24
}

type StyledTextChunk struct {
	Reader *binary_reader.BinaryReader

	Text        string
	Formattings []Formatting

	Unk1 int32
	Unk2 int32
}

func ReadStxtChunkRaw(r *binary_reader.BinaryReader, fonts map[uint32]*Font, endian binary.ByteOrder, isAfterburner bool) (*StyledTextChunk, error) {

	var err error
	chunk := &StyledTextChunk{
		Reader: r,
	}

	chunk.Reader.HexDump(true)

	// Always big endian?
	r.Seek(0, 0)
	err = chunk.Read(binary.BigEndian, fonts)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

func (chunk *StyledTextChunk) Read(endian binary.ByteOrder, fonts map[uint32]*Font) error {
	var err error

	chunk.Unk1, err = chunk.Reader.ReadInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading unk1: %s", err)
	}

	// fmt.Printf("unk1: %d\n", chunk.Unk1)

	textLength, err := chunk.Reader.ReadInt32(endian)
	if err != nil {
		return fmt.Errorf("error reading text length: %s", err)
	}
	utils.DebugMsg("stxt", "textLength: %d (%d)", textLength, int(textLength))

	chunk.Reader.ReadInt32(endian)

	chunk.Text, err = chunk.Reader.ReadString(int(textLength), endian)
	if err != nil {
		return fmt.Errorf("error reading text: %s", err)
	}

	formatcount, err := chunk.Reader.ReadInt16(endian)
	if err != nil {
		return fmt.Errorf("error reading format count: %s", err)
	}

	chunk.Formattings = make([]Formatting, 0)

	for i := 0; i < int(formatcount); i++ {
		format := Formatting{}

		format.Offset, err = chunk.Reader.ReadInt32(endian)
		if err != nil {
			return fmt.Errorf("error reading format offset: %s", err)
		}

		format.Height, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format height: %s", err)
		}

		format.Ascent, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format ascent: %s", err)
		}

		format.FontId, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format font id: %s", err)
		}

		// check if font is in fonts
		if _, ok := fonts[uint32(format.FontId)]; ok {

			format.FontName = fonts[uint32(format.FontId)].Name
			format.FontPlatform = fonts[uint32(format.FontId)].Platform
		} else {
			utils.ErrorMsg("stxt", "font id %d not found in fonts", format.FontId)
			utils.DebugMsg("stxt", "available fonts: %v", fonts)
		}
		formatting, err := chunk.Reader.ReadUByte()
		if err != nil {
			return fmt.Errorf("error reading format formatting: %s", err)
		}

		format.Bold = formatting&1 != 0
		format.Italic = formatting&2 != 0
		format.Underline = formatting&4 != 0

		format.Padding, err = chunk.Reader.ReadUByte()
		if err != nil {
			return fmt.Errorf("error reading format padding: %s", err)
		}

		format.FontSize, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format font size: %s", err)
		}

		red, err := chunk.Reader.ReadUInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format red: %s", err)
		}

		green, err := chunk.Reader.ReadUInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format green: %s", err)
		}

		blue, err := chunk.Reader.ReadUInt16(endian)
		if err != nil {
			return fmt.Errorf("error reading format blue: %s", err)
		}

		format.Colour = palettes.Pixel24{
			R: uint8(red >> 8),
			G: uint8(green >> 8),
			B: uint8(blue >> 8),
		}

		chunk.Formattings = append(chunk.Formattings, format)
	}

	return nil
}

func (c *StyledTextChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %s", err)
	}

	return string(bytes), nil
}
