package chunks

import (
	"bufio"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/utils"
)

type AfterburnerResource struct {
	ResourceId         uint32
	Offset             int32
	CompressedLength   uint32
	DecompressedLength uint32
	CompressionType    uint32
	ChunkType          string
	Data               []byte
}

type ABMPChunk struct {
	Chunk     *BinaryChunk
	Resources []*AfterburnerResource
}

func ReadABMPChunk(r io.ReadSeeker, endian binary.ByteOrder) (*ABMPChunk, error) {
	var err error

	chunk := &ABMPChunk{}

	chunk.Chunk, err = FromBinary(r, "ABMP", endian)
	if err != nil {

		return nil, err
	}

	start, err := chunk.Chunk.r.Seek(0, io.SeekCurrent) // Get current file pointer position
	if err != nil {
		return chunk, err
	}

	compressionType, err := chunk.Chunk.ReadVarInt()
	if err != nil {
		return chunk, fmt.Errorf("error reading compression type: %s", err)
	}

	uncompressedLen, err := chunk.Chunk.ReadVarInt()
	if err != nil {
		return chunk, fmt.Errorf("error reading uncompressed length: %s", err)
	}

	utils.DebugMsg("ABMP", "compressionType: %d\n", compressionType)
	utils.DebugMsg("ABMP", "uncompressedLen: %d\n", uncompressedLen)

	zlibReader, err := zlib.NewReader(bufio.NewReader(chunk.Chunk.r))
	if err != nil {
		return nil, fmt.Errorf("error creating zlib reader: %s", err)
	}
	defer zlibReader.Close()

	limitedReader := &io.LimitedReader{R: zlibReader, N: int64(uncompressedLen)}

	decompressedData, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("error reading decompressed data: %s", err)
	}

	chunk.Chunk.Data = decompressedData

	reader, err := binary_reader.NewBinaryReader(decompressedData, int32(len(decompressedData)))
	if err != nil {
		return nil, fmt.Errorf("error creating binary reader: %s", err)
	}

	reader.ReadVarInt()
	reader.ReadVarInt()

	entries, _, err := reader.ReadVarInt()
	if err != nil {
		return chunk, fmt.Errorf("error reading entries: %s", err)
	}

	var i = 0
	for i < int(entries) {
		resourceId, _, err := reader.ReadVarInt()
		if err != nil {
			return chunk, fmt.Errorf("error reading entry id: %s", err)
		}

		offset, _, err := reader.ReadVarInt()
		if err != nil {
			return chunk, fmt.Errorf("error reading entry offset: %s", err)
		}

		compressedLength, _, err := reader.ReadVarInt()
		if err != nil {
			return chunk, fmt.Errorf("error reading entry compressedLength: %s", err)
		}

		decompressedLength, _, err := reader.ReadVarInt()
		if err != nil {
			return chunk, fmt.Errorf("error reading entry decompressed length: %s", err)
		}

		compressionType, _, err := reader.ReadVarInt()
		if err != nil {
			return chunk, fmt.Errorf("error reading entry compression type: %s", err)
		}

		// Read chunk type, and convert to string.
		chunkTypeBytes, err := reader.ReadBytes(4)
		if err != nil {
			return nil, err
		}

		var chunkType string
		if endian == binary.BigEndian {
			chunkType = string(chunkTypeBytes[:])
		} else {
			chunkType = string(utils.ReverseBytes(chunkTypeBytes[:]))
		}

		var res AfterburnerResource = AfterburnerResource{
			ResourceId:         resourceId,
			Offset:             int32(offset),
			CompressedLength:   compressedLength,
			DecompressedLength: decompressedLength,
			CompressionType:    compressionType,
			ChunkType:          chunkType,
		}

		chunk.Resources = append(chunk.Resources, &res)
		i++

		position, _ := reader.Seek(0, io.SeekCurrent)

		utils.DebugMsg("ABMP", "@%d  resourceId: %d, offset: %d, compressedLength: %d, decompressedLength: %d, compressionType: %d, chunkType: %s\n", position, resourceId, offset, compressedLength, decompressedLength, compressionType, chunkType)

	}

	_, _ = chunk.Chunk.r.Seek(start+int64(chunk.Chunk.Length), io.SeekStart)

	return chunk, nil

}

func (c *ABMPChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (m *ABMPChunk) FindResourceByType(chunkType string) (*AfterburnerResource, error) {
	for _, res := range m.Resources {
		if res.ChunkType == chunkType {
			return res, nil
		}
	}

	return nil, fmt.Errorf("resource %s not found", chunkType)
}

func (m *ABMPChunk) FindResourceByID(resourceId int) (*AfterburnerResource, error) {
	for i := range m.Resources {
		var res = m.Resources[i]
		if int(res.ResourceId) == resourceId {
			return res, nil
		}
	}

	return nil, fmt.Errorf("resource %d not found", resourceId)
}

func (m *ABMPChunk) FindResourcesByType(chunkType string) ([]*AfterburnerResource, error) {
	var resources []*AfterburnerResource

	for i := range m.Resources {
		var res = m.Resources[i]
		if (res.ChunkType) == chunkType {
			resources = append(resources, res)
		}
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("resource %s not found", chunkType)
	} else {
		return resources, nil
	}
}
