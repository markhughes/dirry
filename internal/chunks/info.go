package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/utils"
)

type InfoChunk struct {
	Reader *binary_reader.BinaryReader

	unknown01 int16
	Version   int16

	Rect utils.Rect

	Width  int16
	Height int16

	CastListStart int16
	CastListEnd   int16

	InitialFrameRate int8

	BgColor int8

	LightSwitch int8

	CommentFont  int16
	CommentSize  int16
	CommentStyle int16
}

func (chunk *InfoChunk) Read(endian binary.ByteOrder) error {
	var err error

	chunk.unknown01, err = chunk.Reader.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Version, err = chunk.Reader.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Rect, err = chunk.Reader.ReadRect(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.Width = chunk.Rect.Right - chunk.Rect.Left
	chunk.Height = chunk.Rect.Bottom - chunk.Rect.Top

	chunk.CastListStart, err = chunk.Reader.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.CastListEnd, err = chunk.Reader.ReadInt16(binary.BigEndian)
	if err != nil {
		return err
	}

	chunk.InitialFrameRate, err = chunk.Reader.ReadInt8()
	if err != nil {
		return err
	}

	chunk.LightSwitch, err = chunk.Reader.ReadInt8()
	if err != nil {
		return err
	}

	// TODO: from here on is different between versions and it's a bit irrelevant for extraction purposes rn, so skipping for now
	//       would still like to get the accurate structure though

	chunk.Reader.ReadInt16(endian) // ?

	chunk.CommentFont, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.CommentSize, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.CommentStyle, err = chunk.Reader.ReadInt16(endian)
	if err != nil {
		return err
	}

	chunk.BgColor, err = chunk.Reader.ReadInt8()
	if err != nil {
		return err
	}

	return nil

}

func ReadInfoChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*InfoChunk, error) {
	var err error
	chunk := &InfoChunk{
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

func (c *InfoChunk) Print() {
	fmt.Println()
	fmt.Printf("-- VWCH*\n")
	fmt.Printf("Version .......... %d\n", c.Version)
	fmt.Printf("Rect: ..............\n")
	fmt.Printf("  Top: ............. %v\n", c.Rect.Top)
	fmt.Printf("  Bottom: .......... %v\n", c.Rect.Bottom)
	fmt.Printf("  Left: ............ %v\n", c.Rect.Left)
	fmt.Printf("  Bottom: .......... %v\n", c.Rect.Right)
	fmt.Printf("  Width ............ %d\n", c.Width)
	fmt.Printf("  Height ........... %d\n", c.Height)
	fmt.Printf("Cast List Start .. %d\n", c.CastListStart)
	fmt.Printf("Cast List End .... %d\n", c.CastListEnd)
	fmt.Printf("Init Frame Rate .. %d\n", c.InitialFrameRate)
	fmt.Printf("Background Color . %d\n", c.BgColor)
	fmt.Printf("--\n")
	fmt.Println()

}

func (c *InfoChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
