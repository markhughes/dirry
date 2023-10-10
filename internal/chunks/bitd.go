package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/bitd"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/members"
	"github.com/markhughes/dirry/internal/palettes"
	"github.com/markhughes/dirry/internal/utils"
)

type BitmapChunk struct {
	Reader *binary_reader.BinaryReader

	Encoded bool
	Member  *members.MemberBitmap
	Data    []byte
	Info    bitd.BitmapInfo
}

func (chunk *BitmapChunk) Read(endian binary.ByteOrder, castChunk *CastChunk) error {
	var err error

	if member, ok := castChunk.Member.(*members.MemberBitmap); ok {
		chunk.Member = member
	} else {
		return fmt.Errorf("castChunk.Member is not a BitmapCastMember")
	}

	var clut palettes.PaletteValue
	clut, err = palettes.RetrievePallete(palettes.Clut(chunk.Member.Clut))
	if err != nil {
		return err
	}
	// TODO: handle gray bw etc?
	// 	clut, err = palettes.RetrievePallete(palettes.ClutGrayscale)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	utils.DebugMsg("BITD", "clut.Size: %v\n", clut.Size)

	// content, _ := chunk.Chunk.ReadBytes(int(chunk.Chunk.Length))

	width := chunk.Member.InitialRect.Width
	height := chunk.Member.InitialRect.Height

	if chunk.Member.BitsPerPixel == 1 {
		// 1-bit images' width is multiple of 16
		width = ((width-1)/16 + 1) * 16
	} else if chunk.Member.BitsPerPixel == 4 {
		// 4-bit images' width is multiple of 4
		width = ((width-1)/4 + 1) * 4
	}

	utils.DebugMsg("BITD", "width: %d, height: %d\n", width, height)

	if chunk.Reader.Length == int32(height*width*int16(chunk.Member.BitsPerPixel)/8) {
		utils.DebugMsg("BITD", "chunk.Chunk.Length == %v and that matches height&width*bitsPerPixel/8\n", chunk.Reader.Length)
	} else {
		utils.DebugMsg("BITD", "length: %d\n", chunk.Reader.Length)

	}

	data, err := chunk.Reader.ReadBytes(int(chunk.Reader.Length))
	if err != nil {
		return err
	}

	chunk.Info, chunk.Data, err = bitd.ConvertImage(data, int(width), int(height), int(chunk.Member.BitsPerPixel), clut)
	if err != nil {
		return err
	}
	return nil

}

func ReadBitmapChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, castChunk *CastChunk, isAfterburner bool) (*BitmapChunk, error) {
	var err error

	chunk := &BitmapChunk{
		Reader: r,
	}

	chunk.Reader.HexDump(true)

	r.Seek(0, 0)
	err = chunk.Read(binary.BigEndian, castChunk)
	if err != nil {
		return nil, err
	}

	return chunk, nil
}

func (c *BitmapChunk) Print() {
	fmt.Printf("-- BITd\n")
}

func (c *BitmapChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c *BitmapChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "BITD")
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "BITD")
	}

	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists
	outputFile := filepath.Join(outputFolder, name+".png")

	targetFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	// Write the BMP file to disk.
	_, err = targetFile.Write(c.Data)
	if err != nil {
		panic(err)
	}

	err = targetFile.Sync()
	if err != nil {
		panic(err)
	}

	utils.DebugMsg("BITD", "Saved to %s\n", outputFile)
}
