package tdf

const (
	ELEMENT_TYPE_MAX        = 0x1f
	ID_TERM                 = 0x00
	VARSIZE_MORE            = 0x80
	VARSIZE_NEGATIVE        = 0x40
	VARSIZE_VALUE_MASK      = 0x3f
	VARSIZE_MAX_ENCODE_SIZE = 10
	FLOAT_SIZE              = 4
	HEADER_TYPE_OFFSET      = 3
	HEADER_SIZE             = 4
	MAX_TAG_LENGTH          = 4
	TAG_CHAR_MIN            = 32
	TAG_CHAR_MAX            = 95
	TAG_UNSPECIFIED         = 0
)

func PromoteUnsignedByte(byteVal byte) int16 {
	return int16(byteVal) & 0xff
}

func PromoteUnsignedInteger(intVal int) int64 {
	return int64(intVal) & 0xffffffff
}

func PromoteUnsignedShort(shortVal int16) int {
	return int(shortVal) & 0xffff
}
