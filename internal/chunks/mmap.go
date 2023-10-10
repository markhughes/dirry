package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

type Resource struct {
	ResourceId int
	ChunkType  string
	Size       int32
	Offset     int32
	Unknown1   int16
	Unknown2   int16
	Unknown3   int32

	KeyRecord *KeyRecord
}

func (r *Resource) IsValid() bool {
	_, ok := ChunkTypes[r.ChunkType]
	return ok
}

type MmapChunk struct {
	Chunk           *BinaryChunk
	HeaderSize      int16
	EntrySize       int16
	NumberOfEntries int32
	NonZeroEntries  int32
	FirstFreeEntry  int32
	FirstJunkEntry1 int32
	FirstJunkEntry2 int32
	Resources       []*Resource
}

func ReadMmapChunk(r io.ReadSeeker, endian binary.ByteOrder, offset int64, dirOffset int64) (*MmapChunk, error) {
	var err error

	utils.InfoMsg("mmap", "Reading mmap chunk at %d", offset)
	_, err = r.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	chunk := &MmapChunk{}
	chunk.Chunk, err = FromBinaryAt(r, "mmap", endian, offset)
	if err != nil {
		return nil, err
	}

	chunk.HeaderSize, err = chunk.Chunk.ReadInt16(endian)
	if err != nil {
		return nil, err
	}

	utils.InfoMsg("mmap", "Header size: %d", chunk.HeaderSize)

	chunk.EntrySize, err = chunk.Chunk.ReadInt16(endian)
	if err != nil {
		return nil, err
	}

	chunk.NumberOfEntries, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return nil, err
	}

	chunk.NonZeroEntries, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return nil, err
	}

	chunk.FirstFreeEntry, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return nil, err
	}

	chunk.FirstJunkEntry1, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return nil, err
	}

	chunk.FirstJunkEntry2, err = chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return nil, err
	}

	resources := make([]*Resource, chunk.NumberOfEntries)

	utils.InfoMsg("mmap", "Reading %d resources", chunk.NumberOfEntries)
	for i := 0; i < int(chunk.NumberOfEntries); i++ {
		res := &Resource{}
		res.ResourceId = i

		var resourceChunkType [4]byte
		if err := binary.Read(r, endian, &resourceChunkType); err != nil {
			return nil, err
		}

		if endian == binary.LittleEndian {
			res.ChunkType = string(utils.ReverseBytes(resourceChunkType[:]))
		} else {
			res.ChunkType = string((resourceChunkType[:]))

		}
		if err := binary.Read(r, endian, &res.Size); err != nil {
			return nil, err
		}

		if err := binary.Read(r, endian, &res.Offset); err != nil {
			return nil, err
		}
		if dirOffset > 0 {
			if res.ChunkType == "RIFX" || res.ChunkType == "XIFR" {
				if res.Offset != int32(dirOffset) {
					utils.InfoMsg("mmap", "Warning: RIFX/XIFR offset %d does not match dirOffset %d", res.Offset, dirOffset)
				}
			}

			res.Offset = res.Offset - int32(dirOffset)
			utils.DebugMsg("mmap", "Offset %d - %d = %d", res.Offset, dirOffset, res.Offset)

			if res.Offset < 0 {
				res.Offset = -1
			}
		}

		// fmt.Printf("Resource %d: %s, %d, %d - %d\n", i, res.ChunkType, res.Size, res.Offset, dirOffset)

		if err := binary.Read(r, endian, &res.Unknown1); err != nil {
			return nil, err
		}

		if err := binary.Read(r, endian, &res.Unknown2); err != nil {
			return nil, err
		}

		if err := binary.Read(r, endian, &res.Unknown3); err != nil {
			return nil, err
		}

		res.KeyRecord = &KeyRecord{}
		resources[i] = res
	}

	chunk.Resources = resources

	return chunk, nil
}
func (m *MmapChunk) FindResourceByType(chunkType string) (*Resource, error) {
	for i := range m.Resources {
		if string((m.Resources)[i].ChunkType[:]) == chunkType {
			return (m.Resources)[i], nil
		}
	}
	return nil, fmt.Errorf("resource with type %s not found", chunkType)
}

func (m *MmapChunk) FindResourceByID(resourceId int) (*Resource, error) {
	for i := range m.Resources {
		if (m.Resources)[i].ResourceId == resourceId {
			return (m.Resources)[i], nil
		}
	}
	return nil, fmt.Errorf("resource with resourceId %v not found", resourceId)
}

func (m *MmapChunk) FindResourcesByType(chunkType string) ([]*Resource, error) {
	utils.DebugMsg("mmap", "Looking for %s\n", chunkType)
	var resources []*Resource

	for i := range m.Resources {
		if (m.Resources)[i].ChunkType == chunkType {
			utils.DebugMsg("mmap", "Found %s at %d\n", (m.Resources)[i].ChunkType, i)

			resources = append(resources, (m.Resources)[i])
		}
	}

	if len(resources) == 0 {
		return nil, fmt.Errorf("resources with type %s not found", chunkType)
	}

	return resources, nil
}

func (c *MmapChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
