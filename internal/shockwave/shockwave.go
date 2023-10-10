package shockwave

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/markhughes/dirry/internal/chunks"
	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/utils"
	"github.com/markhughes/dirry/internal/version"
)

type Shockwave struct {
	ID     string
	Length uint32
	Endian binary.ByteOrder

	FilePath string

	FileReader  *os.File
	BytesReader *bytes.Reader

	DirOffset int64

	Codec Codec

	ChunkMap ChunkMap

	Version version.Version

	ProjectName string
	PkgName     string

	Casts map[int32]*chunks.CastChunk
	Fonts map[uint32]*chunks.Font
}

type ShockwaveFile struct {
	Path        string
	MinusOffset int64
}

func (shockwave *Shockwave) Init() {
	shockwave.Casts = make(map[int32]*chunks.CastChunk)
	shockwave.Fonts = make(map[uint32]*chunks.Font)
	shockwave.ProjectName = filepath.Base(shockwave.FilePath)
}

func (shockwave *Shockwave) GetReader() io.ReadSeeker {
	if shockwave.FileReader != nil {
		return shockwave.FileReader
	}

	if shockwave.BytesReader != nil {
		return shockwave.BytesReader
	}

	panic("no reader")
}

func (shockwave *Shockwave) read() (expanded []ShockwaveFile, openError error) {
	// Let's start by reading in the header
	var id [4]byte
	if _, err := io.ReadFull(shockwave.GetReader(), id[:]); err != nil {
		return nil, err
	}

	// Ok lets check are the first two characters in `id` MZ
	if id[0] == 'M' && id[1] == 'Z' {
		// Seek to start
		shockwave.GetReader().Seek(0, io.SeekStart)

		// read entire contents into memory
		buf := new(bytes.Buffer)
		buf.ReadFrom(shockwave.GetReader())

		// TODO: support WASM
		// Try to extract from the EXE
		filePaths, err := shockwave.ExtractExe(buf.Bytes())
		if err != nil {
			return nil, err
		}

		return filePaths, nil

	} else {
		// It's something else...
		shockwave.ID = string(id[:])
	}

	utils.InfoMsg("shockwave", "ID: %s\n", shockwave.ID)
	if shockwave.ID == "XFIR" {
		// XFIR is little endian (but honestly it doesn't seem to apply everywhere?)
		shockwave.Endian = binary.LittleEndian
	} else if shockwave.ID == "RIFX" {
		// RIFX is big endian
		shockwave.Endian = binary.BigEndian
	} else {
		return nil, fmt.Errorf("unknown shockwave type: %s", shockwave.ID)
	}

	if err := binary.Read(shockwave.GetReader(), shockwave.Endian, &shockwave.Length); err != nil {
		utils.ErrorMsg("shockwave", "Error reading length: %s", err)
	}

	utils.InfoMsg("shockwave", "Length: %d", shockwave.Length)

	codecName := make([]byte, 4) // This creates a slice of 4 bytes
	if _, err := io.ReadFull(shockwave.GetReader(), codecName); err != nil {
		utils.ErrorMsg("shockwave", "Error reading codec name: %s", err)
		return
	}
	if shockwave.Endian == binary.LittleEndian {
		codecName = utils.ReverseBytes(codecName)
	}

	utils.InfoMsg("shockwave", "Codec Name: %s", string(codecName[:]))

	var ok bool
	shockwave.Codec, ok = CodecByName(string(codecName[:]))
	if !ok {
		return nil, fmt.Errorf("unknown codec: %s", string(codecName[:]))
	}

	// Codec can determine how file is read e.g. afterburner is different

	utils.InfoMsg("shockwave", "Codec Name: %s", codecName[:])
	utils.InfoMsg("shockwave", "Codec Name: %s", shockwave.Codec.Name)
	utils.DebugMsg("shockwave", "Codec Type: %s", shockwave.Codec.Type)

	if shockwave.Codec.Type == Afterburner {
		ParseAfterburner(shockwave)
	} else {
		ParseStandard(shockwave)
	}

	return nil, nil

}

func (shockwave *Shockwave) OpenContent(filePath string, content []byte) (expanded []ShockwaveFile, openError error) {
	utils.InfoMsg("shockwave", "Opening File: %s\n", filePath)

	shockwave.FilePath = filePath

	shockwave.Init()

	shockwave.BytesReader = bytes.NewReader(content)

	return shockwave.read()
}

/**
 * Opens a shockwave file, or returns a list of files to open separately.
 */
func (shockwave *Shockwave) Open(filePath string) (expanded []ShockwaveFile, openError error) {
	utils.InfoMsg("shockwave", "Opening File: %s\n", filePath)
	var err error
	shockwave.FilePath = filePath
	shockwave.FileReader, err = os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}

	shockwave.Init()

	return shockwave.read()
}

func (shockwave *Shockwave) Zip() error {
	os.MkdirAll(consts.PathZips, os.ModePerm) // Ensure the output directory exists

	zipFileName := path.Join(consts.PathZips, shockwave.ProjectName+".zip")

	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("error creating zip file: %s", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	resources := shockwave.ChunkMap.GetAllResources()
	for _, resource := range resources {
		if resource.ChunkType == "RIFX" || resource.ChunkType == "XFIR" {
			continue
		}
		// create a file header
		header := &zip.FileHeader{
			Name: fmt.Sprintf("%s/%d", resource.ChunkType, resource.ResourceId),
			// Set the compression method to Deflate (the most common method)
			Method: zip.Deflate,
		}

		header.SetMode(os.ModePerm)

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// write the resource binary data to the zip
		_, err = writer.Write(resource.Binary)
		if err != nil {
			return err
		}
	}

	header := &zip.FileHeader{
		Name:   ("map.json"),
		Method: zip.Deflate,
	}

	header.SetMode(os.ModePerm)

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	var jsonMapArray []interface{}

	for _, chunk := range shockwave.ChunkMap.GetAllResources() {
		var item = make(map[string]interface{})
		if chunk.ChunkType == "RIFX" || chunk.ChunkType == "XFIR" {
			continue
		}

		item["chunkType"] = chunk.ChunkType
		item["resourceId"] = chunk.ResourceId
		item["offset"] = chunk.Offset
		item["name"] = chunk.Name
		jsonMapArray = append(jsonMapArray, item)

	}

	data, err := json.MarshalIndent(jsonMapArray, "", " ")
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil

}

func (shockwave *Shockwave) Close() {
	if shockwave.FileReader != nil {
		shockwave.FileReader.Close()
	}
}

func (shockwave *Shockwave) IsAfterburner() bool {
	return shockwave.Codec.Type == Afterburner
}

/**
 * binary dump
 */
func (shockwave *Shockwave) Dump(offset int64, length int64, targetDir string, targetFileName string) {
	targetDir = utils.CleanString(targetDir)
	targetFileName = utils.CleanString(targetFileName)
	if (targetDir) == "" {
		targetDir = "empty_name"
	}

	utils.DebugMsg("shockwave", "Dumping %d bytes from offset %d to %s %s\n", length, offset, targetDir, targetFileName)

	if length < 0 {
		utils.DebugMsg("shockwave", "Length is negative, skipping\n")
		return
	}
	if (offset + length) > int64(shockwave.Length) {
		utils.DebugMsg("shockwave", "Length is greater than file size, skipping\n")
		return
	}
	if offset < 0 {
		utils.DebugMsg("shockwave", "Offset is negative, skipping\n")
		return
	}

	outputFolder := filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "binary", ""+string(targetDir))
	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists
	outputFile := filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "binary", ""+string(targetDir), ""+string(targetFileName))

	// seek to wherever it is that we want to read
	_, err := shockwave.GetReader().Seek(offset, 0)
	if err != nil {
		panic(err) // panic as in what i do daily
	}

	// Create two files: <name>.chunk and <name>.bin
	chunkFile, err := os.Create(outputFile + ".chunk")
	if err != nil {
		fmt.Printf("Error creating file: %s\n", outputFile+".chunk")
		return
	}
	defer chunkFile.Close()

	binFile, err := os.Create(outputFile + ".bin")
	if err != nil {
		fmt.Printf("Error creating file: %s\n", outputFile+".bin")
		return
	}
	defer binFile.Close()

	buf := make([]byte, length+8) // +8 as this will have the header in it too...
	_, err = shockwave.GetReader().Read(buf)
	if err != nil {
		panic(err)
	}

	_, err = chunkFile.Write(buf)
	if err != nil {
		panic(err)
	}

	if length > 0 {

		_, err = binFile.Write(buf[8:]) // we skip the first 8 bytes
		if err != nil {
			panic(err)
		}
		err = binFile.Sync()
		if err != nil {
			panic(err)
		}
	}

	err = chunkFile.Sync()
	if err != nil {
		panic(err)
	}

}
