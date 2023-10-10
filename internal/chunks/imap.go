package chunks

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

type ImapChunk struct {
	Signature            string
	Length               int32
	MemoryMapCount       int32
	MemoryMapOffset      int32
	MemoryMapFileVersion int32
	Reserved             int16
	Unknown              int16
	Reserved2            int32
}

func ReadImapChunk(r io.Reader, endian binary.ByteOrder) (*ImapChunk, error) {
	c := &ImapChunk{}
	var err error

	pos, err := r.(io.Seeker).Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	var data = make([]byte, 32)
	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}

	_, err = r.(io.Seeker).Seek(pos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	c.Signature, err = utils.ReadString(r, 4, endian == binary.LittleEndian)
	if err != nil {
		return nil, err
	}

	if c.Signature != "imap" {
		return nil, fmt.Errorf("invalid imap signature: %s", c.Signature)
	}

	c.Length, err = utils.ReadInt32(r, endian)
	if err != nil {
		return nil, err
	}

	// count    offset   version   reserved  unknown  reserved2
	// 01000000 2C000000 42070000 00000000 00000000 00000000
	//        1       44     1858        0        0        0

	// 01000000 AC002A00 42070000 00000000 00000000 00000000
	// 	  1        44044     1858        0        0        0
	utils.DebugMsg("imap", "imap length: %d\n", c.Length)

	c.MemoryMapCount, err = utils.ReadInt32(r, endian)
	if err != nil {
		return nil, err
	}

	utils.DebugMsg("imap", "imap memory map count: %d\n", c.MemoryMapCount)

	c.MemoryMapOffset, err = utils.ReadInt32(r, endian)
	if err != nil {
		return nil, err
	}

	utils.DebugMsg("imap", "imap memory map offset: %d\n", c.MemoryMapOffset)

	c.MemoryMapFileVersion, err = utils.ReadInt32(r, endian)
	if err != nil {
		return nil, fmt.Errorf("error reading imap memory map file version: %s", err)
	}

	c.Reserved, err = utils.ReadInt16(r, endian)
	if err != nil {
		return nil, fmt.Errorf("error reading imap reserved: %s", err)
	}

	c.Unknown, err = utils.ReadInt16(r, endian)
	if err != nil {
		return nil, fmt.Errorf("error reading imap unknown: %s", err)
	}

	c.Reserved2, err = utils.ReadInt32(r, endian)
	if err != nil {
		return nil, fmt.Errorf("error reading imap reserved2: %s", err)
	}

	return c, nil
}

func (c *ImapChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
