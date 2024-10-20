package tdf

import (
	"encoding/hex"
	"testing"

	"github.com/ddkwork/encoding/struct2table"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/stretchr/testify/assert"
)

func TestGenBaseType(t *testing.T) {
	mylog.FormatAllFiles()
	g := stream.NewGeneratedFile()
	m := stream.NewOrderedMap("", "") // todo 更精确的tips
	m.Set("Integer", "bool,int8,int16,int32,int64,uint8,uint16,uint32,uint64")
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

func TestReflect(t *testing.T) {
	t.Skip()
	type (
		LID struct {
			LNM  string
			TYPE int64
		}
		ListElem struct {
			BOID int64 // Variable todo
			FLGS int64
			LID
		}
		TdfMsg struct {
			ALST []ListElem
			LMS  int64
			PNAM string
			PRID int64
			PRMS int64
			MXRC int64
			OFRC int64
		}
	)
	msg := TdfMsg{
		ALST: []ListElem{
			{
				BOID: 1, //? todo
				FLGS: 1,
				LID: LID{
					LNM:  "OSDKPlatformFriendList",
					TYPE: 1,
				},
			},
		},
		LMS:  0,
		PNAM: "",
		PRID: 0,
		PRMS: 0,
		MXRC: 4294967295,
		OFRC: 0,
	}
	root := struct2table.Marshal(msg)
	println(root.String())
}

func TestMarshal(t *testing.T) {
}

func TestNativeType_IsValid(t *testing.T) {
}

func Test_marshalList(t *testing.T) {
}

func Test_marshalMap(t *testing.T) {
}

func Test_marshalSingular(t *testing.T) {
	b := marshalSingular("BLOB", mylog.Check2(hex.DecodeString("deadbeef")))
	assert.Equal(t, []byte{0x8a, 0xcb, 0xe2, 0x2, 0x4, 0xde, 0xad, 0xbe, 0xef}, b.Bytes())

	expectedBytes := []byte{
		0xda, 0x1b, 0x35, 0x01, 0x0e, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00,
	}
	b = marshalSingular("VALU", "Hello, World!")
	assert.Equal(t, expectedBytes, b.Bytes())

	b = marshalSingular("VALU", 1337)
	assert.Equal(t, []byte{0xda, 0x1b, 0x35, 0x00, 0xb9, 0x14}, b.Bytes())
}

func Test_unmarshalSingular(t *testing.T) {
	tag, wireType, data := unmarshalSingular([]byte{0x8a, 0xcb, 0xe2, 0x2, 0x4, 0xde, 0xad, 0xbe, 0xef})
	assert.Equal(t, "BLOB", tag)
	assert.Equal(t, BinaryType, wireType) // 不使用idt机制，应该按类型取值保平安
	assert.Equal(t, mylog.Check2(hex.DecodeString("deadbeef")), data)

	tag, wireType, data = unmarshalSingular([]byte{
		0xda, 0x1b, 0x35, 0x01, 0x0e, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
		0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00,
	})
	assert.Equal(t, "VALU", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Hello, World!", data)

	tag, wireType, data = unmarshalSingular([]byte{0xda, 0x1b, 0x35, 0x00, 0xb9, 0x14})
	assert.Equal(t, "VALU", tag)
	assert.Equal(t, IntegerType, wireType)
	assert.Equal(t, uint32(1337), data)
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
	//if !bytes.Equal(buffer, expectedBytes) {
	//	t.Errorf("Expected %+v, but got %+v", expectedBytes, buffer)
	//}
}

func Test_marshalStruct(t *testing.T) {
	t.Skip()
	//{
	//        ALST<List>: type:3,size:1
	//        {
	//                BOID<Variable3>:v1:0,v2:0,v3:0
	//                FLGS<Integer>:1
	//                LID@<Struct>:{
	//                        LNM@<String>:"OSDKPlatformFriendList"
	//                        TYPE<Integer>:1
	//                }
	//        }
	//        LMS@<Integer>:0
	//        PNAM<String>:""
	//        PRID<Integer>:0
	//        PRMS<Integer>:0
	//        MXRC<Integer>:4294967295
	//        OFRC<Integer>:0
	//}
	stream.WriteTruncate("tdf.bin", packedBuf)
	type (
		LID struct {
			LNM  string `json:"lnm"`
			TYPE int64  `json:"type"`
		}
		ListElem struct {
			BOID int64 // Variable todo
			FLGS int64
			LID
		}
		TdfMsg struct {
			ALST []ListElem
			LMS  int64
			PNAM string
			PRID int64
			PRMS int64
			MXRC int64
			OFRC int64
		}
	)
	msg := TdfMsg{
		ALST: []ListElem{
			{
				BOID: 1, //? todo
				FLGS: 1,
				LID: LID{
					LNM:  "OSDKPlatformFriendList",
					TYPE: 1,
				},
			},
		},
		LMS:  0,
		PNAM: "",
		PRID: 0,
		PRMS: 0,
		MXRC: 4294967295,
		OFRC: 0,
	}

	marshal := Marshal(msg) // todo
	// todo merge header 16 bytes
	mylog.HexDump("marshal", marshal.Bytes())

	// assert.Equal(packedBuf, encoder.w)
	assert.Equal(t, packedBuf[16:], marshal)
}

//func Test_encode(t *testing.T) {
//	encoder := Heat2Encoder{
//		TdfEncoder: TdfEncoder{
//			w:              bytes.NewBuffer(nil),
//			mOwnByteBuffer: false,
//			mEncodeResult:  false,
//		},
//		mEncodeHeader: true,
//		mErrorCount:   0,
//		mBuf:          make([]byte, 10), // todo
//	}
//	// todo 构造结构体,生成tree和json，解码
//	encoder.encodeHeader(MakeTag("ALST"), TDF_TYPE_LIST)
//	encoder.encodeType(TDF_TYPE_STRUCT)
//	encoder.encodeVarsizeInteger(1)
//	encoder.encodeHeader(MakeTag("BOID"), TDF_TYPE_BLAZE_OBJECT_ID)
//	encoder.encodeVarsizeInteger(0)
//	encoder.encodeVarsizeInteger(0)
//	encoder.encodeVarsizeInteger(0)
//	encoder.encodeHeaderAndVarsizeInteger(MakeTag("FLGS"), 1)
//
//	encoder.encodeHeader(MakeTag("LID"), TDF_TYPE_STRUCT)
//	encoder.encodeHeader(MakeTag("LNM"), TDF_TYPE_STRING)
//	encoder.TDF_TYPE_STRING(MakeTag("LNM"), "OSDKPlatformFriendList")
//	encoder.encodeHeaderAndVarsizeInteger(MakeTag("TYPE"), 1)
//	encoder.encodeVarsizeInteger(0) // TdfStruct end
//
//	encoder.TDF_TYPE_STRING(MakeTag("PNAM"), "")
//
//	//		0xC2, 0xE8, 0x6D, //PNAM<String>:""
//	//		0x01,       //
//	//		0x01, 0x00, //""
//	encoder.encodeHeaderAndVarsizeInteger(MakeTag("PRID"), 0)
//	encoder.encodeHeaderAndVarsizeInteger(MakeTag("PRMS"), 0)
//	//		0x00, //end  list end ?
//	encoder.encodeHeaderAndVarsizeInteger(MakeTag("MXRC"), 4294967295) // todo test
//	encoder.encodeHeaderAndVarsizeInteger(MakeTag("OFRC"), 0)
//
//	// todo merge header 16 bytes
//
//	mylog.HexDump("masher", encoder.w)
//
//	assert.Equal(t, packedBuf, encoder.w)
//}

var packedBuf = []byte{
	0x00, //== ID_TERM 0 Reset ?
	0x00, //
	0x00, //

	0x5C,             // bufType if bufType >= int(tdf.TDF_TYPE_MAX) error
	0x00, 0x00, 0x00, // bufTag Uint32 << 8
	// bufType  bufTag 和传入的比较，不等或者大于退出
	// SkipElement

	0x19,                                           //
	0x00, 0x06, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, //

	//////////////////////////////////////////////
	//算法，头部的16字节，前8字节用于校验，
	//后8字节是一个大数，用于根据校验通过的类型跳过？
	//valueIsNegative := (buf[0] & VARSIZE_NEGATIVE) == VARSIZE_NEGATIVE
	//检查是否是负数
	//
	//hasMore := (buf[0] & VARSIZE_MORE) == VARSIZE_MORE
	//有没有更多的字节？
	//总之就是解码大整数  DecodeVarsizeInteger
	//
	//综合来看，钱8字节校验buf信息，后续的都是轮训解码大数和对应的字符串类型
	//
	//////////////////////////////////////////////
	//0x00, 0x00, 0x00, 0x5C, 0x00, 0x00, 0x00, 0x19,
	//0x00, 0x06, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00,

	0x86, 0xCC, 0xF4, // ALST<List>: type:3,size:1
	0x04, // List
	0x03, // yes
	0x01, // size yes

	0x8A, 0xFA, 0x64, // BOID<Variable3>:v1:0,v2:0,v3:0
	0x09, //
	0x00, //
	0x00, //
	0x00, //

	0x9A, 0xC9, 0xF3, // FLGS<Integer>:1
	0x00, //
	0x01, //

	0xB2, 0x99, 0x00, // LID@<Struct>:{
	0x03, //

	0xB2, 0xEB, 0x40, // LNM@<String>:"OSDKPlatformFriendList"
	0x01, //

	// OSDKPlatformFriendList ?
	0x17,                                                       //
	0x4F, 0x53, 0x44, 0x4B, 0x50, 0x6C, 0x61, 0x74, 0x66, 0x6F, //
	0x72, 0x6D, 0x46, 0x72, 0x69, 0x65, 0x6E, 0x64, 0x4C, 0x69, 0x73, //
	0x74, //
	0x00, //

	0xD3, 0x9C, 0x25, // TYPE<Integer>:1
	0x00, // type id
	0x01, // value

	0x00, //? TdfStruct end yes

	0xB2, 0xDC, 0xC0, // LMS@<Integer>:0
	0x00, //
	0x00, //

	0xC2, 0xE8, 0x6D, // PNAM<String>:""
	0x01,       //
	0x01, 0x00, //""

	0xC3, 0x2A, 0x64, // PRID<Integer>:0
	0x00, //
	0x00, //

	0xC3, 0x2B, 0x73, // PRMS<Integer>:0
	0x00, //
	0x00, //

	0x00, // end

	0xB7, 0x8C, 0xA3, // MXRC<Integer>:4294967295
	0x00,                   //
	0xBF,                   //? 感觉有问题，没开源细节不清楚
	0xFF, 0xFF, 0xFF, 0x1F, //

	0xBE, 0x6C, 0xA3, // OFRC<Integer>:0
	0x00, //
	0x00, //
}
