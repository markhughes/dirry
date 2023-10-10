package shockwave

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/utils"
)

const (
	imapPos    = 0xC
	intMmapPos = 0x18
	mmapPos    = 0x2C
)

// create file strucutre
type File struct {
	Tag     string
	Size    uint32
	Offset  uint32
	Content []byte
}

/**
 * Thanks to n0samu for documenting this, was hard to find info about this online.
 * https://github.com/n0samu/director-files-extract
 */
func (shockwave *Shockwave) ExtractExe(fBytes []byte) (shockwavePath []ShockwaveFile, err error) {

	var outFiles = make([]ShockwaveFile, 0)

	// Find the APPL/LPPA file
	winFile := regexp.MustCompile(`XFIR.{4}LPPA`).FindIndex(fBytes)
	macFile := regexp.MustCompile(`RIFX.{4}APPL`).FindIndex(fBytes)

	var off int
	if winFile != nil {
		off = winFile[0]
	} else if macFile != nil {
		off = macFile[0]
	} else {
		// TODO: seems like D4 EXEs of a NE (not PE) are a little different.. not sure how to handle this
		return outFiles, fmt.Errorf("not a director application")
	}

	utils.InfoMsg("exe", "Confirmed Director file at %d", off)
	f := bytes.NewReader(fBytes[off:])

	sig := make([]byte, 4)
	f.Read(sig)
	var signature = string(sig)
	var endian binary.ByteOrder
	switch signature {
	case "XFIR", "FFIR":
		endian = binary.LittleEndian
	case "RIFX", "RIFF":
		endian = binary.BigEndian
	default:
		return outFiles, fmt.Errorf("cannot determine codec from signature")
	}

	utils.DebugMsg("exe", "Executable director file signature: %s", signature)

	// IMAP
	f.Seek(imapPos, io.SeekStart)
	tag, err := utils.ReadString(f, 4, endian == binary.LittleEndian)
	if err != nil {
		return outFiles, fmt.Errorf("error reading imap tag: %s", err)
	}

	if tag != "imap" {
		return outFiles, fmt.Errorf("expected 'imap' tag not found, got %s", tag)
	}

	f.Seek(8, io.SeekCurrent)

	// Find the mmap offset in the imap chunk
	mmapOff, err := utils.ReadUInt32(f, endian)
	if err != nil {
		return outFiles, fmt.Errorf("error reading mmap offset: %s", err)
	}

	mmapOff -= mmapPos

	// MMAP
	_, err = f.Seek(int64(mmapPos), io.SeekStart)
	if err != nil {
		return outFiles, fmt.Errorf("error seeking to mmap: %s", err)
	}

	tag, err = utils.ReadString(f, 4, endian == binary.LittleEndian)
	if err != nil {
		return outFiles, fmt.Errorf("error reading mmap tag: %s", err)
	}

	if tag != "mmap" {
		return outFiles, fmt.Errorf("expected 'mmap' tag not found")
	}

	_, err = f.Seek(int64(mmapPos+10), io.SeekStart)
	if err != nil {
		return outFiles, fmt.Errorf("error seeking to mmap resource count: %s", err)
	}

	mmapResLen, err := utils.ReadInt16(f, endian)
	if err != nil {
		return outFiles, fmt.Errorf("error reading mmap resource length: %s", err)
	}

	_, err = f.Seek(int64(mmapPos+0x10), io.SeekStart)
	if err != nil {
		return outFiles, fmt.Errorf("error seeking to mmap resource count: %s", err)
	}

	mmapResCount, err := utils.ReadUInt32(f, endian)
	if err != nil {
		return outFiles, fmt.Errorf("error reading mmap resource count: %s", err)
	}

	// MMAP resources
	mmapRessPos := mmapPos + 32
	_, err = f.Seek(int64(mmapRessPos+8), io.SeekStart)
	if err != nil {
		return outFiles, fmt.Errorf("error seeking to mmap resource count: %s", err)
	}

	REL, err := utils.ReadUInt32(f, endian)
	if err != nil {
		return outFiles, fmt.Errorf("error reading REL: %s", err)
	}

	// This seems heavily inspired by the mac clipboard names
	resources := make([]File, mmapResCount)
	for i := range resources {
		f.Seek(int64(i*int(mmapResLen))+int64(mmapRessPos), io.SeekStart)
		tag, err := utils.ReadString(f, 4, endian == binary.LittleEndian)
		if err != nil {
			return outFiles, err
		}

		size, err := utils.ReadUInt32(f, endian)
		if err != nil {
			return outFiles, err
		}

		resoff, err := utils.ReadUInt32(f, endian)
		if err != nil {
			return outFiles, err
		}

		// size and offset corrections
		size += 8
		if resoff != 0 {
			resoff -= REL
		}

		f.Seek(int64(resoff), io.SeekStart)
		fileContent := make([]byte, size)
		f.Read(fileContent)

		resources[i] = File{
			Tag:     tag,
			Size:    size,
			Offset:  resoff,
			Content: fileContent,
		}

	}

	var files []File = make([]File, 0)
	// Files
	{
		for i := range resources {
			if resources[i].Tag == "File" {
				files = append(files, resources[i])
			}
		}
	}

	// Read Dict
	{
		var dict File
		for i := range resources {
			if resources[i].Tag == "Dict" {
				dict = resources[i]
				break
			}
		}

		if dict.Tag == "" {
			return outFiles, fmt.Errorf("no Dict file found")
		}

		dictReader, err := parseDict(dict.Content, endian)
		if err != nil {
			return outFiles, err
		}

		for i := range dictReader {
			var currentFile = files[i]
			var currentPath = dictReader[i]

			var directory = filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "extracted")
			if strings.Contains(currentPath, ".x32") || strings.Contains(currentPath, ".x16") {
				directory = filepath.Join(directory, "Xtras")
			}
			os.MkdirAll(directory, os.ModePerm)

			var file = filepath.Join(directory, path.Base(strings.ReplaceAll(currentPath, "\\", "/")))
			ioutil.WriteFile(file, currentFile.Content, 0644)

			fmt.Printf("Found file: %s at %d\n", file, int64(currentFile.Offset)+int64(off))

			var sfile = ShockwaveFile{
				Path:        file,
				MinusOffset: int64(currentFile.Offset) + int64(off),
			}
			// check if the content starts with RIFX or XFIR
			if bytes.HasPrefix(currentFile.Content, []byte("RIFX")) || bytes.HasPrefix(currentFile.Content, []byte("XFIR")) {
				outFiles = append(outFiles, sfile)

			}
		}

	}

	for i := range resources {
		var directory = filepath.Join(consts.PathDump, filepath.Base(shockwave.FilePath), "exe_resources")
		os.MkdirAll(directory, os.ModePerm)

		var file = filepath.Join(directory, fmt.Sprintf("%d_%s", i, resources[i].Tag))
		ioutil.WriteFile(file, resources[i].Content, 0644)

	}

	if len(outFiles) == 0 {
		return outFiles, fmt.Errorf("no Director files found")
	}

	return outFiles, nil

}

func parseDict(data []byte, endian binary.ByteOrder) ([]string, error) {
	var err error
	r := bytes.NewReader(data[8:])
	var toclen uint32
	binary.Read(r, binary.LittleEndian, &toclen)
	if toclen > 0x10000 {
		// Win16 EXEs swap endianness after the tag size
		if endian == binary.LittleEndian {
			endian = binary.BigEndian
		} else {
			endian = binary.LittleEndian
		}
		r.Seek(0, 0)
		binary.Read(r, binary.LittleEndian, &toclen)
	}
	r.Seek(0x10, 0)
	var len_names uint32
	err = binary.Read(r, endian, &len_names)
	if err != nil {
		return nil, err
	}
	r.Seek(0x18, 0)
	// r.Read(toclen)
	r.Seek(int64(toclen), 1)
	unk1, err := utils.ReadInt16(r, endian)
	if err != nil {
		return nil, err
	}
	r.Seek(int64(unk1-0x12), 1)
	names := make([]string, len_names)
	for i := range names {
		var lname uint32
		binary.Read(r, endian, &lname)
		fname := make([]byte, lname)
		r.Read(fname)
		// assert lname == len(fname)
		// r.Read(-lname % 4)
		r.Seek(int64((-lname)%4), 1)
		names[i] = string(fname)
	}
	return names, nil
}
