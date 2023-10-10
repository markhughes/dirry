package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"regexp"
)

func ReadString(r io.Reader, n int, reverse bool) (string, error) {
	buf := make([]byte, n)

	err := binary.Read(r, binary.LittleEndian, buf)
	if err != nil {
		return "", err
	}

	if reverse {
		buf = ReverseBytes(buf)
	}
	return string(buf), nil
}

func ReadVarInt(r io.Reader) (val uint32, bytesRead int, err error) {
	var b byte

	for {
		b = 0
		// endian doesn't matter here as we're only reading one byte
		err = binary.Read(r, binary.LittleEndian, &b)
		if err != nil {
			return 0, 0, err
		}
		bytesRead += 1

		val = (val << 7) | uint32(b&0x7f)

		if b>>7 == 0 {
			break
		}
	}

	return val, bytesRead, nil
}

func ReadUInt32(r io.Reader, endian binary.ByteOrder) (uint32, error) {
	var n uint32
	err := binary.Read(r, endian, &n)
	return n, err
}

func ReadInt32(r io.Reader, endian binary.ByteOrder) (int32, error) {
	var n int32
	err := binary.Read(r, endian, &n)
	return n, err
}

func ReadInt16(r io.Reader, endian binary.ByteOrder) (int16, error) {
	var n int16
	err := binary.Read(r, endian, &n)
	return n, err
}
func ReadUInt16(r io.Reader, endian binary.ByteOrder) (uint16, error) {
	var n uint16
	err := binary.Read(r, endian, &n)
	return n, err
}

func ReadInt8(r io.Reader) (int8, error) {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return 0, err
	}
	return int8(b[0]), nil
}

func ReverseBytes(bytes []byte) []byte {
	length := len(bytes)
	reversed := make([]byte, length)
	for i, b := range bytes {
		reversed[length-1-i] = b
	}
	return reversed
}

func Hexdump(r *os.File, size int64, print bool) (string, error) {

	buf := make([]byte, size)
	n, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	buf = buf[:n]
	hexdump := hex.Dump(buf)

	if print {
		fmt.Println(hexdump)
	}

	return hexdump, nil
}

type Rect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
	Width  int16
	Height int16
}

func ReadRect(r io.Reader, endian binary.ByteOrder) (Rect, error) {
	var rect Rect
	rect.Top, _ = ReadInt16(r, endian)
	rect.Left, _ = ReadInt16(r, endian)
	rect.Bottom, _ = ReadInt16(r, endian)
	rect.Right, _ = ReadInt16(r, endian)

	rect.Width = rect.Right - rect.Left
	rect.Height = rect.Bottom - rect.Top

	return rect, nil
}

/*
	rect.left = stream->readSint16BE();
	rect.top = stream->readSint16BE();
	rect.right = stream->readSint16BE();
	rect.bottom = stream->readSint16BE();

*/

// TODO: move to binary chunk struct
func ReadUInt8(r io.Reader) (uint8, error) {
	var res uint8
	err := binary.Read(r, binary.BigEndian, &res)
	return res, err
}

func CleanString(str string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(str, "")
}
