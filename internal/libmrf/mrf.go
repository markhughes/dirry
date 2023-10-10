package libmrf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/markhughes/dirry/internal/macroman"
	"github.com/markhughes/dirry/internal/utils"
)

type Resource struct {
	Id         uint16
	NameOffset uint16
	PackedAttr uint32
	Flags      uint32
	Junk       uint32
	Name       string
	Data       []byte
	DataSize   uint32
	DataOffset uint32
	Type       string
}

type ResourceFork struct {
	JunkNextresmap uint32
	JunkFilerefnum uint16
	FileAttributes uint16
	Resources      []Resource
}

func FromFile(path string) (*ResourceFork, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromBytes(data)
}

func FromBytes(data []byte) (*ResourceFork, error) {
	rf := &ResourceFork{}
	if len(data) == 0 {
		return rf, nil
	}

	r := bytes.NewReader(data)
	var dataOffset, mapOffset, dataLength, mapLength uint32

	binary.Read(r, binary.BigEndian, &dataOffset)
	binary.Read(r, binary.BigEndian, &mapOffset)
	binary.Read(r, binary.BigEndian, &dataLength)
	binary.Read(r, binary.BigEndian, &mapLength)

	utils.DebugMsg("mrf", "dataOffset: %d\n", dataOffset)
	utils.DebugMsg("mrf", "mapOffset: %d\n", mapOffset)
	utils.DebugMsg("mrf", "dataLength: %d\n", dataLength)
	utils.DebugMsg("mrf", "mapLength: %d\n", mapLength)
	dataSection := data[dataOffset : dataOffset+dataLength]
	mapSection := data[mapOffset : mapOffset+mapLength]

	r = bytes.NewReader(mapSection[16:]) // skipping first 16 bytes

	binary.Read(r, binary.BigEndian, &rf.JunkNextresmap)
	binary.Read(r, binary.BigEndian, &rf.JunkFilerefnum)
	binary.Read(r, binary.BigEndian, &rf.FileAttributes)
	utils.DebugMsg("mrf", "JunkNextresmap: %d\n", rf.JunkNextresmap)
	utils.DebugMsg("mrf", "JunkFilerefnum: %d\n", rf.JunkFilerefnum)
	utils.DebugMsg("mrf", "FileAttributes: %d\n", rf.FileAttributes)

	var typelistOffsetInMap, namelistOffsetInMap, numTypes uint16
	binary.Read(r, binary.BigEndian, &typelistOffsetInMap)
	binary.Read(r, binary.BigEndian, &namelistOffsetInMap)
	binary.Read(r, binary.BigEndian, &numTypes)
	numTypes++

	utils.DebugMsg("mrf", "typelistOffsetInMap: %d\n", typelistOffsetInMap)
	utils.DebugMsg("mrf", "namelistOffsetInMap: %d\n", namelistOffsetInMap)
	utils.DebugMsg("mrf", "numTypes: %d\n", numTypes)

	utils.DebugMsg("mrf", "\n")

	uTypes := bytes.NewReader(mapSection[typelistOffsetInMap:])
	uNames := bytes.NewReader(mapSection[namelistOffsetInMap:])

	for i := uint16(0); i < numTypes; i++ {
		var resType [4]byte
		var resCount, reslistOffset uint16

		r.Read(resType[:])
		binary.Read(r, binary.BigEndian, &resCount)
		binary.Read(r, binary.BigEndian, &reslistOffset)
		resCount++

		utils.DebugMsg("mrf", "resType: %s\n", resType)
		utils.DebugMsg("mrf", "resCount: %d\n", resCount)
		utils.DebugMsg("mrf", "reslistOffset: %d\n", reslistOffset)

		_, err := uTypes.Seek(int64(reslistOffset), io.SeekStart)
		if err != nil {
			return nil, err
		}

		println()
		for j := uint16(0); j < resCount; j++ {
			var resource Resource
			binary.Read(uTypes, binary.BigEndian, &resource.Id)
			binary.Read(uTypes, binary.BigEndian, &resource.NameOffset)
			binary.Read(uTypes, binary.BigEndian, &resource.PackedAttr)
			binary.Read(uTypes, binary.BigEndian, &resource.Junk)

			resource.Flags = resource.PackedAttr >> 24
			resource.DataOffset = resource.PackedAttr & 0x00FFFFFF

			if resource.NameOffset != 0xFFFF {
				_, err := uNames.Seek(int64(resource.NameOffset), io.SeekStart)
				if err != nil {
					return nil, err
				}

				var nameLength uint8
				binary.Read(uNames, binary.BigEndian, &nameLength)
				name := make([]byte, nameLength)
				utils.DebugMsg("mrf", "resource.namebytes: %v\n", name)

				uNames.Read(name)
				resource.Name = macroman.ConvertMacRomanToUTF8(string(name))
			} else {
				resource.Name = fmt.Sprintf("Unknown %d %d", resource.Id, j)
			}

			resource.Type = string(resType[:])

			utils.DebugMsg("mrf", "resource.Id: %d\n", resource.Id)
			utils.DebugMsg("mrf", "resource.NameOffset: %d\n", resource.NameOffset)
			utils.DebugMsg("mrf", "resource.PackedAttr: %d\n", resource.PackedAttr)
			utils.DebugMsg("mrf", "resource.Junk: %d\n", resource.Junk)

			dataR := bytes.NewReader(dataSection[resource.DataOffset:])
			binary.Read(dataR, binary.BigEndian, &resource.DataSize)
			resourceData := make([]byte, resource.DataSize)
			dataR.Read(resourceData)
			resource.Data = resourceData

			utils.DebugMsg("mrf", "resource.Junk: %d\n", resource.Junk)
			utils.DebugMsg("mrf", "resource.DataSize: %d\n", resource.DataSize)
			utils.DebugMsg("mrf", "resource.DataOffset: %d\n", resource.DataOffset)

			rf.Resources = append(rf.Resources, resource)
			utils.DebugMsg("mrf", "[done] inner\n")

		}

		utils.DebugMsg("mrf", "[done] outter\n")

	}

	sort.Slice(rf.Resources, func(i, j int) bool {
		return rf.Resources[i].DataOffset < rf.Resources[j].DataOffset
	})

	return rf, nil
}
