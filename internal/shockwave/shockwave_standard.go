package shockwave

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/chunks"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
)

func ParseStandard(shockwave *Shockwave) {
	var err error
	var json string

	var outputFolder string
	if shockwave.PkgName == "" {
		outputFolder = filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "chunks_mmap")
	} else {
		outputFolder = filepath.Join(consts.PathDump, shockwave.PkgName, "file", filepath.Base(shockwave.FilePath), "chunks_mmap")
	}

	os.MkdirAll(outputFolder, os.ModePerm)

	var chunkMap ChunkMap = &StandardChunkMap{}
	shockwave.ChunkMap = chunkMap
	shockwave.ChunkMap.SetShockwave(shockwave)

	// --------------------------------------------------
	//  IMAP CHUNK
	// --------------------------------------------------
	// The imap chunk has some basic version information
	// such as the position of the memory map, etc.

	// After the header is the IMAP chunk, which contains the offset to the memory map
	// Handle imap and mmap

	imap, err := chunks.ReadImapChunk(shockwave.GetReader(), shockwave.Endian)
	if err != nil {
		log.Fatal(err)
	}

	imap.MemoryMapOffset -= int32(shockwave.DirOffset)

	shockwave.Version = version.ParseVersion(imap.MemoryMapFileVersion)

	json, err = imap.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	utils.SaveChunkToFileBetter("imap", 0, 1, shockwave.FilePath, json, shockwave.PkgName, "")

	// --------------------------------------------------
	//  MMAP CHUNK
	// --------------------------------------------------

	mmap, err := chunks.ReadMmapChunk(shockwave.GetReader(), shockwave.Endian, int64(imap.MemoryMapOffset), shockwave.DirOffset)
	if err != nil {
		log.Fatal(err)
	}

	json, err = mmap.ToJSON()
	if err != nil {
		log.Fatal(err)
	}
	utils.SaveChunkToFileBetter("mmap", 0, 1, shockwave.FilePath, json, shockwave.PkgName, "")

	utils.DebugMsg("shockwave", "Resource count: %d", len(mmap.Resources))

	// --------------------------------------------------
	//  KEY* CHUNK
	// --------------------------------------------------

	keyFromMap, err := mmap.FindResourceByType("KEY*")
	if err != nil {
		log.Fatal(err)
	}

	keys, err := chunks.ReadKeyChunk(shockwave.GetReader(), shockwave.Endian, int64(keyFromMap.Offset))
	if err != nil {
		log.Fatal(err)
	}

	for i := range keys.Records {
		record := keys.Records[i] // Make a copy of the record
		if !record.IsValid() {
			continue
		}
		if record.ElementIndex > 0 {
			resource, err := mmap.FindResourceByID(int(record.ElementIndex))
			if err == nil {

				resource.KeyRecord = record
				utils.DebugMsg("shockwave", "mapped %s to %d", resource.ChunkType, record.CastNumber)
			} else {
				utils.WarnMsg("shockwave", "(?) Could not find ResourceId for KEY* mapping: %v", record.ElementIndex)
			}
		}
	}

	for i := range mmap.Resources {
		mresource := mmap.Resources[i]

		if mresource.Offset == -1 {
			continue
		}

		var binaryData = make([]byte, mresource.Size+8)
		_, err = shockwave.GetReader().Seek(int64(mresource.Offset), 0)
		if err != nil {
			log.Fatal(fmt.Errorf("error seeking to offset %d: %s", mresource.Offset, err))
		}

		_, err = shockwave.GetReader().Read(binaryData)
		if err != nil {
			log.Fatal(fmt.Errorf("error reading %d bytes at offset %d: %s", mresource.Size, mresource.Offset, err))
		}

		var resource = &ShockwaveResource{
			ResourceId:       int32(mresource.ResourceId),
			Offset:           int32(mresource.Offset),
			CompressedSize:   mresource.Size,
			UncompressedSize: mresource.Size,
			CompressionType:  1,
			ChunkType:        mresource.ChunkType,
			Binary:           binaryData[8:],
		}

		if mresource.KeyRecord.ElementIndex > 0 {
			resource.CastId = mresource.KeyRecord.CastIndex
		}

		resource.DumpBinary(outputFolder)
		chunkMap.AddResource(resource)
	}

}
