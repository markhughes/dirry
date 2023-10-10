package chunks

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/errors"
	"github.com/markhughes/dirry/internal/members"
	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
	"golang.org/x/text/encoding/charmap"
)

type Rect struct {
	Top, Left, Bottom, Right int16
}

type CastMember interface {
	ToJson() (string, error)
	FromBytes(b []byte, v version.Version, flags uint8) error
}

type CastType int

const (
	Bitmap CastType = iota + 1
	FilmLoop
	StyledText
	Palette
	Picture
	Sound
	Button
	Shape
	Movie
	DigitalVideo
	Script
	Text
	OLE
	Transition
	Xtra
)

func (t CastType) String() string {
	switch t {
	case Bitmap:
		return "Bitmap"
	case FilmLoop:
		return "FilmLoop"
	case StyledText:
		return "StyledText"
	case Palette:
		return "Palette"
	case Picture:
		return "Picture"
	case Sound:
		return "Sound"
	case Button:
		return "Button"
	case Shape:
		return "Shape"
	case Movie:
		return "Movie"
	case DigitalVideo:
		return "DigitalVideo"
	case Script:
		return "Script"
	case Text:
		return "Text"
	case OLE:
		return "OLE"
	case Transition:
		return "Transition"
	case Xtra:
		return "Xtra"
	default:
		return fmt.Sprintf("CastType: %d", int(t))
	}
}

type CommonMemberProperties struct {
	ScriptText         string
	Name               string
	FilePath           string
	FileName           string
	FileType           string
	XtraGUID           [16]byte
	XtraName           string
	RegistrationPoints []int32
	ClipboardFormat    string
	CreationDate       int32
	ModifiedDate       int32
	ModifiedBy         string
	Comments           string
	ImageCompression   int
	ImageQuality       int
}

type CastChunk struct {
	Chunk          *BinaryChunk
	Type           CastType
	CastDataLength uint32
	CastInfoLength uint32

	Properties CommonMemberProperties

	Member CastMember

	BasicContent []string
	ExtraContent []string

	PropertyOffsets []int32
}

func (chunk *CastChunk) Save(projectName string, name string, pkg string) {
	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "converted", "CASt", name)
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "converted", "CASt", name)
	}

	os.MkdirAll(outputFolder, os.ModePerm)

	// save name to name.txt
	if chunk.Properties.Name != "" {
		var fileName = filepath.Join(outputFolder, "label.txt")
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("could not create file...: %v\n", err)
			return
		}
		// write chunk.Name to file

		_, err = file.WriteString(chunk.Properties.Name)
		if err != nil {
			fmt.Printf("could not write file...: %v\n", err)
			panic(err)
		}

		file.Close()
	}

	// save script script.ls
	if chunk.Properties.ScriptText != "" {
		var fileName = filepath.Join(outputFolder, "script.ls")
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("could not create file...: %v\n", err)
			return
		}

		// write chunk.Member.ScriptText to file
		_, err = file.WriteString(chunk.Properties.ScriptText)
		if err != nil {
			fmt.Printf("could not write file...: %v\n", err)
			panic(err)
		}

		file.Close()
	}

}

// TODO: very broken atm
func ReadCastChunkRaw(r *bytes.Reader, v version.Version, endian binary.ByteOrder, isAfterburner bool) (*CastChunk, error) {

	var err error
	chunk := &CastChunk{}
	chunk.Chunk, err = FromBinaryAtPartHeadless(r, "CASt", endian)
	if err != nil {
		return nil, err
	}

	headerData := make([]byte, 0)

	basicData := make([]byte, 0)

	dataType, err := chunk.Chunk.ReadInt32(endian)
	if err != nil {
		return nil, fmt.Errorf("error reading cast data type: %v", err)
	}

	var headerSize int32
	var additionalSize int32

	if (int64(dataType) & 0xFFFFFF00) != 0 {
		utils.DebugMsg("CASt", "Director 4 cast chunk detected")
		// Director 4
		chunk.Chunk.Seek(0, io.SeekCurrent)

		headerSize16, err := chunk.Chunk.ReadInt16(endian)
		if err != nil {
			return nil, fmt.Errorf("[d4] error reading cast header size: %v", err)
		}

		headerSize = int32(headerSize16)

		additionalSize, err = chunk.Chunk.ReadInt32(endian)
		if err != nil {
			return nil, fmt.Errorf("[d4] error reading cast additional size: %v", err)
		}

		dataType, err = chunk.Chunk.ReadInt32(endian)
		if err != nil {
			return nil, fmt.Errorf("[d4] error reading cast data type: %v", err)
		}

		if dataType < 1 || dataType > 15 {
			//
			chunk.Chunk.Seek(-4, io.SeekCurrent)
			dataTypeByte, err := chunk.Chunk.ReadUInt8()
			if err != nil {
				return nil, fmt.Errorf("[d4] error reading cast data type: %v", err)
			}

			dataType2 := int32(dataTypeByte)

			if dataType2 < 1 || dataType2 > 15 {

				return nil, fmt.Errorf("[d4] invalid cast data type, tried: %d and %d", dataType, dataType2)
			}

			dataType = dataType2
		}

		if headerSize > 0 {
			headerData = make([]byte, headerSize-1)
			_, err = chunk.Chunk.Read(headerData)
			if err != nil {
				return nil, fmt.Errorf("[d4] error reading cast header data: %v", err)
			}
		}

		if additionalSize > 0 {
			basicData = make([]byte, additionalSize)
			_, err = chunk.Chunk.Read(basicData)
			if err != nil {
				return nil, fmt.Errorf("[d4] error reading cast basic data: %v", err)
			}
		}

	} else {
		// Director 5+
		utils.DebugMsg("CASt", "Director 5+ cast chunk detected")

		if dataType < 1 || dataType > 15 {
			return nil, fmt.Errorf("[d4] invalid cast data type: %d", dataType)
		}

		additionalSize, err = chunk.Chunk.ReadInt32(endian)
		if err != nil {
			return nil, fmt.Errorf("[d5] error reading cast additional size: %v", err)
		}

		headerSize, err = chunk.Chunk.ReadInt32(endian)
		if err != nil {
			return nil, fmt.Errorf("[d5] error reading cast header size: %v", err)
		}

		if additionalSize > 0 {
			basicData = make([]byte, additionalSize)
			_, err = chunk.Chunk.Read(basicData)
			if err != nil {
				return nil, fmt.Errorf("[d5] error reading cast basic data: %v", err)
			}
		}

		if headerSize > 0 {
			headerData = make([]byte, headerSize)
			_, err = chunk.Chunk.Read(headerData)
			if err != nil {
				return nil, fmt.Errorf("[d5] error reading cast header data: %v", err)
			}
		}

	}

	if len(basicData) > 0 {
		// Parse main data
		contentMarker := binary.BigEndian.Uint32(basicData[:4])
		// log.Printf("contentMarker? = %08x\n", contentMarker)

		if contentMarker != 0x00000014 {
			// Just return current chunk as-is and let the caller handle it
			return chunk, fmt.Errorf("content marker is not 0x00000014")
		}

		basicData00 := binary.BigEndian.Uint32(basicData[4:8])
		// log.Printf("basicData00 = %08x\n", basicData00)

		basicData01 := binary.BigEndian.Uint32(basicData[8:12])
		// log.Printf("basicData01 = %08x\n", basicData01)

		basicData02 := binary.BigEndian.Uint32(basicData[12:16])
		// log.Printf("basicData02 = %08x\n", basicData02)

		basicData03 := binary.BigEndian.Uint32(basicData[16:20])
		// log.Printf("basicData03 = %08x\n", basicData03)

		// fmt.Printf("basicData = %v\n", hex.Dump(basicData[21:]))
		chunk.BasicContent = make([]string, 0)

		chunk.BasicContent = append(chunk.BasicContent, fmt.Sprintf("0x%08x", basicData00&0xFFFFFFFF))
		chunk.BasicContent = append(chunk.BasicContent, fmt.Sprintf("0x%08x", basicData01&0xFFFFFFFF))
		chunk.BasicContent = append(chunk.BasicContent, fmt.Sprintf("0x%08x", basicData02&0xFFFFFFFF))
		chunk.BasicContent = append(chunk.BasicContent, fmt.Sprintf("0x%08x", basicData03&0xFFFFFFFF))

		chunk.Properties = CommonMemberProperties{}
		// chunk.Properties.ReadProperties(basicData)

		nstruct := binary.BigEndian.Uint16(basicData[20:22])
		// log.Printf("number of structures contained = %d\n", nstruct)

		if nstruct > 0 {
			reg, err := regexp.Compile(`[^A-Za-z0-9\-_. ]+`)

			if err != nil {
				log.Fatal(err)
			}
			chunk.ExtraContent = make([]string, 0)

			chunk.PropertyOffsets = make([]int32, nstruct+1)
			for i := range chunk.PropertyOffsets {
				chunk.PropertyOffsets[i] = int32(binary.BigEndian.Uint32(basicData[22+i*4 : 26+i*4]))
				// log.Printf("stindx[%d] = %08x\n", i, chunk.PropertyOffsets[i])
			}

			var basicDataReader = bytes.NewReader(basicData)
			// 22 + nstruct*4
			var what = make([]byte, 22+nstruct*4)
			basicDataReader.Read(what)

			for i := 0; i < int(nstruct); i++ {
				stlen := int(chunk.PropertyOffsets[i+1] - chunk.PropertyOffsets[i])
				if stlen > 0 {
					data := make([]byte, stlen)
					basicDataReader.Read(data)

					chunk.ExtraContent = append(chunk.ExtraContent, string(data))
					value, err := charmap.ISO8859_1.NewDecoder().String(string(data))
					if err != nil {
						return nil, fmt.Errorf("error decoding string in properties: %v", err)
					}

					switch i {
					case 0:
						chunk.Properties.ScriptText = value
					case 1:

						chunk.Properties.Name = reg.ReplaceAllString(value, "_")

					case 2:
						chunk.Properties.FilePath = value

					case 3:
						chunk.Properties.FileName = value

					case 4:
						chunk.Properties.FileType = value

					case 9:
						var guid [16]byte
						copy(guid[:], data[:16])
						chunk.Properties.XtraGUID = guid

					case 10:
						chunk.Properties.XtraName = value

					case 12:
						for i := 0; i < stlen/4; i++ {
							var point int32
							binary.Read(bytes.NewReader(basicData[chunk.PropertyOffsets[i]+22:chunk.PropertyOffsets[i]+26]), binary.BigEndian, &point)
							chunk.Properties.RegistrationPoints = append(chunk.Properties.RegistrationPoints, point)
						}

					case 16:
						chunk.Properties.ClipboardFormat = value

					case 17:
						var date int32
						binary.Read(bytes.NewReader(basicData[chunk.PropertyOffsets[i]+22:chunk.PropertyOffsets[i]+26]), binary.BigEndian, &date)
						chunk.Properties.CreationDate = date * 1000

					case 18:
						var date int32
						binary.Read(bytes.NewReader(basicData[chunk.PropertyOffsets[i]+22:chunk.PropertyOffsets[i]+26]), binary.BigEndian, &date)

						chunk.Properties.ModifiedDate = date * 1000

					case 19:
						chunk.Properties.ModifiedBy = value

					case 20:
						chunk.Properties.Comments = value

					case 21:
						var compression int32
						binary.Read(bytes.NewReader(basicData[chunk.PropertyOffsets[i]+22:chunk.PropertyOffsets[i]+26]), binary.BigEndian, &compression)
						chunk.Properties.ImageCompression = int(compression)

					case 22:
						var quality int32
						binary.Read(bytes.NewReader(basicData[chunk.PropertyOffsets[i]+22:chunk.PropertyOffsets[i]+26]), binary.BigEndian, &quality)
						chunk.Properties.ImageQuality = int(quality)

					}
				} else {
					chunk.ExtraContent = append(chunk.ExtraContent, "")

				}
			}
		}
	}

	chunk.Type = CastType(dataType)

	chunk.CastDataLength = uint32(additionalSize)
	chunk.CastInfoLength = uint32(headerSize)

	if len(headerData) > 0 {
		switch CastType(dataType) {
		case Script:
			scriptMember := &members.MemberScript{}
			chunk.Member = scriptMember
		case Bitmap, OLE:
			bitmapMember := &members.MemberBitmap{}
			chunk.Member = bitmapMember
		case Xtra:
			xtraMember := &members.MemberXtra{}
			chunk.Member = xtraMember
		default:
			return chunk, &errors.UnhandledCastTypeError{CastTypeName: CastType(dataType).String(), CastType: dataType}
		}

		if chunk.Member != nil {
			chunk.Member.FromBytes(headerData, v, 0)
		}
	}

	return chunk, nil

}

func (c *CastChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
