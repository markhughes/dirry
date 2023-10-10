package bitd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func unpackPackbits1(input []byte) ([]byte, error) {
	// For 1-bit images,  copy the input data to the output?
	return input, nil
}

func unpackPackbits2(input []byte) ([]byte, error) {
	if len(input)%4 != 0 {
		return nil, fmt.Errorf("input length must be a multiple of 4")
	}

	output := make([]byte, len(input))

	// Iterate over the input data in chunks of 4 bytes..
	for i := 0; i < len(input); i += 4 {
		tmp := input[i]

		output[i] = input[i+1]
		output[i+1] = input[i+2]
		output[i+2] = input[i+3]

		output[i+3] = 255 - tmp
	}

	return output, nil
}

func unpackPackbits8(input []byte) ([]byte, error) {
	reader := bytes.NewReader(input)

	var output bytes.Buffer

	for {
		var count int8
		err := binary.Read(reader, binary.BigEndian, &count)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read count byte: %s", err)
		}

		if count >= 0 {
			// If the count is non-negative, read that many bytes and write them to the output?
			data := make([]byte, count+1)
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return nil, fmt.Errorf("failed to read data bytes: %s", err)
			}
			output.Write(data)
		} else {
			// If the count is negative, read one byte and write it to the output that many times?
			var data byte
			err := binary.Read(reader, binary.BigEndian, &data)
			if err != nil {
				return nil, fmt.Errorf("failed to read data byte: %s", err)
			}
			for i := int8(0); i < 1-count; i++ {
				output.WriteByte(data)
			}
		}
	}

	return output.Bytes(), nil
}

func unpackPackbits24(input []byte, width int) ([]byte, error) {
	reader := bytes.NewReader(input)

	var output bytes.Buffer

	buffer := make([]byte, width*3)

	j := 0
	for {
		var count int8
		err := binary.Read(reader, binary.BigEndian, &count)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read count byte: %s", err)
		}

		if count >= 0 {
			// Againm, if the count is non-negative, read that many bytes and write them to the buffer
			data := make([]byte, count+1)
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return nil, fmt.Errorf("failed to read data bytes: %s", err)
			}
			copy(buffer[j:], data)
			j += len(data)
		} else {
			// if negative read one byte and write it to the buffer that many times
			var data byte
			err := binary.Read(reader, binary.BigEndian, &data)
			if err != nil {
				return nil, fmt.Errorf("failed to read data byte: %s", err)
			}
			for i := int8(0); i < 1-count; i++ {
				buffer[j] = data
				j++
			}
		}

		if j == len(buffer) {
			// success-o!
			for k := 0; k < width; k++ {
				output.WriteByte(buffer[k])
				output.WriteByte(buffer[k+width])
				output.WriteByte(buffer[k+width*2])
			}
			j = 0
		}
	}

	return output.Bytes(), nil
}

func unpackPackbits32(input []byte, width int) ([]byte, error) {
	reader := bytes.NewReader(input)

	var output bytes.Buffer

	buffer := make([]byte, width*4)

	j := 0
	for {
		var count int8
		err := binary.Read(reader, binary.BigEndian, &count)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read count byte: %s", err)
		}

		if count >= 0 {
			// non-negative we read that many bytes and write them to the buffer
			data := make([]byte, count+1)
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return nil, fmt.Errorf("failed to read data bytes: %s", err)
			}
			copy(buffer[j:], data)
			j += len(data)
		} else {
			// Inegative, read one byte, write it to the buffer {count times
			var data byte
			err := binary.Read(reader, binary.BigEndian, &data)
			if err != nil {
				return nil, fmt.Errorf("failed to read data byte: %s", err)
			}
			for i := int8(0); i < 1-count; i++ {
				buffer[j] = data
				j++
			}
		}

		if j == len(buffer) {
			// cccccoool
			for k := 0; k < width; k++ {
				output.WriteByte(buffer[k+width])
				output.WriteByte(buffer[k+width*2])
				output.WriteByte(buffer[k+width*3])
				alpha := 255 - buffer[k]
				output.WriteByte(alpha)
			}
			j = 0
		}
	}

	return output.Bytes(), nil
}
