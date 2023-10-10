package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/palettes"
	"github.com/markhughes/dirry/internal/utils"
)

type ClutChunk struct {
	Reader  *binary_reader.BinaryReader
	Palette palettes.PaletteValue
}

func (c *ClutChunk) Read(endian binary.ByteOrder) error {
	var err error

	c.Palette.Size, err = c.Reader.ReadInt32(endian)
	if err != nil {
		return err
	}

	utils.DebugMsg("CLUT", fmt.Sprintf("Palette size: %d", c.Palette.Size))

	if c.Palette.Size > 256 {
		utils.InfoMsg("CLUT", fmt.Sprintf("Palette size is bigger than %d: %d", 256, c.Palette.Size))
		c.Palette.Size = 256
	}

	for i := 0; i < int(c.Palette.Size); i++ {
		red, err := c.Reader.ReadUInt16(endian)
		if err != nil {
			return err
		}
		green, err := c.Reader.ReadUInt16(endian)
		if err != nil {
			return err
		}
		blue, err := c.Reader.ReadUInt16(endian)
		if err != nil {
			return err
		}

		c.Palette.Palette[i] = palettes.Pixel24{
			R: uint8(red >> 8),
			G: uint8(green >> 8),
			B: uint8(blue >> 8),
		}

	}

	return nil
}

func ReadClutChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*ClutChunk, error) {
	var err error

	chunk := &ClutChunk{
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

func (c *ClutChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c *ClutChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "CLUT")
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "CLUT")
	}

	os.MkdirAll(outputFolder, os.ModePerm)
	outputFile := filepath.Join(outputFolder, name+".act")

	targetFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	// Initialize the ACT file with zeroes, the file should be 772 bytes long.
	act := make([]byte, 772)
	for i := range act {
		act[i] = 0
	}

	for i, color := range c.Palette.Palette {
		act[i*3] = color.R
		act[i*3+1] = color.G
		act[i*3+2] = color.B
	}

	binary.BigEndian.PutUint16(act[768:], uint16(c.Palette.Size))

	_, err = targetFile.Write(act)
	if err != nil {
		panic(err)
	}

	err = targetFile.Sync()
	if err != nil {
		panic(err)
	}
}
