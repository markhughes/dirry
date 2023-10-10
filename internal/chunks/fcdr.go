package chunks

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type CompressionType struct {
	GUID string
	Name string
}

type FcdrChunk struct {
	Chunk            *BinaryChunk
	Count            uint16
	CompressionTypes []CompressionType
}

func ReadFcdrChunk(r io.ReadSeeker, endian binary.ByteOrder) (*FcdrChunk, error) {
	var err error

	chunk := &FcdrChunk{}

	chunk.Chunk, err = FromBinary(r, "Fcdr", endian)
	if err != nil {

		return nil, err
	}

	reader := bufio.NewReader(chunk.Chunk.r)
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	// Set the limit on the zlib reader to the original uncompressed data length
	limitedReader := &io.LimitedReader{R: zlibReader, N: int64(chunk.Chunk.Length * 10)}

	// Read all data from the limited reader (this will be the decompressed data)
	decompressedData, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	// add reader + decompressed data to chunk
	chunk.Chunk.Data = decompressedData
	chunk.Chunk.BytesReader = bytes.NewReader(chunk.Chunk.Data)

	err = binary.Read(chunk.Chunk.BytesReader, endian, &chunk.Count)
	if err != nil {
		return nil, err
	}

	var guids = make([]string, 0)
	for i := 0; i < int(chunk.Count); i++ {

		// guid

		var data1 uint32
		var data2 uint16
		var data3 uint16

		err = binary.Read(chunk.Chunk.BytesReader, endian, &data1)
		if err != nil {
			return nil, err
		}

		err = binary.Read(chunk.Chunk.BytesReader, endian, &data2)
		if err != nil {
			return nil, err
		}

		err = binary.Read(chunk.Chunk.BytesReader, endian, &data3)
		if err != nil {
			return nil, err
		}

		var data4 [8]uint8 = [8]uint8{}
		for i := 0; i < 8; i++ {
			err = binary.Read(chunk.Chunk.BytesReader, endian, &data4[i])
			if err != nil {
				return nil, err
			}
		}

		var guid = fmt.Sprintf("%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X", data1, data2, data3, data4[0], data4[1], data4[2], data4[3], data4[4], data4[5], data4[6], data4[7])
		guids = append(guids, guid)

	}

	for i := 0; i < int(chunk.Count); i++ {
		compressionDesc, err := chunk.Chunk.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}

		var compressionType = CompressionType{
			GUID: guids[i],
			Name: compressionDesc,
		}

		chunk.CompressionTypes = append(chunk.CompressionTypes, compressionType)
	}

	// go to end
	chunk.Chunk.r.Seek(chunk.Chunk.StartPosition+int64(chunk.Chunk.Length), io.SeekStart)

	return chunk, nil

}

func (c *FcdrChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
