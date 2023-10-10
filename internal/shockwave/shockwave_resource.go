package shockwave

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/binary_reader"
)

type ShockwaveResource struct {
	ResourceId int32

	// -1 = if using afterburner, it's from the ils
	Offset           int32
	CompressedSize   int32
	UncompressedSize int32

	// 0 = compressed
	// 1 = uncompressed
	CompressionType int32

	CastId    int32
	LibId     int32
	ChunkType string
	Name      string
	Children  []ShockwaveResource
	Binary    []byte
}

func (resource *ShockwaveResource) DumpBinary(outputFolder string) error {
	if string(resource.ChunkType) == "" {
		resource.ChunkType = "unknown"
	}
	var size = resource.UncompressedSize

	outputFolder = filepath.Join(outputFolder, resource.ChunkType)

	err := os.MkdirAll(outputFolder, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return (err)
		}
	}

	var fileName = filepath.Join(outputFolder, fmt.Sprint(resource.ResourceId)+"_"+fmt.Sprint(size)+".bin")
	file, err := os.Create(fileName)
	if err != nil {
		return (err)
	}

	// fmt.Printf("Dumping %d bytes to %s\n", len(resource.Binary), fileName)
	// fmt.Printf("Directory %s\n", outputFolder)
	defer file.Close()

	_, err = file.Write(resource.Binary)
	if err != nil {
		return (err)
	}

	return nil

}

func (resource *ShockwaveResource) GetReader() (*binary_reader.BinaryReader, error) {
	return binary_reader.NewBinaryReader(resource.Binary, resource.UncompressedSize)

}
