package tdf

import (
	"encoding/hex"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenBaseType(t *testing.T) {
	g := stream.NewGeneratedFile()
	m := stream.NewOrderedMap("", "") //todo 更精确的tips,似乎缺少bool的绑定
	m.Set("Integer", "int8,int16,int32,int64,uint8,uint16,uint32,uint64")
	m.Set("String", "string")
	m.Set("Binary", "[]byte")
	m.Set("Struct", "TDFStruct")
	m.Set("List", "[]any")
	m.Set("Map", "map[any]any")
	m.Set("Union", "Union")
	m.Set("Variable", "Variable")
	m.Set("BlazeObjectType", "BlazeObjectType")
	m.Set("BlazeObjectID", "BlazeObjectID")
	m.Set("Float", "float32,float64")
	m.Set("TimeValue", "time.Time")
	m.Set("Max", "must less it")
	g.Types("base", m)
}

func TestNode_Marshal(t *testing.T) {

}

func TestInteger_decode(t *testing.T) {
	bytesData := []byte{0xda, 0x1b, 0x35, 0x00, 0xb9, 0x14}
	n := NewNode(bytesData)
	assert.Equal(t, "VALU", n.tag)
	assert.Equal(t, IntegerType, n.wireType)
	assert.Equal(t, uint32(1337), n.DecodeInteger())
}

func TestInteger_encode(t *testing.T) {
	expectedBytes := []byte{0xda, 0x1b, 0x35, 0x00, 0xb9, 0x14}
	assert.Equal(t, expectedBytes, NewNode(nil).EncodeInteger("VALU", 1337))
}

func TestList_encode(t *testing.T) {

}

func TestString_decode(t *testing.T) {
	bytesData := []byte{
		0xda, 0x1b, 0x35, 0x01, 0x0e, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00,
	}
	n := NewNode(bytesData)
	assert.Equal(t, "VALU", n.tag)
	assert.Equal(t, StringType, n.wireType)
	assert.Equal(t, "Hello, World!", n.DecodeString())
}

func TestString_encode(t1 *testing.T) {
	expectedBytes := []byte{
		0xda, 0x1b, 0x35, 0x01, 0x0e, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00,
	}
	assert.Equal(t1, expectedBytes, NewNode(nil).EncodeString("VALU", "Hello, World!"))
}

func TestReadBlob(t *testing.T) {
	bytesData := []byte{0x8a, 0xcb, 0xe2, 0x2, 0x4, 0xde, 0xad, 0xbe, 0xef}
	n := NewNode(bytesData)
	assert.Equal(t, "BLOB", n.tag)
	assert.Equal(t, BinaryType, n.wireType)
	assert.Equal(t, mylog.Check2(hex.DecodeString("deadbeef")), n.DecodeBlob())
}

func TestWriteBlob(t *testing.T) {
	expectedBytes := []byte{0x8a, 0xcb, 0xe2, 0x2, 0x4, 0xde, 0xad, 0xbe, 0xef}
	assert.Equal(t, expectedBytes, NewNode(nil).EncodeBlob("BLOB", mylog.Check2(hex.DecodeString("deadbeef"))))
}

func TestUnionRead(t *testing.T) {

}

func TestReadStruct(t *testing.T) {
	//b := []byte{
	//	0xcf, 0x4c, 0xa3, 0x03, 0xa6, 0xed, 0x00, 0x00, 0x2a, 0xcf,
	//	0x4c, 0x80, 0x01, 0x0e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	//	0x20, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x00, 0x00,
	//}

	// 假设 Node.UmMarshal 是一个可以从字节切片读取结构的函数
	//tdf := Node.ReadTDF(b)
	//
	//expected := TDFStruct{
	//	Label: "STRC",
	//	Content: []any{
	//		TDFInteger{"INT ", 42},
	//		TDFString{"STR ", "Nested struct"},
	//	},
	//}
	//
	//if !tdf.Equals(expected) { // 需要实现 Equals 方法来比较两个 TDFStruct
	//	t.Errorf("Expected %+v, but got %+v", expected, tdf)
	//}
}

func TestWriteStruct(t *testing.T) {
	//payload := []TDFStruct{
	//	{
	//		Label: "STRC",
	//		Content: []any{
	//			TDFInteger{"INT ", 42},
	//			TDFString{"STR ", "Nested struct"},
	//		},
	//	},
	//}
	//
	//var buffer bytes.Buffer
	//for _, x := range payload {
	//	buffer.Write(x.Write())
	//}
	//
	//expectedBytes := []byte{
	//	0xcf, 0x4c, 0xa3, 0x03, 0xa6, 0xed, 0x00, 0x00, 0x2a, 0xcf,
	//	0x4c, 0x80, 0x01, 0x0e, 0x4e, 0x65, 0x73, 0x74, 0x65, 0x64,
	//	0x20, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x00, 0x00,
	//}
	//
	//if !bytes.Equal(buffer.Bytes(), expectedBytes) {
	//	t.Errorf("Expected %+v, but got %+v", expectedBytes, buffer.Bytes())
	//}
}
