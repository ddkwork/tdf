package tdf

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/golibrary/mylog"
)

func Test_Tag(t *testing.T) {
	tag := MakeTag("CFID")

	mylog.Hex("makeTag", tag)
	assert.Equal(t, uint32(0x8e6a6400), tag)
	// mylog.Info("decodeTag", DecodeTag(tag))
}

func Test2(t *testing.T) {
	p := new(ObjAscii)
	ascii := "CFID"
	packedBuf := p.AsciiPack(ascii)
	mylog.HexDump("Pack "+strconv.Quote(ascii), packedBuf)
	mylog.Success("UnPack", strconv.Quote(p.AsciiUnPack(packedBuf)))
}
