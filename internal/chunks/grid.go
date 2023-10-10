package chunks

import (
	"encoding/binary"
	"encoding/json"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/utils"
)

// Axiss can be "horizontal" or "vertical
type GridAxis string

const (
	HorizontalAxis GridAxis = "horizontal"
	VerticalAxis   GridAxis = "vertical"
)

// Axiss can be "horizontal" or "vertical
type GuideDisplay string

const (
	Dots  GuideDisplay = "dots"
	Lines GuideDisplay = "lines"
)

type Guide struct {
	Axis     GridAxis
	Position int16
}

type GridChunk struct {
	Reader *binary_reader.BinaryReader

	Width  int16
	Height int16

	Display GuideDisplay

	GridColour int16

	Guides       []Guide
	GuidesColour int16
}

func ReadGridChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*GridChunk, error) {

	var err error
	chunk := &GridChunk{
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

func (chunk *GridChunk) Read(endian binary.ByteOrder) error {
	var err error

	_, err = chunk.Reader.ReadBytes(4)
	if err != nil {
		return err
	}

	chunk.Width, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.Height, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	display, err := chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	if display == 2 {
		chunk.Display = Dots
	} else {
		chunk.Display = Lines
	}

	chunk.GridColour, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	guideCount, err := chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.GuidesColour, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.Guides = make([]Guide, guideCount)

	utils.DebugMsg("grid", "guideCount = %d\n", guideCount)
	utils.DebugMsg("grid", "guidColour = %d\n", chunk.GuidesColour)

	for i := 0; i < int(guideCount); i++ {
		guide := Guide{}

		axis, err := chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

		if axis == 1 {
			guide.Axis = VerticalAxis
		} else {
			guide.Axis = HorizontalAxis
		}

		guide.Position, err = chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

		utils.DebugMsg("grid", "guide %d: %s %d\n", i, guide.Axis, guide.Position)

		chunk.Guides[i] = guide

	}

	return nil
}

func (c *GridChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
