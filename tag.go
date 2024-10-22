package tdf

import (
	"encoding/binary"
	"reflect"
	"slices"
	"strings"
)

func MakeTag(tag string) uint32 {
	if tag == "" {
		return 0
	}

	var tmptag [MAX_TAG_LENGTH]byte
	tmptag[0] = '\x00'

	// ensure tag string is valid
	size := len(tag)
	if size > MAX_TAG_LENGTH {
		return 0
	}
	if size < 1 {
		return 0
	}

	for idx := 0; idx < size; idx++ {
		tmptag[idx] = strings.ToUpper(string(tag[idx]))[0]

		// Check valid range for tag characters
		if tmptag[idx] < byte(TAG_CHAR_MIN) || tmptag[idx] > byte(TAG_CHAR_MAX) {
			return 0
		}
	}

	if tmptag[0] < 'A' || tmptag[0] > 'Z' {
		return 0
	}

	t := uint32((tmptag[0]-byte(TAG_CHAR_MIN))&0x3f) << 26
	if size > 1 {
		t |= uint32((tmptag[1]-byte(TAG_CHAR_MIN))&0x3f) << 20
		if size > 2 {
			t |= uint32((tmptag[2]-byte(TAG_CHAR_MIN))&0x3f) << 14
			if size > 3 {
				t |= uint32(tmptag[3]-byte(TAG_CHAR_MIN)) & 0x3f << 8
			}
		}
	}
	return t
}

func DecodeTag11(tag uint32) string {
	buf := make([]byte, MAX_TAG_LENGTH+1)
	tag >>= 8
	for i := 3; i >= 0; i-- {
		sixbits := tag & 0x3f
		if sixbits != 0 {
			buf[i] = byte(sixbits + 32)
		} else {
			buf[i] = '\x00'
		}
		tag >>= 6
	}
	buf[4] = '\x00'
	buf = buf[:MAX_TAG_LENGTH]
	convertToLowercase := false
	if convertToLowercase {
		for i := 0; i < 4; i++ {
			buf[i] = strings.ToLower(string(buf[i]))[0]
		}
	}
	return string(buf)
}

func DecodeTag_(tag int64, convertToLowercase bool) string {
	buf := make([]byte, 5)
	// size := 4
	tag >>= 8
	for i := 3; i >= 0; i-- {
		sixbits := tag & 0x3f
		if sixbits != 0 {
			buf[i] = byte(sixbits + 32)
		} else {
			buf[i] = '\x00'
			// size = i
		}
		tag >>= 6
	}
	if convertToLowercase {
		buf = []byte(strings.ToLower(string(buf)))
		//for i := 0; i < size; i++ {
		//	buf[i] = byte(strings.ToLower(int(buf[i])))
		//}
	}
	actualLength := 0
	for actualLength = 0; actualLength < len(buf); actualLength++ {
		if buf[actualLength] == 0 {
			break
		}
	}
	return string(buf[:actualLength])
}

type TagInfo struct {
	MemberNativeType reflect.Type
	Tag              int64
	MemberIndex      int64
	Member           string
	Description      string
	Default          string
}

func NewTagInfo(memberNativeType reflect.Type, tag int64, member string, memberIndex int64, description, def string) *TagInfo {
	return &TagInfo{
		MemberNativeType: memberNativeType,
		Tag:              tag,
		MemberIndex:      memberIndex,
		Member:           member,
		Description:      description,
		Default:          def,
	}
}

func (t *TagInfo) GetMemberNativeType() reflect.Type {
	return t.MemberNativeType
}

func (t *TagInfo) GetTag() int64 {
	return t.Tag
}

func (t *TagInfo) GetMemberIndex() int64 {
	return t.MemberIndex
}

func (t *TagInfo) GetMember() string {
	return t.Member
}

func (t *TagInfo) GetDescription() string {
	return t.Description
}

func (t *TagInfo) GetDefault() string {
	return t.Default
}

/////////////////////////////////////

type (
	InterfaceAscii interface {
		AsciiPack(ascii string) []byte
		AsciiUnPack(packedBuf []byte) string
	}
	ObjAscii struct{}
	ctx      struct {
		ascii string
		value uint32
		buf   []byte
	}
)

func (ObjAscii) AsciiPack(ascii string) []byte {
	return new(ctx).AsciiPack(ascii)
}

func (ObjAscii) AsciiUnPack(packedBuf []byte) string {
	return new(ctx).AsciiUnPack(packedBuf)
}

const maxAsciiLen = 4

func (p *ctx) AsciiPack(ascii string) []byte {
	*p = ctx{
		ascii: ascii,
		value: 0,
		buf:   make([]byte, maxAsciiLen),
	}
	for i := 0; i < len(p.ascii); i++ {
		p.value |= (32 | uint32(p.ascii[i])&31) << ((3 - i) * 6)
	}
	binary.BigEndian.PutUint32(p.buf, p.value)
	p.buf = p.buf[1:]
	return p.buf
}

func (p *ctx) AsciiUnPack(packedBuf []byte) string {
	*p = ctx{
		ascii: "",
		value: 0,
		buf:   packedBuf,
	}
	buf0 := make([]byte, 1)
	p.buf = slices.Concat(buf0, p.buf)
	p.value = binary.BigEndian.Uint32(p.buf)
	tmp := make([]byte, maxAsciiLen)
	for i := 0; i < maxAsciiLen; i++ {
		v := (p.value >> ((3 - i) * 6)) & 63
		tmp[i] = byte(0x40 | (v & 0x1F))
	}
	p.ascii = string(tmp)
	return p.ascii
}
