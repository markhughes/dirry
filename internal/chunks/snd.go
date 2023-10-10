package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/consts"
)

type SndChunk struct {
	Chunk *BinaryChunk
	Type  int16
}

func ReadSndChunk(r *os.File, endian binary.ByteOrder, offset int64) (*SndChunk, error) {
	_, err := r.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	chunk := &SndChunk{}

	chunk.Chunk, err = FromBinaryAt(r, "snd ", endian, offset)
	if err != nil {
		return nil, err
	}

	err = binary.Read(chunk.Chunk.r, binary.LittleEndian, &chunk.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file type: %v", err)
	}

	return chunk, nil
}

func (c *SndChunk) Save(projectName string, name string) {
	outputFolder := filepath.Join(consts.PathDump, projectName, "converted", "snd")
	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists
	outputFile := filepath.Join(outputFolder, name+".act")

	targetFile, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	snd := make([]byte, 0)

	// Write the snd file to disk.
	_, err = targetFile.Write(snd)
	if err != nil {
		panic(err)
	}

	err = targetFile.Sync()
	if err != nil {
		panic(err)
	}
}

func (c *SndChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
