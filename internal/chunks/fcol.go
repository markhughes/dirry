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

type FcolChunk struct {
	Reader  *binary_reader.BinaryReader
	Colours []palettes.Pixel24
}

func (chunk *FcolChunk) Read(endian binary.ByteOrder) error {

	chunk.Reader.ReadInt32(endian)
	chunk.Reader.ReadInt32(endian)

	chunk.Colours = make([]palettes.Pixel24, 0)
	// there are 16 colours
	for i := 0; i < 16; i++ {
		red, err := chunk.Reader.ReadInt32(endian)
		if err != nil {
			if i > 0 {
				break
			}
			return fmt.Errorf("failed to read fcol chunk red: %s", err)
		}

		green, err := chunk.Reader.ReadInt32(endian)
		if err != nil {
			return fmt.Errorf("failed to read fcol chunk green: %s", err)
		}

		blue, err := chunk.Reader.ReadInt32(endian)
		if err != nil {
			return fmt.Errorf("failed to read fcol chunk blue: %s", err)
		}

		chunk.Colours = append(chunk.Colours, palettes.Pixel24{
			R: uint8(red >> 8),
			G: uint8(green >> 8),
			B: uint8(blue >> 8),
		})

		utils.DebugMsg("fcol", "Colour %d: %d, %d, %d\n", i, red, green, blue)
	}

	return nil
}

func ReadFcolChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*FcolChunk, error) {
	var err error
	chunk := &FcolChunk{
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

func (c *FcolChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "FCOL")
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "FCOL")
	}

	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists
	outputFile := filepath.Join(outputFolder, name+".act")

	targetFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	// Initialize the ACT file with zeroes. The file should be 772 bytes long.
	act := make([]byte, 772)
	for i := range act {
		act[i] = 0
	}

	// Add the colors to the ACT file.
	for i, color := range c.Colours {
		act[i*3] = color.R
		act[i*3+1] = color.G
		act[i*3+2] = color.B
	}

	// Add the palette size to the ACT file.
	binary.BigEndian.PutUint16(act[768:], uint16(len(c.Colours)))

	// Write the ACT file to disk.
	_, err = targetFile.Write(act)
	if err != nil {
		panic(err)
	}

	err = targetFile.Sync()
	if err != nil {
		panic(err)
	}
}

func (c *FcolChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
