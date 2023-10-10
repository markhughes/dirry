package chunks

/*
This is the old binary reader that was used in the first build of dirry, ideally we should move away from this
*/
import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/markhughes/dirry/internal/consts"
	"github.com/markhughes/dirry/internal/utils"
)

// Deprecated: move away for hacked together BinaryChunk
type BinaryChunk struct {
	Type          string
	StartPosition int64
	Length        int32
	Data          []byte
	Reader        *bufio.Reader
	r             io.ReadSeeker
	BytesReader   *bytes.Reader
}

/**
 * gets chunk at current offset
 */
// Deprecated: move away for hacked together BinaryChunk
func FromBinary(r io.ReadSeeker, typeExpected string, endian binary.ByteOrder) (*BinaryChunk, error) {
	return FromBinaryAt(r, typeExpected, endian, -1)
}

/**
 * Gets chunk at a specific offset
 */
// Deprecated: move away for hacked together BinaryChunk
func FromBinaryAt(r io.ReadSeeker, typeExpected string, endian binary.ByteOrder, offset int64) (*BinaryChunk, error) {
	return fromBinaryAt(r, typeExpected, endian, offset, false, false)
}

/**
 * Gets chunk at a specific offset
 * It works with aftershock binaries.
 */
// Deprecated: move away for hacked together BinaryChunk
func FromBinaryAtHeadless(r io.ReadSeeker, typeExpected string, endian binary.ByteOrder, offset int64, skipLength bool) (*BinaryChunk, error) {
	return fromBinaryAt(r, typeExpected, endian, offset, true, skipLength)
}

/**
 * Gets chunk at a specific offset
 * It works with aftershock binaries.
 */
// Deprecated: move away for hacked together BinaryChunk
func FromBinaryAtPartHeadless(r *bytes.Reader, typeExpected string, endian binary.ByteOrder) (*BinaryChunk, error) {
	chunk := &BinaryChunk{}
	chunk.BytesReader = r
	chunk.Reader = bufio.NewReader(chunk.BytesReader)
	chunk.Type = typeExpected

	return chunk, nil
}

/**
 * Gets chunk at a specific offset
 * It works with aftershock binaries.
 */
// Deprecated: move away for hacked together BinaryChunk
func fromBinaryAt(r io.ReadSeeker, typeExpected string, endian binary.ByteOrder, offset int64, headless bool, skipLength bool) (*BinaryChunk, error) {
	chunk := &BinaryChunk{}
	chunk.r = (r)

	if offset != -1 {
		_, err := chunk.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("error seeking to offset %d: %s", offset, err)
		}

	}

	if !headless {
		// Read chunk type, and convert to string.
		var chunkTypeBytes [4]byte
		if err := binary.Read(r, endian, &chunkTypeBytes); err != nil {
			return nil, fmt.Errorf("error reading chunk type (expecting %s) at %d: %s", typeExpected, offset, err)
		}

		if endian == binary.BigEndian {
			chunk.Type = string(chunkTypeBytes[:])
		} else {
			chunk.Type = string(utils.ReverseBytes(chunkTypeBytes[:]))
		}

		// Sanity check
		if chunk.Type != typeExpected {
			return nil, fmt.Errorf("expected %s chunk, got %s at offset %d", typeExpected, chunk.Type, offset)
		}
	} else {
		chunk.Type = typeExpected
	}

	if !skipLength {
		if chunk.IsVarLength() || headless {
			if chunk.IsCompressed() || headless {
				val, err := chunk.ReadVarInt()
				if err != nil {
					return nil, fmt.Errorf("error reading headless varint chunk length: %s", err)
				}
				chunk.Length = int32(val)
				// chunk.HexDump(true)
			} else {
				// If chunk is not compressed, proceed as before
				val, err := chunk.ReadVarInt()
				if err != nil {
					return nil, fmt.Errorf("error reading chunk varint length: %s", err)
				}
				chunk.Length = int32(val)
			}

		} else {
			if err := binary.Read(r, endian, &chunk.Length); err != nil {
				return nil, fmt.Errorf("error reading chunk int32 length: %s", err)
			}

		}
	}

	// Set start position
	chunk.StartPosition, _ = chunk.r.Seek(0, io.SeekCurrent)

	return chunk, nil

}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) IsCompressed() bool {
	return chunk.Type == "Fcdr" || chunk.Type == "ABMP"
}

// Deprecated: move away for hacked together BinaryChunk
func (c *BinaryChunk) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) HexDump(print bool) (string, error) {
	buf := make([]byte, chunk.Length)

	var currentPos int64
	if chunk.Data == nil {
		currentPos, _ = chunk.r.Seek(0, io.SeekCurrent)
		n, err := chunk.r.Read(buf)
		if err != nil {
			return "", err
		}
		buf = buf[:n]
		chunk.Data = buf

	} else {
		buf = chunk.Data
	}

	hexdump := hex.Dump(buf)

	if print {
		utils.DebugMsg(chunk.Type, "")
		utils.DebugMsg(chunk.Type, "Chunk type: %s", chunk.Type)
		utils.DebugMsg(chunk.Type, "Chunk length: %d", chunk.Length)
		utils.DebugMsg(chunk.Type, "hexdump: "+hexdump)
		utils.DebugMsg(chunk.Type, "")
	}

	if chunk.Data == nil {
		chunk.r.Seek(currentPos, io.SeekStart)
	}

	return hexdump, nil
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) HasMore() bool {
	currentPos, _ := chunk.r.Seek(0, io.SeekCurrent)

	return currentPos <= chunk.StartPosition+int64(chunk.Length)+8
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) Position() int64 {
	var currentPos int64

	if chunk.Reader != nil {
		if chunk.BytesReader != nil {
			currentPos = chunk.BytesReader.Size() - int64(chunk.BytesReader.Len())
		} else {
			currentPos = chunk.BytesReader.Size() - int64(chunk.Length)
		}
	}
	if chunk.r != nil {
		currentPos, _ = chunk.r.Seek(0, io.SeekCurrent)
	}

	if chunk.BytesReader != nil {
		currentPos, _ = chunk.BytesReader.Seek(0, io.SeekCurrent)
	}

	return currentPos
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadBytes(number int) ([]byte, error) {
	bytes := make([]byte, number)

	_, err := chunk.GetReader().Read(bytes)

	if err != nil {
		return bytes, err
	}

	return bytes, nil

}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) Read(p []byte) (n int, err error) {
	return chunk.GetReader().Read(p)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadByte() (byte, error) {
	var b []byte = make([]byte, 1)

	_, err := chunk.GetReader().Read(b)
	if err != nil {
		return 0, err
	}

	return b[0], nil

}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := chunk.GetReader().(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	} else {
		return -1, fmt.Errorf("reader does not support seeking")
	}
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) GetPosition() (int64, error) {
	return chunk.Seek(0, io.SeekCurrent)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadAllBytes() ([]byte, error) {
	if chunk.Data != nil {
		return chunk.Data, nil
	}

	currentPos, _ := chunk.GetPosition()

	chunk.Seek(0, io.SeekStart)

	bytes, err := chunk.ReadBytes(int(chunk.Length))
	if err != nil {
		return bytes, err
	}

	chunk.Seek(currentPos, io.SeekStart)

	return bytes, err

}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) GetReader() io.Reader {
	if chunk.BytesReader != nil {
		return chunk.BytesReader
	}
	if chunk.Reader == nil {
		return chunk.r
	} else {
		return chunk.Reader
	}
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadVarInt() (uint32, error) {
	var reader = chunk.GetReader()

	var val uint32
	var b byte
	var err error

	for {
		b = 0
		err = binary.Read(reader, binary.LittleEndian, &b)
		if err != nil {
			return 0, err
		}

		val = (val << 7) | uint32(b&0x7f)

		if b>>7 == 0 {
			break
		}
	}

	return (val), nil
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadInt32(endian binary.ByteOrder) (int32, error) {
	return utils.ReadInt32(chunk.GetReader(), endian)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadInt16(endian binary.ByteOrder) (int16, error) {
	return utils.ReadInt16(chunk.GetReader(), endian)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadInt8() (int8, error) {
	return utils.ReadInt8(chunk.GetReader())
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadUInt8() (uint8, error) {
	return utils.ReadUInt8(chunk.GetReader())
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadUInt32(endian binary.ByteOrder) (uint32, error) {
	return utils.ReadUInt32(chunk.GetReader(), endian)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadString(length int, endian binary.ByteOrder) (string, error) {
	return utils.ReadString(chunk.GetReader(), length, endian == binary.LittleEndian)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadNullTerminatedString() (string, error) {
	bytes := make([]byte, 0)
	for {
		b, err := chunk.ReadByte()
		if err != nil {
			return "", err
		}

		// Stop if we find the null character
		if b == 0 || b == 00 {
			break
		}

		bytes = append(bytes, b)
	}
	return string(bytes), nil
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadUInt16(endian binary.ByteOrder) (uint16, error) {
	return utils.ReadUInt16(chunk.GetReader(), endian)
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) ReadUByte() (int, error) {
	var reader = chunk.GetReader()

	b := make([]byte, 1)

	_, err := reader.Read(b)
	if err != nil {
		return 0, err
	}

	return int(b[0]), nil
}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) IsVarLength() bool {
	if chunk.Type == "Fver" || chunk.Type == "Fcdr" || chunk.Type == "ABMP" || chunk.Type == "FGEI" {
		return true
	}

	return false

}

// Deprecated: move away for hacked together BinaryChunk
func (chunk *BinaryChunk) DecompressedDump(projectName string, targetDir string, targetFileName string, pkg string) {
	targetDir = utils.CleanString(targetDir)
	targetFileName = utils.CleanString(targetFileName)
	if (targetDir) == "" {
		targetDir = "empty_name"
	}

	if chunk == nil {
		fmt.Printf("Chunk is nil\n")
		return
	}
	utils.DebugMsg("shockwave", "Dumping decompressed %s chunk to %s as %s", chunk.Type, targetDir, targetFileName)

	var outputFolder string
	if pkg == "" {
		outputFolder = filepath.Join(consts.PathDump, projectName, "chunks_abmp", string(targetDir))
	} else {
		outputFolder = filepath.Join(consts.PathDump, pkg, "file", projectName, "chunks_abmp", string(targetDir))
	}

	os.MkdirAll(outputFolder, os.ModePerm) // Ensure the output directory exists
	outputFile := filepath.Join(outputFolder, ""+string(targetFileName))

	chunkFile, err := os.Create(outputFile + ".bin")
	if err != nil {
		fmt.Printf("Error creating file: %s", outputFile+".bin")
		return
	}
	defer chunkFile.Close()

	_, err = chunkFile.Write(chunk.Data)
	if err != nil {
		panic(err)
	}

	err = chunkFile.Sync()
	if err != nil {
		panic(err)
	}

}
