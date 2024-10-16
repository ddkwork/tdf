package tdf

import (
	"bytes"
	"encoding/binary"
	"github.com/ddkwork/golibrary/mylog"
	"slices"
	"strings"
)

func DecodeTag(tag []byte) string {
	if len(tag) != 3 {
		mylog.Check("tagBuf length must be 3")
	}
	tag = slices.Concat([]byte{0}, tag)
	encodedTag := binary.BigEndian.Uint32(tag)
	decodedTag := make([]byte, 4)
	for i := 3; i >= 0; i-- {
		sixbits := encodedTag & 0x3F
		if sixbits != 0 {
			decodedTag[i] = byte(sixbits + 32)
		} else {
			decodedTag[i] = ' '
		}
		encodedTag >>= 6
	}
	return string(decodedTag)
}

func EncodeTag(tag string) []byte {
	if len(tag) != 4 {
		mylog.Check("tag must be 4 characters long")
	}

	tag = strings.ToUpper(tag)
	var encodedTag uint32 = 0

	for i := 0; i < len(tag); i++ {
		char := tag[i]
		if char == ' ' {
			continue
		}
		encodedTag |= (uint32(char) - 32) << (6 * (3 - i))
	}

	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, encodedTag)
	return result[1:]
}

func decompressInteger(b *bytes.Buffer) uint32 {
	var result uint64 = 0
	var currentShift uint64 = 6
	buffer := make([]byte, 1)
	mylog.Check2(b.Read(buffer))
	result += uint64(buffer[0]) & 0x3F
	for buffer[0]&0x80 != 0 {
		if _, err := b.Read(buffer); err != nil {
			mylog.CheckEof(err) //todo test  Check2 是否可以工作
		}
		result |= (uint64(buffer[0]) & 0x7F) << currentShift
		currentShift += 7
	}
	return uint32(result)
}

func compressInteger(value uint32) []byte {
	var result []byte
	if value < 0x40 {
		result = append(result, byte(value))
	} else {
		currentByte := (value & 0x3F) | 0x80
		result = append(result, byte(currentByte))
		currentShift := value >> 6
		for currentShift >= 0x80 {
			currentByte = (currentShift & 0x7F) | 0x80
			currentShift >>= 7
			result = append(result, byte(currentByte))
		}
		result = append(result, byte(currentShift))
	}

	return result
}
