package binary_reader

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/markhughes/dirry/internal/utils"
)

type BinaryReader struct {
	bytes       []byte
	bytesReader *bytes.Reader
	Length      int32
}

func NewBinaryReader(b []byte, length int32) (*BinaryReader, error) {
	reader := &BinaryReader{}
	err := reader.setBinary(b)
	if err != nil {
		return nil, err
	}

	reader.Length = length

	return reader, nil
}

func (br *BinaryReader) Size() int {
	return len(br.bytes)
}

func (br *BinaryReader) setBinary(b []byte) error {
	br.bytes = b

	br.bytesReader = bytes.NewReader(br.bytes)
	return nil
}

func (br *BinaryReader) ReadAllBytes() ([]byte, error) {
	if br.bytes == nil {
		return nil, fmt.Errorf("bytes is nil")
	}
	return br.bytes, nil
}

func (br *BinaryReader) ReadVarInt() (val uint32, bytesRead int, err error) {
	return utils.ReadVarInt(br.bytesReader)
}

func (br *BinaryReader) ReadInt32(endian binary.ByteOrder) (int32, error) {
	return utils.ReadInt32(br.bytesReader, endian)
}

func (br *BinaryReader) ReadInt16(endian binary.ByteOrder) (int16, error) {
	return utils.ReadInt16(br.bytesReader, endian)
}

func (br *BinaryReader) ReadInt8() (int8, error) {
	return utils.ReadInt8(br.bytesReader)
}

func (br *BinaryReader) ReadUInt32(endian binary.ByteOrder) (uint32, error) {
	return utils.ReadUInt32(br.bytesReader, endian)
}

func (br *BinaryReader) ReadUInt16(endian binary.ByteOrder) (uint16, error) {
	var b [2]byte
	_, err := io.ReadFull(br.bytesReader, b[:])
	if err != nil {
		return 0, err
	}
	return endian.Uint16(b[:]), nil
}

func (br *BinaryReader) ReadUInt8() (uint8, error) {
	var b [1]byte
	_, err := io.ReadFull(br.bytesReader, b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (br *BinaryReader) ReadString(length int, endian binary.ByteOrder) (string, error) {
	return utils.ReadString(br.bytesReader, length, endian == binary.LittleEndian)
}

func (br *BinaryReader) ReadUByte() (int, error) {
	b := make([]byte, 1)

	_, err := br.bytesReader.Read(b)
	if err != nil {
		return 0, err
	}

	return int(b[0]), nil

}

func (br *BinaryReader) ReadBytes(length int) ([]byte, error) {
	b := make([]byte, length)

	_, err := br.bytesReader.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (br *BinaryReader) ReadRect(endian binary.ByteOrder) (utils.Rect, error) {
	return utils.ReadRect(br.bytesReader, endian)
}

// Utils

func (br *BinaryReader) Seek(offset int64, whence int) (int64, error) {
	return br.bytesReader.Seek(offset, whence)
}

func (br *BinaryReader) Pos() int64 {
	pos, _ := br.bytesReader.Seek(0, io.SeekCurrent)
	return pos
}

func (br *BinaryReader) HexDump(print bool) string {
	var data = hex.Dump(br.bytes)

	if print {
		utils.DebugMsg("binary_reader", "HexDump:\n%s", data)
	}

	return data
}

// Deprecated: this is unsafe.
func (br *BinaryReader) GetUnsafeBytesReader() *bytes.Reader {
	return br.bytesReader
}

func (br *BinaryReader) GetBytes() []byte {
	return br.bytes
}
