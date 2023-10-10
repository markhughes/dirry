package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/binary_reader"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/errors"
	"github.com/markhughes/dirry/internal/members"
	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/xmed"
)

type XmedChunk struct {
	Reader    *binary_reader.BinaryReader
	Member    *members.MemberXtra
	Extension string
	Decoded   bool

	Data []byte
	Meta []byte
}

func (chunk *XmedChunk) Read(castChunk *CastChunk, endian binary.ByteOrder) error {
	var err error

	var member *members.MemberXtra
	var ok bool

	if castChunk == nil {

		bytes, err := chunk.Reader.ReadBytes(15)
		if err != nil {
			return fmt.Errorf("could not read xmed chunk type: %v", err)
		}

		chunk.Reader.Seek(0, 0)

		if bytes[0] == 'P' && bytes[1] == 'F' && bytes[2] == 'R' {
			utils.InfoMsg("xmed", "found PFR1")
			chunk.Member = &members.MemberXtra{
				Type: "font",
			}
		} else if bytes[0] == '3' && bytes[1] == 'D' && bytes[2] == 'E' && bytes[3] == 'M' {
			utils.InfoMsg("xmed", "found 3DEM")
			chunk.Member = &members.MemberXtra{
				Type: "shockwave3d",
			}
		} else if bytes[12] == 'F' && bytes[13] == 'W' && bytes[14] == 'S' {
			utils.InfoMsg("xmed", "found FWS")
			chunk.Member = &members.MemberXtra{
				Type: "flash",
			}
		} else {
			utils.InfoMsg("xmed", "XMED chunk has no cast chunk, and we could not determine the type (%s).", bytes)

			return fmt.Errorf("could not determine xmed chunk type")
		}

	} else {
		if member, ok = castChunk.Member.(*members.MemberXtra); ok {
			chunk.Member = member
		} else {
			return fmt.Errorf("castChunk.Member is not a MemberXtra")
		}
	}

	chunk.Extension = "unknown." + chunk.Member.Type + ".bin"

	switch chunk.Member.Type {
	case "flash":
		chunk.Data, err = xmed.CreateFlashBinary(chunk.Reader.GetUnsafeBytesReader(), endian)
		if err != nil {
			return fmt.Errorf("could not create flash binary: %v", err)
		}
		chunk.Extension = "swf"
		chunk.Decoded = true

	case "vectorShape": // director used flash for vectorShapes lol
		chunk.Data, err = xmed.CreateFlashBinary(chunk.Reader.GetUnsafeBytesReader(), endian)
		if err != nil {
			return fmt.Errorf("could not create flash binary: %v", err)
		}
		chunk.Extension = "swf"
		chunk.Decoded = true

	case "havok": // -- is this just a 3d file?
	case "shockwave3d":
		chunk.Data, err = xmed.CreateShockwave3DBinary(chunk.Reader.GetUnsafeBytesReader(), endian)
		if err != nil {
			return fmt.Errorf("could not create shockwave3d binary: %v", err)
		}

		chunk.Extension = "ifx"
		chunk.Decoded = true

	case "font":
		chunk.Data, chunk.Meta, err = xmed.CreateFontBinary(chunk.Reader, endian)
		if err != nil {
			return fmt.Errorf("could not create font binary: %v", err)
		}

		chunk.Extension = "pfr"
		chunk.Decoded = true

	default:
		chunk.Decoded = true

		return &errors.UnhandledXtraTypeError{XtraType: chunk.Member.Type}
	}
	return nil

}

func ReadXmedChunkRaw(r *binary_reader.BinaryReader, castChunk *CastChunk, endian binary.ByteOrder, isAfterburner bool) (*XmedChunk, error) {
	var err error
	chunk := &XmedChunk{
		Reader: r,
	}

	chunk.Reader.HexDump(true)

	// Always big endian?
	r.Seek(0, 0)
	err = chunk.Read(castChunk, binary.BigEndian)
	if err != nil {
		return chunk, err
	}

	return chunk, nil
}

func (c *XmedChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "XMED")
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "XMED")
	}

	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists

	{
		outputFile := filepath.Join(outputFolder, name+"."+c.Extension)

		targetFile, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("could not create file...: %v\n", err)
			return
		}
		defer targetFile.Close()

		// Write the file to disk.
		_, err = targetFile.Write(c.Data)
		if err != nil {
			fmt.Printf("could not write file...: %v\n", err)
			panic(err)
		}

		err = targetFile.Sync()
		if err != nil {
			fmt.Printf("could not sync file...: %v\n", err)
			panic(err)
		}
	}

	{
		metaOutputFile := filepath.Join(outputFolder, name+"."+c.Extension+".json")

		metaTargetFile, err := os.Create(metaOutputFile)
		if err != nil {
			fmt.Printf("could not create metafile...: %v\n", err)
			return
		}
		defer metaTargetFile.Close()

		// Write the file to disk.
		_, err = metaTargetFile.Write(c.Meta)
		if err != nil {
			fmt.Printf("could not write meta file...: %v\n", err)
			panic(err)
		}

		err = metaTargetFile.Sync()
		if err != nil {
			fmt.Printf("could not sync meta file...: %v\n", err)
			panic(err)
		}
	}
}

func (c *XmedChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
