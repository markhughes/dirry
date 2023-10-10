package chunks

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/utils"
)

type FontXMapChunk struct {
	Reader *binary_reader.BinaryReader

	FontMappings map[string]*FontMapping
	CharMappings map[string][]*CharMapping
}

type FontMapping struct {
	OriginalFont   string
	SubstituteFont string
	SizeMapping    map[int]int
	MapNone        bool
}

type CharMapping struct {
	OriginalPlatform string
	TargetPlatform   string
	From             int
	To               int
}

// TODO: does not always work
func (chunk *FontXMapChunk) Read(endian binary.ByteOrder) error {

	chunk.FontMappings = make(map[string]*FontMapping)
	chunk.CharMappings = make(map[string][]*CharMapping)

	scanner := bufio.NewScanner(bytes.NewReader(chunk.Reader.GetBytes())) // skip the first 8 bytes with are the 4cc and chunk size
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, ";") || line == "" {
			continue
		}

		if strings.Contains(line, ": =>") {
			// Character Mappings
			// Platform: => Platform:  oldChar => oldChar ...
			utils.DebugMsg("fxmp", "handling char mapping: %s\n", line)
			parts := strings.SplitN(line, ": =>", 2)
			if len(parts) < 2 {
				continue
			}

			// Get the platforms
			platformParts := strings.SplitN(parts[1], ":", 2)
			if len(platformParts) < 2 {
				continue
			}
			originalPlatform := strings.TrimSpace(parts[0])
			targetPlatform := strings.TrimSpace(platformParts[0])

			// remove target platform and split by spaces
			mappingParts := strings.Fields(strings.TrimPrefix(platformParts[1], targetPlatform+":"))
			for _, pair := range mappingParts {
				charParts := strings.Split(pair, "=>")
				if len(charParts) < 2 {
					continue
				}

				oldChar, err := strconv.Atoi(charParts[0])
				if err != nil {
					continue
				}

				newChar, err := strconv.Atoi(charParts[1])
				if err != nil {
					continue
				}

				chunk.CharMappings[originalPlatform+":"+targetPlatform] = append(chunk.CharMappings[originalPlatform+":"+targetPlatform], &CharMapping{
					OriginalPlatform: originalPlatform,
					TargetPlatform:   targetPlatform,
					From:             oldChar,
					To:               newChar,
				})
			}
		} else if strings.Contains(line, "=>") {
			// Font Mapping
			// Platform:FontName => Platform:FontName [MAP NONE] [oldSize => newSize]

			// Mac:"New York"    => Win:"MS Serif"
			// Mac:Symbol        => Win:Symbol  Map None
			// Mac:Times         => Win:"Times New Roman" 14=>12 18=>14 24=>18 30=>24
			// Mac:Palatino      => Win:"Times New Roman"

			// this is a font mapping

			//  >>> Mac:Times         => Win:"Times New Roman" 14=>12 18=>14 24=>18 30=>24
			parts := strings.SplitN(line, "=>", 2)
			if len(parts) < 2 {
				continue
			}
			// >>> Mac:Times
			originalFont := strings.TrimSpace(parts[0])

			// >>> Win:"Times New Roman" 14=>12 18=>14 24=>18 30=>24
			substituteFont := strings.TrimSpace(parts[1])

			mapNone := false
			sizeMapping := make(map[int]int)

			if strings.Contains(substituteFont, "Map None") {
				mapNone = true
				substituteFont = strings.ReplaceAll(substituteFont, "Map None", "")
				substituteFont = strings.TrimSpace(substituteFont)
			}

			// >>> Win:"Times New Roman" 14=>12 18=>14 24=>18 30=>24

			// TODO: now we need to break the size mapping out into a map
			// substituteFont should be come Win:"Times New Roman"
			// sizeMapping should become [14=,12],[18,14],[24,18],[30,24]

			var sizeMappingsStr = ""
			if !strings.Contains(substituteFont, "\"") {
				// Handle case for non-spaced font names like 'Win:Symbol'

				substituteFont = strings.TrimSpace(substituteFont)

				parts := strings.SplitN(substituteFont, " ", 2)
				if len(parts) < 2 {
					continue
				}

				substituteFont = parts[0]
				sizeMappingsStr = parts[1]
			} else {

				// Find the index of the last quotation mark
				lastQuoteIndex := strings.LastIndex(substituteFont, "\"")

				// If there are no quotation marks, we have a problem
				if lastQuoteIndex == -1 {
					fmt.Println("Invalid format: no closing quotation mark in substitute font")
					continue
				}

				// Get the font name (including the space after it) and the size mappings
				sizeMappingsStr = substituteFont[lastQuoteIndex+1:]
				substituteFont = substituteFont[:lastQuoteIndex+1]

			}

			// Trim spaces
			substituteFont = strings.TrimSpace(substituteFont)
			sizeMappingsStr = strings.TrimSpace(sizeMappingsStr)

			if sizeMappingsStr != "" {
				sizeMappingsParts := strings.Split(sizeMappingsStr, " ")
				for _, mapping := range sizeMappingsParts {
					sizes := strings.Split(mapping, "=>")
					if len(sizes) != 2 {
						fmt.Println("Invalid format: expected size mappings in the form 'old=>new'")
						continue
					}

					oldSize, err := strconv.Atoi(sizes[0])
					if err != nil {
						fmt.Println("Invalid format: expected old size to be an integer")
						continue
					}

					newSize, err := strconv.Atoi(sizes[1])
					if err != nil {
						fmt.Println("Invalid format: expected new size to be an integer")
						continue
					}

					sizeMapping[oldSize] = newSize
				}
			}
			chunk.FontMappings[originalFont] = &FontMapping{
				OriginalFont:   originalFont,
				SubstituteFont: substituteFont,
				SizeMapping:    sizeMapping,
				MapNone:        mapNone,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil

}

func ReadFontXMapChunkRaw(r *binary_reader.BinaryReader, endian binary.ByteOrder, isAfterburner bool) (*FontXMapChunk, error) {
	var err error
	chunk := &FontXMapChunk{
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

func (c *FontXMapChunk) Print() {
	fmt.Println("fxmp", "-- FXmp")
}

func (c *FontXMapChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (c *FontXMapChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "FXmp")
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "FXmp")
	}

	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists

	// drop the raw TXT file
	targetFile, err := os.Create(filepath.Join(outputFolder, name+".FONTMAP.txt"))
	if err != nil {
		panic(err)
	}
	defer targetFile.Close()

	_, err = targetFile.Write(c.Reader.GetBytes())
	if err != nil {
		panic(err)
	}

	err = targetFile.Sync()
	if err != nil {
		panic(err)
	}

	// charmap
	charmapFile, err := os.Create(filepath.Join(outputFolder, name+".charmap.json"))
	if err != nil {
		panic(err)
	}
	defer charmapFile.Close()

	bytes, err := json.MarshalIndent(c.CharMappings, "", "  ")
	if err != nil {
		panic(err)
	}
	_, err = charmapFile.Write(bytes)
	if err != nil {
		panic(err)
	}

	err = charmapFile.Sync()
	if err != nil {
		panic(err)
	}

	// fontmap
	fontmapFile, err := os.Create(filepath.Join(outputFolder, name+".fontmap.json"))
	if err != nil {
		panic(err)
	}
	defer fontmapFile.Close()

	bytes, err = json.MarshalIndent(c.FontMappings, "", "  ")
	if err != nil {
		panic(err)
	}

	_, err = fontmapFile.Write(bytes)
	if err != nil {
		panic(err)
	}

	err = fontmapFile.Sync()
	if err != nil {
		panic(err)
	}

}
