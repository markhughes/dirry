package chunks

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	"github.com/markhughes/dirry/internal/consts"
)

type EdimChunk struct {
	Chunk     *BinaryChunk
	Extension string
	MIME      string
	Binary    []byte
}

func ReadEdimChunkRaw(r *bytes.Reader, length int, endian binary.ByteOrder) (*EdimChunk, error) {

	var err error
	chunk := &EdimChunk{}
	chunk.Chunk, err = FromBinaryAtPartHeadless(r, "ediM", endian)
	chunk.Chunk.Length = int32(length)
	if err != nil {
		return nil, err
	}

	_, err = processEdimChunk(chunk, endian)
	if err != nil {
		return nil, err
	}

	return chunk, nil

}

func processEdimChunk(chunk *EdimChunk, endian binary.ByteOrder) (*EdimChunk, error) {
	var err error

	// Read the binary data into chunk.Binary
	chunk.Binary, err = chunk.Chunk.ReadAllBytes()
	if err != nil {
		return nil, err
	}

	if len(chunk.Binary) < 4 {
		return nil, fmt.Errorf("binary data is too short for a magic number")
	}

	kind, _ := filetype.Match(chunk.Binary)
	chunk.Extension = kind.Extension
	chunk.MIME = kind.MIME.Value

	return chunk, nil
}

func (c *EdimChunk) Print() {
	fmt.Printf("-- BITd\n")
}

func (c *EdimChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c *EdimChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "ediM")
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "ediM")
	}

	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists
	outputFile := filepath.Join(outputFolder, name+"."+c.Extension)

	targetFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	_, err = targetFile.Write(c.Binary)
	if err != nil {
		panic(err)
	}

	err = targetFile.Sync()
	if err != nil {
		panic(err)
	}

}
