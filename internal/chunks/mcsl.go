package chunks

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"strings"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/utils"
)

type CastLib struct {
	Name        string
	Path        string
	ItemCount   int16
	NumId       int16
	StorageType int
	LibId       int
}

type MCsLChunk struct {
	Reader *binary_reader.BinaryReader

	Count    uint32
	UnkCount uint32
	CastLibs []CastLib
}

func (chunk *MCsLChunk) Read(endian binary.ByteOrder) error {
	chunk.Reader.Seek(4, io.SeekCurrent)

	count, err := chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}
	chunk.Count = count
	utils.DebugMsg("mcsl", "count: %d\n", count)

	if _, err := chunk.Reader.Seek(2, io.SeekCurrent); err != nil {
		return err
	}

	unkCount, err := chunk.Reader.ReadUInt32(endian)
	if err != nil {
		return err
	}
	chunk.UnkCount = unkCount + 1 // add 1

	// Ignore 4*UnkCount bytes
	if _, err := chunk.Reader.Seek(4*int64(chunk.UnkCount), io.SeekCurrent); err != nil {
		return err
	}

	for i := 0; i < int(chunk.Count); i++ {
		pos, _ := chunk.Reader.Seek(0, io.SeekCurrent)
		utils.DebugMsg("mcsl", "Reading castlib %d @ %d\n", i, pos)

		// Read nameSize
		nameSize, err := chunk.Reader.ReadUInt8()
		if err != nil {
			return err
		}
		utils.DebugMsg("mcsl", "nameSize: %d\n", nameSize)

		var name string
		if nameSize > 0 {
			// Read name
			nameBytes, err := chunk.Reader.ReadBytes(int(nameSize))
			if err != nil {
				return err
			}

			name = string(nameBytes)

			// remove trailing null characters if any
			name = strings.TrimRight(name, "\x00")
		}
		utils.DebugMsg("mcsl", "name: %s\n", name)

		// Ignore 1 byte
		if _, err := chunk.Reader.Seek(1, io.SeekCurrent); err != nil {
			return err
		}

		// Read pathSize
		pathSize, err := chunk.Reader.ReadUInt8()
		if err != nil {
			return err
		}
		utils.DebugMsg("mcsl", "pathSize: %d\n", pathSize)

		var path string
		if pathSize > 1 {
			// Read path
			pathBytes, err := chunk.Reader.ReadBytes(int(pathSize) + 1) // Add 1 to also read the trailing null byte
			if err != nil {
				return err
			}
			path = string(pathBytes)

			// remove trailing null characters if any
			path = strings.TrimRight(path, "\x00")
		} else if pathSize == 1 {
			// If pathSize is 1, we skip 1 byte
			if _, err := chunk.Reader.Seek(1, io.SeekCurrent); err != nil {
				return err
			}
		}

		utils.DebugMsg("mcsl", "path: %s\n", path)

		// unknown purpose of this
		if _, err := chunk.Reader.Seek(1, io.SeekCurrent); err != nil {
			return err
		}

		// unknown purpose of this
		chunk.Reader.ReadInt8()

		// unknown purpose of this
		chunk.Reader.ReadInt8()

		// itemCount
		itemCount, err := chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}
		utils.DebugMsg("mcsl", "itemCount: %d\n", itemCount)
		// Ignore 2 bytes?

		numId, err := chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}

		// Read libId
		libId, err := chunk.Reader.ReadInt16(endian)
		if err != nil {
			return err
		}
		libId -= 1023

		var storageType = 1
		if path != "" {
			storageType = 0
		}

		chunk.CastLibs = append(chunk.CastLibs, CastLib{Name: name, Path: path, ItemCount: itemCount, LibId: int(libId), StorageType: storageType, NumId: (numId)})
	}

	return nil
}

func ReadMcslChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*MCsLChunk, error) {
	var err error
	chunk := &MCsLChunk{
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

func (c *MCsLChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
