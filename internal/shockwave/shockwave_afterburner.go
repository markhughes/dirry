package shockwave

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"path/filepath"

	"github.com/markhughes/dirry/internal/chunks"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
)

func ParseAfterburner(shockwave *Shockwave) {

	// --------------------------------------------------
	//  FVER CHUNK
	// --------------------------------------------------
	// The fver chunk has some basic version information

	fver, err := chunks.ReadFverChunk(shockwave.GetReader(), shockwave.Endian)
	if err != nil {
		utils.ErrorMsg("shockwave", "Error reading FVER chunk: %s", err)
		return
	}

	fver.Chunk.DecompressedDump(filepath.Base(shockwave.FilePath), "Fver", "fver", shockwave.PkgName)
	json, err := fver.ToJSON()
	if err != nil {
		utils.ErrorMsg("shockwave", "Error converting FVER chunk to JSON: %s", err)
		return
	} else {
		utils.SaveChunkToFileBetter("fver", int(fver.Chunk.StartPosition), -9, shockwave.FilePath, json, shockwave.PkgName, "")

	}
	shockwave.Version = version.ParseVersion(int32(fver.Version))
	utils.InfoMsg("shockwave", "Shockwave version: %s", shockwave.Version.ToString())

	// --------------------------------------------------
	//  FCDR CHUNK
	// --------------------------------------------------
	// The fcdr chunk has info about compressions

	fcdr, err := chunks.ReadFcdrChunk(shockwave.GetReader(), shockwave.Endian)
	if err != nil {
		utils.ErrorMsg("shockwave", "Error reading FCDR chunk: %s", err)
	}
	fcdr.Chunk.DecompressedDump(filepath.Base(shockwave.FilePath), "Fcdr", "fcdr", shockwave.PkgName)

	json, err = fcdr.ToJSON()
	if err != nil {
		utils.ErrorMsg("shockwave", "Error converting FCDR chunk to JSON: %s", err)
	} else {
		utils.SaveChunkToFileBetter("fcdr", int(fver.Chunk.StartPosition), -9, shockwave.FilePath, json, shockwave.PkgName, "")

	}

	// --------------------------------------------------
	//  ABMP CHUNK
	// --------------------------------------------------
	// the abmp chunk has info about the resources in the
	// file.

	// ABMP after fcdr
	abmp, err := chunks.ReadABMPChunk(shockwave.GetReader(), shockwave.Endian)
	if err != nil {
		utils.ErrorMsg("shockwave", "Error reading ABMP chunk: %s", err)
	}

	// create a count of resources
	var count = make(map[string]int)
	for i := range abmp.Resources {
		count[abmp.Resources[i].ChunkType]++
	}

	// print out the count
	for k, v := range count {
		utils.DebugMsg("shockwave", "Resource type: %s, count: %d", k, v)
	}

	json, err = abmp.ToJSON()
	if err != nil {
		utils.ErrorMsg("shockwave", "Error converting abmp chunk to JSON: %s", err)
	} else {
		utils.SaveChunkToFileBetter("abmp", int(abmp.Chunk.StartPosition), 1, shockwave.FilePath, json, shockwave.PkgName, "")

	}

	abmp.Chunk.DecompressedDump(filepath.Base(shockwave.FilePath), "ABMP", "abmp", shockwave.PkgName)

	// --------------------------------------------------
	//  FGEI CHUNK
	// --------------------------------------------------
	// This is more of an entry point it seems to use as
	// an offset to the resources in the file.

	fgei, err := chunks.ReadFGEIChunk(shockwave.GetReader(), shockwave.Endian, abmp)
	if err != nil {
		utils.ErrorMsg("shockwave", "Error reading FGEI chunk: %s", err)
	}

	fgei.Chunk.DecompressedDump(filepath.Base(shockwave.FilePath), "FGEI", "FGEI", shockwave.PkgName)

	for i := range abmp.Resources {
		if abmp.Resources[i].Offset != -1 {
			abmp.Resources[i].Offset = int32(abmp.Resources[i].Offset) + int32(fgei.Chunk.Position()) + fgei.Chunk.Length

		}
	}

	// --------------------------------------------------
	//  CREATE CHUNK MAP + ILS PREP
	// --------------------------------------------------
	// Create a chunkmap so later we can easily seek
	// resources from inside directory, also, we need
	// to keep a record of what to pull from the ILS

	var outputFolder string
	if shockwave.PkgName == "" {
		outputFolder = filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "chunks_abmp")
	} else {
		outputFolder = filepath.Join(consts.PathDump, shockwave.PkgName, "file", filepath.Base(shockwave.FilePath), "chunks_abmp")
	}

	// RESOURCE MAPPING
	var ilsResourcesMap = make(map[uint32]*chunks.AfterburnerResource)

	{
		var chunkMap ChunkMap = &AfterburnerChunkMap{}

		shockwave.ChunkMap = chunkMap
		shockwave.ChunkMap.SetShockwave(shockwave)

		for i := range abmp.Resources {
			var resource = abmp.Resources[i]

			if resource.Offset == -1 {
				// If the offset is -1, it means the resource is in the ILS
				// so there is nothing else to read from here, remember it
				// and move to the next one.
				ilsResourcesMap[resource.ResourceId] = resource
				continue
			}

			// in the file, move to offset, read size
			// then add to chunkmap

			_, err := shockwave.GetReader().Seek(int64(resource.Offset), io.SeekStart)
			if err != nil {
				fmt.Printf("Error seeking to resource offset: %s\n", err)
			}

			var data = make([]byte, resource.CompressedLength)
			if len(data) != int(resource.CompressedLength) {
				fmt.Printf("Size mismatch: %d != %d\n", len(data), resource.CompressedLength)
			}
			shockwave.GetReader().Read(data)

			if resource.CompressionType == 0 {
				zlibReader, err := zlib.NewReader(bytes.NewReader(data))
				if err != nil {
					fmt.Printf("error creating zlib reader: %s", err)
					continue
				}
				defer zlibReader.Close()

				// Set the limit on the zlib reader to the original uncompressed data length
				limitedReader := &io.LimitedReader{R: zlibReader, N: int64(resource.DecompressedLength)}

				// Read all data from the limited reader (this will be the decompressed data)
				data, err = io.ReadAll(limitedReader)
				if err != nil {
					fmt.Printf("error reading decompressed data: %s", err)
					continue
				}

				if len(data) != int(resource.DecompressedLength) {
					fmt.Printf("Size mismatch: %d != %d\n", len(data), resource.DecompressedLength)
				}

			}

			var res = &ShockwaveResource{
				ResourceId:       int32(resource.ResourceId),
				Offset:           int32(resource.Offset),
				CompressedSize:   int32(resource.CompressedLength),
				UncompressedSize: int32(resource.DecompressedLength),
				CompressionType:  int32(resource.CompressionType),
				ChunkType:        resource.ChunkType,
				Binary:           data,
			}
			res, err = chunkMap.AddResource(res)
			if err != nil {
				fmt.Printf("Error adding resource: %s\n", err)
			}

			err = res.DumpBinary(outputFolder)
			if err != nil {
				fmt.Printf("Error dumping binary from ABMP (size: %d): %s\n", resource.Offset, err)
			}

		}
	}

	// --------------------------------------------------
	//  ILS CHUNK
	// --------------------------------------------------
	// From the chunk map we can look up the ILS chunk
	// and read the resources from it.
	// These resources are just one after another.

	for i, v := range ilsResourcesMap {
		utils.DebugMsg("shockwave", "ILS resource: %d %s %d %d %d", i, v.ChunkType, v.CompressionType, v.CompressedLength, v.DecompressedLength)
	}
	if shockwave.PkgName == "" {
		outputFolder = filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "chunks_ils")
	} else {
		outputFolder = filepath.Join(consts.PathDump, shockwave.PkgName, "file", filepath.Base(shockwave.FilePath), "chunks_ils")
	}

	res := shockwave.ChunkMap.GetResourcesByTag("ILS ")

	if (res == nil) || (len(res) == 0) {
		fmt.Printf("No ILS resources found\n")
		return
	} else {
		var ils = res[0]
		ilsReader, err := ils.GetReader()
		if err != nil {
			fmt.Printf("Error getting reader for ILS resource: %s\n", err)
			return
		}
		for {
			resourceId, _, err := ilsReader.ReadVarInt()
			if err != nil {
				if err == io.EOF {
					break
				}

				utils.ErrorMsg("shockwave", "Error reading resource id from ils: %s", err)
				return
			}

			var afterburnerResource = ilsResourcesMap[resourceId]
			fmt.Printf("ILS resource: %d %s %d %d %d\n", resourceId, afterburnerResource.ChunkType, afterburnerResource.CompressionType, afterburnerResource.CompressedLength, afterburnerResource.DecompressedLength)
			if afterburnerResource == nil || afterburnerResource.CompressedLength == 0 {
				break
			}

			bytes, err := ilsReader.ReadBytes(int(afterburnerResource.DecompressedLength))
			if err != nil {
				utils.ErrorMsg("shockwave", "Error reading bytes: %s", err)
			}

			var res = &ShockwaveResource{
				ResourceId:       int32(resourceId),
				Offset:           int32(afterburnerResource.Offset),
				CompressedSize:   int32(afterburnerResource.CompressedLength),
				UncompressedSize: int32(afterburnerResource.DecompressedLength),
				CompressionType:  int32(afterburnerResource.CompressionType),
				ChunkType:        afterburnerResource.ChunkType,
				Binary:           bytes,
			}
			res, err = shockwave.ChunkMap.AddResource(res)
			if err != nil {
				utils.ErrorMsg("shockwave", "Error adding resource: %s", err)
			}

			err = res.DumpBinary(outputFolder)
			if err != nil {
				utils.ErrorMsg("shockwave", "Error dumping binary from ILS: %s", err)
			}

			if len(res.Binary) != int(afterburnerResource.DecompressedLength) {
				utils.WarnMsg("shockwave", "Warning: size mismatch: %d != %d\n", len(res.Binary), afterburnerResource.DecompressedLength)
			}

		}

	}

	// --------------------------------------------------
	//  KEY* CHUNK
	// --------------------------------------------------

	keysFromMap := shockwave.ChunkMap.GetResourcesByTag("KEY*")
	if err != nil {
		utils.ErrorMsg("shockwave", "Error finding KEY* resource: %s", err)
		return
	}

	if len(keysFromMap) == 0 {
		utils.ErrorMsg("shockwave", "Could not find KEY* resource")
		return
	}

	keyReader, err := keysFromMap[0].GetReader()
	if err != nil {
		utils.ErrorMsg("shockwave", "Error creating reader for KEY* resource: %s", err)
		return
	}

	keys, err := chunks.ReadKeyChunkRaw(keyReader.GetUnsafeBytesReader(), shockwave.Endian)
	if err != nil {
		utils.ErrorMsg("shockwave", "Error reading KEY* chunk: %s", err)
		return
	}

	for i := range keys.Records {
		record := keys.Records[i] // Make a copy of the record
		if record != nil && !record.IsValid() {
			continue
		}

		if record.ElementIndex > 0 {
			resource := shockwave.ChunkMap.GetResourceById(int32(record.ElementIndex))
			if err != nil {

				resource.CastId = record.CastIndex
				utils.DebugMsg("shockwave", "mapped %s to %d", resource.ChunkType, record.CastNumber)
			} else {
				utils.WarnMsg("shockwave", "(?) Could not find ResourceId for KEY* mapping: %v", record.ElementIndex)
			}
		}
	}

}
