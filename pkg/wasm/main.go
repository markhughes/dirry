//go:build js
// +build js

package main

import (
	"fmt"
	"path/filepath"
	"syscall/js"

	"github.com/markhughes/dirry/internal/chunks"
	"github.com/markhughes/dirry/internal/errors"
	"github.com/markhughes/dirry/internal/palettes"
	"github.com/markhughes/dirry/internal/shockwave"
	"github.com/markhughes/dirry/internal/utils"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("processFile", js.FuncOf(wasmExtract))

	utils.EnabledDebugAll = true

	<-c
}

func wasmExtract(this js.Value, args []js.Value) interface{} {
	filePath := args[0].String()

	arrayBuffer := args[1]

	// Determine the length of the ArrayBuffer
	length := arrayBuffer.Get("byteLength").Int()

	// Create a Go byte slice to hold the data
	fileContent := make([]byte, length)

	// Copy the bytes from the ArrayBuffer to the Go slice
	js.CopyBytesToGo(fileContent, arrayBuffer)

	callbackChunk := args[2]
	callbackBinary := args[3]

	extractor(filePath, fileContent, "", 0, callbackChunk, callbackBinary)

	return nil
}

func extractor(filePath string, fileContent []byte, pkg string, extraOffset int64, callbackChunk js.Value, callbackBinary js.Value) {

	utils.PrintHeader()
	utils.InfoMsg("dump", " > The Directory Utility <\n\n")

	var err error

	var shockwave shockwave.Shockwave
	shockwave.PkgName = pkg
	shockwave.DirOffset = (extraOffset)

	expanded, err := shockwave.OpenContent(filePath, fileContent)
	if err != nil {
		utils.ErrorMsg("dump", "Error opening file %s: %s\n", filePath, err)
		return
	}

	if len(expanded) > 0 {
		utils.InfoMsg("dump", "Expanded %s to %d files\n", filePath, len(expanded))

		for i := range expanded {
			utils.InfoMsg("dump", "Dumping %s with offset %d\n", expanded[i].Path, int64(expanded[i].MinusOffset))
			// Extractor(expanded[i].Path, filepath.Base(filePath), int64(expanded[i].MinusOffset))
		}
		return
	}

	chunkMapJson, err := shockwave.ChunkMap.ToJson()
	if err != nil {
		utils.ErrorMsg("dump", "Error converting chunkmap to JSON: %s\n", err)
		return
	}

	callbackChunk.Invoke(
		js.ValueOf("+chunkmap"),
		js.ValueOf(0),
		js.ValueOf(0),
		js.ValueOf(filePath),
		js.ValueOf(chunkMapJson),
		js.ValueOf(shockwave.PkgName),
		js.ValueOf(""),
	)

	//callbackChunk.Invoke("+chunkmap", 0, 0, filePath, chunkMapJson, shockwave.PkgName, "")
	// utils.SaveChunkToFileBetter("+chunkmap", 0, 0, filePath, chunkMapJson, shockwave.PkgName, "")

	var resources = shockwave.ChunkMap.GetAllResources()

	var content = ""
	var pendingResourceIds = make([]int, 0)

	for i := range resources {
		content = ""
		var resource = resources[i]

		// if its an empty chunk we cannot do much with it
		if resource.UncompressedSize == 0 {
			continue
		}

		utils.DebugMsg("dump", "Attending %s", resource.ChunkType)
		switch resources[i].ChunkType {

		// Skippable
		case "ILS ":
			continue
		case "mmap":
			continue

		// Chunks
		case "KEY*":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for KEY* resource: %s\n", err)
				break
			}

			keychunk, err := chunks.ReadKeyChunkRaw(reader.GetUnsafeBytesReader(), shockwave.Endian)
			if err != nil {
				utils.ErrorMsg("dump", "Error reading KEY* chunk: %s\n", err)
				break
			}

			content, err = keychunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting KEY* chunk to JSON: %s\n", err)
				break
			}

		case "CASt":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for CASt resource: %s\n", err)
				break
			}

			castchunk, err := chunks.ReadCastChunkRaw(reader.GetUnsafeBytesReader(), shockwave.Version, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				if _, ok := err.(*errors.UnhandledCastTypeError); ok {
					utils.WarnMsg("dump", "Could not decode CASt chunk: %s", err)
				} else {
					utils.ErrorMsg("dump", "Error reading CASt chunk: %s\n", err)
					break

				}
			}

			content, err = castchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting CASt chunk to JSON: %s\n", err)
				break
			}

			shockwave.Casts[resource.ResourceId] = castchunk

			// castchunk.Save(filepath.Base(shockwave.FilePath), fmt.Sprint(resource.ResourceId), shockwave.PkgName)

		case "ediM":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for ediM resource: %s\n", err)
				break
			}

			edimchunk, err := chunks.ReadEdimChunkRaw(reader.GetUnsafeBytesReader(), reader.GetUnsafeBytesReader().Len(), shockwave.Endian)
			if err != nil {
				utils.ErrorMsg("dump", "Error reading ediM chunk: %s\n", err)
				break
			}

			content, err = edimchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting ediM chunk to JSON: %s\n", err)
				break
			}

			// edimchunk.Save(filepath.Base(shockwave.Reader.Name()), fmt.Sprint(resource.ResourceId)+"_"+fmt.Sprint(i), shockwave.PkgName)

		case "XTRl":
			reader, err := resource.GetReader()

			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for Xtrl resource: %s\n", err)
				break
			}

			xtrlchunk, err := chunks.ReadXtrlChunkRaw(reader.GetUnsafeBytesReader(), shockwave.Endian)
			if err != nil {
				utils.ErrorMsg("dump", "Error reading Xtrl chunk: %s\n", err)
				break
			}

			content, err = xtrlchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting Xtrl chunk to JSON: %s\n", err)
				break
			}

		case "GRID":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for GRID resource: %s\n", err)
				break
			}

			gridchunk, err := chunks.ReadGridChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading GRID chunk: %s\n", err)
				break
			}

			content, err = gridchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting GRID chunk to JSON: %s\n", err)
				break
			}

		case "CAS*":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for CAS* resource: %s\n", err)
				break
			}

			caschunk, err := chunks.ReadCasChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading CAS* chunk: %s\n", err)
				break
			}

			content, err = caschunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting CAS* chunk to JSON: %s\n", err)
				break
			}

		case "Sord":

			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for Sord resource: %s\n", err)
				break
			}

			sordchunk, err := chunks.ReadSordChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading Sord chunk: %s", err)
				break
			}

			content, err = sordchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting Sord chunk to JSON: %s\n", err)
				break
			}

		case "FCOL":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for FCOL resource: %s\n", err)
				break
			}

			fcol, err := chunks.ReadFcolChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading FCOL chunk: %s\n", err)
				break
			}
			content, _ = fcol.ToJSON()

			// fcol.Save(filepath.Base(shockwave.Reader.Name()), filepath.Base(shockwave.Reader.Name())+" favourites", shockwave.PkgName)

		case "CLUT":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for CLUT resource: %s\n", err)
				break
			}

			clutchunk, err := chunks.ReadClutChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading CLUT chunk: %s\n", err)
				break
			}

			content, err = clutchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting CLUT chunk to JSON: %s\n", err)
				break
			}

			palettes.RegisterPallete(palettes.Clut(resource.ResourceId), clutchunk.Palette)

			// clutchunk.Save(filepath.Base(shockwave.Reader.Name()), fmt.Sprint(resource.ResourceId), "")

		case "LctX":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for LctX resource: %s\n", err)
				break
			}

			lctxchunk, err := chunks.ReadLctxChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading LctX chunk: %s\n", err)
				break
			}

			content, err = lctxchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting LctX chunk to JSON: %s", err)
				break
			}

		case "Lnam":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for Lnam resource: %s", err)
				break
			}

			lnamchunk, err := chunks.ReadLnamChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading Lnam chunk: %s", err)
				break
			}

			content, err = lnamchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting Lnam chunk to JSON: %s", err)
				break
			}

		case "Lscr":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for Lscr resource: %s", err)
				break
			}

			lscrchunk, err := chunks.ReadLscrChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading Lscr chunk: %s", err)
				break
			}

			content, err = lscrchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting Lscr chunk to JSON: %s", err)
				break
			}

		case "VWFI":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for VWFI resource: %s", err)
				break
			}

			vwfichunk, err := chunks.ReadVwfiChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading VWFI chunk: %s", err)
				break
			}

			content, err = vwfichunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting VWFI chunk to JSON: %s", err)
				break
			}

		case "VWLB":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for VWLB resource: %s", err)
				break
			}

			vwlbchunk, err := chunks.ReadVwlbChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading VWLB chunk: %s", err)
				break
			}

			content, err = vwlbchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting VWLB chunk to JSON: %s", err)
				break
			}

		case "VWCF", "DRCF":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for DRCF resource: %s", err)
				break
			}

			drcfchunk, err := chunks.ReadInfoChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading DRCF chunk: %s", err)
				break
			}

			content, err = drcfchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting DRCF chunk to JSON: %s", err)
				break
			}

		case "MCsL":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for MCsL resource: %s", err)
				break
			}

			mcslchunk, err := chunks.ReadMcslChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading MCsL chunk: %s", err)
				break
			}

			content, err = mcslchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting MCsL chunk to JSON: %s", err)
				break
			}

		case "FXmp":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for FXmp resource: %s", err)
				break
			}

			fxmpchunk, err := chunks.ReadFontXMapChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading FXmp chunk: %s", err)
				break
			}

			content, err = fxmpchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting FXmp chunk to JSON: %s", err)
				break
			}

			// fxmpchunk.Save(filepath.Base(shockwave.Reader.Name()), fmt.Sprint(resource.ResourceId), shockwave.PkgName)

		case "Fmap":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for Fmap resource: %s", err)
				break
			}

			fmapchunk, err := chunks.ReadFmapChunkRaw(reader, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading Fmap chunk: %s", err)
				break
			}

			content, err = fmapchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting Fmap chunk to JSON: %s", err)
				break
			}

			for i := range fmapchunk.Fonts {
				var font = fmapchunk.Fonts[i]
				shockwave.Fonts[uint32(font.FontID)] = font
				utils.DebugMsg("dump", "Adding font %d: %s", font.FontID, font.Name)

			}
		}

		// save to json
		if resource.UncompressedSize > 0 {
			if content != "" {
				utils.SuccessMsg("dump", "Processed %s", resource.ChunkType)
				callbackChunk.Invoke(
					js.ValueOf(resource.ChunkType),
					js.ValueOf(int(resource.Offset)),
					js.ValueOf(int(resource.UncompressedSize)),
					js.ValueOf(filePath),
					js.ValueOf(content),
					js.ValueOf(shockwave.PkgName),
					js.ValueOf(""),
				)

				// utils.SaveChunkToFileBetter(resource.ChunkType, int(resource.Offset), int(resource.UncompressedSize), filePath, content, shockwave.PkgName, "")
			} else {
				pendingResourceIds = append(pendingResourceIds, i)
			}

		}
	}

	// These resources are dependent on a cast chunk or something else being parsed first
	for _, i := range pendingResourceIds {
		content = ""
		var resource = resources[i]

		switch resource.ChunkType {
		case "STXT":
			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for STXT resource: %s\n", err)
				break
			}

			stxtchunk, err := chunks.ReadStxtChunkRaw(reader, shockwave.Fonts, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading STXT chunk: %s\n", err)
				break
			}

			content, err = stxtchunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting STXT chunk to JSON: %s\n", err)
				break
			}

		case "snd ":
			// TODO
			break

		case "BITD":
			var cast = shockwave.Casts[resource.CastId]
			if cast == nil {
				utils.ErrorMsg("dump", "BITD cast %d not found for resource %d", resource.CastId, resource.ResourceId)
				break
			}

			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for BITD resource: %s", err)
				break
			}

			chunk, err := chunks.ReadBitmapChunkRaw(reader, shockwave.Endian, cast, shockwave.IsAfterburner())
			if err != nil {
				utils.ErrorMsg("dump", "Error reading BITD chunk: %s", err)
				break
			}

			content, err = chunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting BITD chunk to JSON: %s", err)
			}

			chunk.Save(filepath.Base(shockwave.FilePath), fmt.Sprint(resource.ResourceId), shockwave.PkgName)

			// TODO

			break

		case "XMED":
			// even if the cast is not found, we will try to detect it and parse it ourselves
			var cast = shockwave.Casts[resource.CastId]

			reader, err := resource.GetReader()
			if err != nil {
				utils.ErrorMsg("dump", "Error getting reader for XMED resource: %s", err)
				break
			}

			chunk, err := chunks.ReadXmedChunkRaw(reader, cast, shockwave.Endian, shockwave.IsAfterburner())
			if err != nil {
				if _, ok := err.(*errors.UnhandledXtraTypeError); ok {
					utils.WarnMsg("dump", "Could not decode XMED chunk: %s", err)
				} else {
					utils.ErrorMsg("dump", "Error reading XMED chunk: %s", err)
					break
				}
			}

			content, err = chunk.ToJSON()
			if err != nil {
				utils.ErrorMsg("dump", "Error converting XMED chunk to JSON: %s", err)
				break
			}

			if chunk.Decoded {

				// chunk.Save(filepath.Base(shockwave.FilePath), fmt.Sprint(resource.ResourceId), shockwave.PkgName)
				utils.SuccessMsg("dump", "XMED chunk decoded")
			} else {
				utils.ErrorMsg("dump", "XMED chunk not decoded?")
			}
		}

		if resource.UncompressedSize > 0 {
			if content != "" {
				utils.SuccessMsg("dump", "Processed %s", resource.ChunkType)
				callbackChunk.Invoke(
					js.ValueOf(resource.ChunkType),
					js.ValueOf(int(resource.Offset)),
					js.ValueOf(int(resource.UncompressedSize)),
					js.ValueOf(filePath),
					js.ValueOf(content),
					js.ValueOf(shockwave.PkgName),
					js.ValueOf(""),
				)

				// utils.SaveChunkToFileBetter(resource.ChunkType, int(resource.Offset), int(resource.UncompressedSize), filePath, content, shockwave.PkgName, "")
			} else {
				utils.ErrorMsg("dump", "Did not convert chunk: %s", resource.ChunkType)

				callbackChunk.Invoke(
					js.ValueOf("incomplete_"+resource.ChunkType),
					js.ValueOf(int(resource.Offset)),
					js.ValueOf(int(resource.UncompressedSize)),
					js.ValueOf(filePath),
					js.ValueOf(content),
					js.ValueOf(shockwave.PkgName),
					js.ValueOf(""),
				)

				// utils.SaveChunkToFileBetter("incomplete_"+resource.ChunkType, int(resource.Offset), int(resource.UncompressedSize), filePath, content, shockwave.PkgName, "")

			}
		}

	}
}
