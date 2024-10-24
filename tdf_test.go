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

type Packet struct {
	Sid         int
	Time        string
	Action      string
	DataSummary string
	Note        string
	Size        int
	decoded     string
	hexdump     stream.HexDumpString
}

func TestName(t *testing.T) {
	type ListElem struct {
		Tag   string
		Value any
	}
	//s:=[]string{"a","b","c"}
	//tags := []string{"LNAM", "SNAM", "DESC"}//?如果不使用结构体，切片需要手动传入tag
	//lists := []ListElem{
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "Expert Scouting"},
	//	{Tag: "SNAM", Value: "Trait Expert Scout"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "Increased Player Weekly Goal XP"},
	//	{Tag: "SNAM", Value: "Trait Increased Experience"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "Become Predictable"},
	//	{Tag: "SNAM", Value: "Trait Predictability"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "Master Trade Negotiator"},
	//	{Tag: "SNAM", Value: "Trait Trade Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "QB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait QB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "RB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait RB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "WR Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait WR Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "TE Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait TE Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "OL Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait OL Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DL Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DL Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "LB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait LB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "K Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait K Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "P Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait P Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "S Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait S Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "CB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait CB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DE Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DE Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "NT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait NT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "FS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait FS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "SS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait SS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "K Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait K Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "P Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait P Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "S Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait S Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "CB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait CB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DE Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DE Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "NT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait NT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "FS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait FS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "SS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait SS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "K Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait K Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "P Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait P Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "S Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait S Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "CB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait CB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DE Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DE Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "NT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait NT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "FS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait FS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "SS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait SS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "K Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait K Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "P Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait P Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "S Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait S Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "CB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait CB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DE Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DE Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "NT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait NT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "FS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait FS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "SS Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait SS Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "K Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait K Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "P Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait P Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "S Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait S Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "CB Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait CB Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DE Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DE Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "DT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait DT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "NT Free Agency Influence"},
	//	{Tag: "SNAM", Value: "Trait NT Free Agent Influence"},
	//	{Tag: "DESC", Value: ""},
	//	{Tag: "LNAM", Value: "FS Free Agency Influence"},
	//}

	exp := stream.NewHexDump(tdf_test_data.hexdump)
	exp.ReadN(16)
	tag, baseType := decodeTagAndWireType(exp)
	assert.Equal(t, "CTET", tag)
	assert.Equal(t, ListType, baseType)

	elemType := BaseType(decompressInteger(exp))
	assert.Equal(t, StructType, elemType)

	length := decompressInteger(exp)
	assert.Equal(t, uint32(44), length) //2c 结构体大小

	//所以这里相当于rootRows的每一行，深度1下的所有字段，前面大部分是结构体字段，后面有3个整数和1个map
	//解码应该改成结构体切片
	//切片 tlv
	//结构体 tlv，因为根据java，因为父级是切片，结构体类型只需要编码一次
	//三个字符串类型字段的结构体，klv，多个元素，注意结构体类型因为父级不编码了，但是id-term需要编码
	//深度为1对应的id-term,至此切片编解码结束

	//然后来到这里{FORM:<9A FC AD >:TdfStruct:size=:TdfStruct:size=5:
	//{DICT:<92 98 F4 >:TdfMap:0-3 size=16:  结构体进了map ？
	//... struct

	//tag, baseType = decodeTagAndWireType(exp)
	//assert.Equal(t, "DESC", tag)
	//assert.Equal(t, StringType, baseType)
	//这里收尾
	//	{RIBC:<CA 98 A3 >:TdfInteger:0x11/17}
	//	{ROOT:<CA FB F4 >:TdfInteger:0x5280016/86507542}
	//	{SIBC:<CE 98 A3 >:TdfInteger:0xE/14}
	//	{TABL:<D2 18 AC >:TdfMap:0-3 size=15:
	//		{460:
	//			{:<>:TdfStruct:size=:TdfStruct:size=2:

	tag, wireType, data := unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Expert Scouting", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait Expert Scout", data)

	exp.ReadByte() //struct end ID_TERM

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Increased Player Weekly Goal XP", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait Increased Experience", data)

	exp.ReadByte()

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Become Predictable", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait Predictability", data)

	exp.ReadByte()

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Master Trade Negotiator", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait Trade Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "QB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait QB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "RB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait RB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "WR Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait WR Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "TE Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait TE Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "OL Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait OL Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "DL Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait DL Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "LB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait LB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "DB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait DB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "K Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait K Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "P Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait P Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "S Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait S Free Agent Influence", data)

	exp.ReadByte()
	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "DESC", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "LNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "CB Free Agency Influence", data)

	tag, wireType, data = unmarshalSingular(exp)
	assert.Equal(t, "SNAM", tag)
	assert.Equal(t, StringType, wireType)
	assert.Equal(t, "Trait CB Free Agent Influence", data)

}

// f3 index
var tdf_test_data = Packet{
	Sid:         0,
	Time:        "",
	Action:      "",
	DataSummary: "",
	Note:        "",
	Size:        0,
	decoded: `
===========================begin
{CTET:<8F 49 74 >:TdfList:Type:3: size=44:

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"Expert Scouting"}
	{SNAM:<CE E8 6D >:TdfString:"Trait Expert Scout"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"Increased Player Weekly Goal XP"}
	{SNAM:<CE E8 6D >:TdfString:"Trait Increased Experience"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"Become Predictable"}
	{SNAM:<CE E8 6D >:TdfString:"Trait Predictability"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"Master Trade Negotiator"}
	{SNAM:<CE E8 6D >:TdfString:"Trait Trade Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"QB Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait QB Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"RB Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait RB Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"WR Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait WR Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"TE Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait TE Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"OL Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait OL Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DL Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DL Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"LB Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait LB Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DB Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DB Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"K Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait K Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"P Free Agency Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait P Free Agent Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"QB Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait QB Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"RB Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait RB Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"WR Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait WR Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"TE Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait TE Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"OL Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait OL Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DL Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DL Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"LB Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait LB Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DB Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DB Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"K Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait K Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"P Training Boost"}
	{SNAM:<CE E8 6D >:TdfString:"Trait P Player Progression"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"QB Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait QB Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"RB Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait RB Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"WR Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait WR Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"TE Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait TE Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"OL Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait OL Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DL Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DL Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"LB Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait LB Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DB Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DB Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"K Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait K Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"P Re-Sign Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait P Contract Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"QB Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait QB Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"RB Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait RB Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"WR Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait WR Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"TE Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait TE Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"OL Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait OL Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DL Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DL Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"LB Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait LB Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"DB Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait DB Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"K Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait K Retirement Influence"}
	}

	{DESC:<92 5C E3 >:TdfString:nullptr}
	{LNAM:<B2 E8 6D >:TdfString:"P Retirement Influence"}
	{SNAM:<CE E8 6D >:TdfString:"Trait P Retirement Influence"}
	}
}
{FORM:<9A FC AD >:TdfStruct:size=:TdfStruct:size=5:
	{DICT:<92 98 F4 >:TdfMap:0-3 size=16:
		{59:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=1:
					{"IsUserControlled":0}
				}
				{NAME:<BA 1B 65 >:TdfString:"UserEntity"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{65:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=1:
					{"IsSubmittable":0}
				}
				{NAME:<BA 1B 65 >:TdfString:"ResponseForm"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{82:
			{:<>:TdfStruct:size=:TdfStruct:size=3:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{NAME:<BA 1B 65 >:TdfString:"Team"}
				{TYPE:<D3 9C 25 >:TdfInteger:0x1/1}
			}
		}
		{83:
			{:<>:TdfStruct:size=:TdfStruct:size=3:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{NAME:<BA 1B 65 >:TdfString:"Team[]"}
				{TYPE:<D3 9C 25 >:TdfInteger:0x1/1}
			}
		}
		{398:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x3B/59}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=307:
					{"AccelerationRating":1}
					{"Age":2}
					{"AgilityRating":3}
					{"AwarenessRating":4}
					{"BCVisionRating":6}
					{"Background":5}
					{"BlockSheddingRating":7}
					{"BreakSackRating":8}
					{"BreakTackleRating":9}
					{"CaptainsPatch":10}
					{"CareerStats":11}
					{"CarryingRating":12}
					{"CatchInTrafficRating":14}
					{"CatchingRating":13}
					{"College":15}
					{"ConfidenceRating":16}
					{"ContractBonus0":17}
					{"ContractBonus1":18}
					{"ContractBonus2":19}
					{"ContractBonus3":20}
					{"ContractBonus4":21}
					{"ContractBonus5":22}
					{"ContractBonus6":23}
					{"ContractBonus7":24}
					{"ContractExtraYearOption":25}
					{"ContractLength":26}
					{"ContractSalary0":27}
					{"ContractSalary1":28}
					{"ContractSalary2":29}
					{"ContractSalary3":30}
					{"ContractSalary4":31}
					{"ContractSalary5":32}
					{"ContractSalary6":33}
					{"ContractSalary7":34}
					{"ContractStatus":35}
					{"ContractYear":36}
					{"CurrentYearSeasonEndingInjuryWeek":37}
					{"DeepRouteRunningRating":38}
					{"ENDPLAY_ANIMVAL":40}
					{"ElusivenessRating":39}
					{"ExperiencePoints":41}
					{"FinesseMovesRating":42}
					{"FirstName":43}
					{"GameStats":44}
					{"Gender":45}
					{"Height":46}
					{"HitPowerRating":47}
					{"ImpactBlockingRating":48}
					{"InjuryRating":49}
					{"InjurySeverity":50}
					{"InjuryStatus":51}
					{"InjuryType":52}
					{"IsCreated":53}
					{"IsDemandRelease":54}
					{"IsInjuredReserve":55}
					{"IsLegend":56}
					{"IsUserControlled":0}
					{"JerseyNum":57}
					{"JukeMoveRating":58}
					{"JumpingRating":59}
					{"KickAccuracyRating":60}
					{"KickPowerRating":61}
					{"KickReturnRating":62}
					{"LastName":63}
					{"LastSeasonProBowlAppearances":64}
					{"LastYearSeasonEndingInjuryWeek":65}
					{"LatestInjuryStage":66}
					{"LatestInjuryWeek":67}
					{"LatestInjuryYear":68}
					{"LeadBlockRating":69}
					{"LegacyScore":70}
					{"ManCoverageRating":71}
					{"MaxInjuryDuration":72}
					{"MediumRouteRunningRating":73}
					{"MetaMorph_ArmsBarycentric":74}
					{"MetaMorph_ArmsBase":75}
					{"MetaMorph_BackPlateAmount":76}
					{"MetaMorph_CalfsBarycentric":77}
					{"MetaMorph_CalfsBase":78}
					{"MetaMorph_ChestBarycentric":79}
					{"MetaMorph_ChestBase":80}
					{"MetaMorph_FeetBarycentric":81}
					{"MetaMorph_FeetBase":82}
					{"MetaMorph_FlakJacketAmount":83}
					{"MetaMorph_GluteBarycentric":84}
					{"MetaMorph_GluteBase":85}
					{"MetaMorph_GutBarycentric":86}
					{"MetaMorph_GutBase":87}
					{"MetaMorph_ShoulderHeight":88}
					{"MetaMorph_ShoulderWidth":89}
					{"MetaMorph_ThighsBarycentric":90}
					{"MetaMorph_ThighsBase":91}
					{"MinInjuryDuration":92}
					{"OriginalAccelerationRating":93}
					{"OriginalAgilityRating":94}
					{"OriginalAwarenessRating":95}
					{"OriginalBCVisionRating":96}
					{"OriginalBlockSheddingRating":97}
					{"OriginalBreakSackRating":98}
					{"OriginalBreakTackleRating":99}
					{"OriginalCarryingRating":100}
					{"OriginalCatchInTrafficRating":102}
					{"OriginalCatchingRating":101}
					{"OriginalDeepRouteRunningRating":103}
					{"OriginalElusivenessRating":104}
					{"OriginalFinesseMovesRating":105}
					{"OriginalHitPowerRating":106}
					{"OriginalImpactBlockingRating":107}
					{"OriginalInjuryRating":108}
					{"OriginalJukeMoveRating":109}
					{"OriginalJumpingRating":110}
					{"OriginalKickAccuracyRating":111}
					{"OriginalKickPowerRating":112}
					{"OriginalKickReturnRating":113}
					{"OriginalLeadBlockRating":114}
					{"OriginalManCoverageRating":115}
					{"OriginalMediumRouteRunningRating":116}
					{"OriginalOverallRating":117}
					{"OriginalPassBlockFinesseRating":118}
					{"OriginalPassBlockPowerRating":119}
					{"OriginalPassBlockRating":120}
					{"OriginalPlayActionRating":121}
					{"OriginalPlayRecognitionRating":122}
					{"OriginalPowerMovesRating":123}
					{"OriginalPressRating":124}
					{"OriginalPursuitRating":125}
					{"OriginalReleaseRating":126}
					{"OriginalRunBlockFinesseRating":127}
					{"OriginalRunBlockPowerRating":128}
					{"OriginalRunBlockRating":129}
					{"OriginalShortRouteRunningRating":130}
					{"OriginalSpectacularCatchRating":131}
					{"OriginalSpeedRating":132}
					{"OriginalSpinMoveRating":133}
					{"OriginalStaminaRating":134}
					{"OriginalStiffArmRating":135}
					{"OriginalStrengthRating":136}
					{"OriginalTackleRating":137}
					{"OriginalThrowAccuracyDeepRating":138}
					{"OriginalThrowAccuracyMidRating":139}
					{"OriginalThrowAccuracyRating":140}
					{"OriginalThrowAccuracyShortRating":141}
					{"OriginalThrowOnTheRunRating":142}
					{"OriginalThrowPowerRating":143}
					{"OriginalThrowUnderPressureRating":144}
					{"OriginalToughnessRating":145}
					{"OriginalTruckingRating":146}
					{"OriginalZoneCoverageRating":147}
					{"OverallGrade0":148}
					{"OverallGrade1":149}
					{"OverallGrade2":150}
					{"OverallGrade3":151}
					{"OverallGrade4":152}
					{"OverallRating":153}
					{"PLYR_ASSETNAME":167}
					{"PLYR_BACKPLATE":168}
					{"PLYR_BIRTHDATE":169}
					{"PLYR_BREATHERITE":170}
					{"PLYR_CAPSALARY":171}
					{"PLYR_CAREERPHASE":172}
					{"PLYR_CELEBRATION":173}
					{"PLYR_COMMENT":174}
					{"PLYR_CONSECYEARSWITHTEAM":175}
					{"PLYR_DRAFTPICK":176}
					{"PLYR_DRAFTROUND":177}
					{"PLYR_DRAFTTEAM":178}
					{"PLYR_EGO":179}
					{"PLYR_EYEPAINT":180}
					{"PLYR_FACEMASK":181}
					{"PLYR_FLAGPROBOWL":182}
					{"PLYR_FLAKJACKET":183}
					{"PLYR_GENERICHEAD":184}
					{"PLYR_GRASSLEFTELBOW":185}
					{"PLYR_GRASSLEFTHAND":186}
					{"PLYR_GRASSLEFTWRIST":187}
					{"PLYR_GRASSRIGHTELBOW":188}
					{"PLYR_GRASSRIGHTHAND":189}
					{"PLYR_GRASSRIGHTWRIST":190}
					{"PLYR_HAIRCOLOR":191}
					{"PLYR_HANDEDNESS":192}
					{"PLYR_HANDWARMER":193}
					{"PLYR_HELMET":194}
					{"PLYR_HOME_STATE":195}
					{"PLYR_HOME_TOWN":196}
					{"PLYR_ICON":197}
					{"PLYR_ISCAPTAIN":198}
					{"PLYR_JERSEYSLEEVE":199}
					{"PLYR_JERSEYTYPE":200}
					{"PLYR_LASTHOLDOUTYEAR":201}
					{"PLYR_LEFTARMSLEEVE":202}
					{"PLYR_LEFTKNEE":203}
					{"PLYR_LEFTSHOE":204}
					{"PLYR_LEFTSPAT":205}
					{"PLYR_LEFTTHIGH":206}
					{"PLYR_MOUTHPIECE":207}
					{"PLYR_NECKPAD":208}
					{"PLYR_NECKTYPE":209}
					{"PLYR_PERFORMLEVEL":210}
					{"PLYR_PORTRAIT":211}
					{"PLYR_PREVTEAMID":212}
					{"PLYR_QBSTYLE":213}
					{"PLYR_RIGHTARMSLEEVE":214}
					{"PLYR_RIGHTKNEE":215}
					{"PLYR_RIGHTSHOE":216}
					{"PLYR_RIGHTSPAT":217}
					{"PLYR_RIGHTTHIGH":218}
					{"PLYR_SIDELINE_HEADGEAR":219}
					{"PLYR_SKIN":220}
					{"PLYR_SLEEVETEMPERATURE":221}
					{"PLYR_SOCK_HEIGHT":222}
					{"PLYR_STANCE":223}
					{"PLYR_STYLE":224}
					{"PLYR_TENDENCY":225}
					{"PLYR_TOWEL":226}
					{"PLYR_UNDERSHIRT":227}
					{"PLYR_VISOR":228}
					{"PassBlockFinesseRating":154}
					{"PassBlockPowerRating":155}
					{"PassBlockRating":156}
					{"PersonalityRating":157}
					{"PlayActionRating":158}
					{"PlayRecognitionRating":166}
					{"PlayerBottom":159}
					{"PlayerPercentage":160}
					{"PlayerTop":161}
					{"PlayerType":162}
					{"PlayoffConferenceWins":163}
					{"PlayoffDivisionWins":164}
					{"PlayoffRoundReached":165}
					{"Position":229}
					{"PowerMovesRating":230}
					{"PracticeSquadYears":231}
					{"PresentationId":232}
					{"PressRating":233}
					{"PrevTeamIndex":234}
					{"ProBowlAppearences":235}
					{"PursuitRating":236}
					{"RegressionPoints":237}
					{"ReleaseRating":238}
					{"ReservedRating":239}
					{"RunBlockFinesseRating":240}
					{"RunBlockPowerRating":241}
					{"RunBlockRating":242}
					{"RunningStyleRating":243}
					{"Scheme":244}
					{"SeasonStats":246}
					{"SeasonalGoal":245}
					{"ShortRouteRunningRating":247}
					{"SkillPoints":248}
					{"SpectacularCatchRating":249}
					{"SpeedRating":250}
					{"SpinMoveRating":251}
					{"StaminaRating":252}
					{"StiffArmRating":253}
					{"StrengthRating":254}
					{"SuperBowlWins":255}
					{"TEAM_TYPE":257}
					{"TRAIT_BIGHITTER":269}
					{"TRAIT_CLUTCH":270}
					{"TRAIT_COVER_BALL":271}
					{"TRAIT_DEEP_BALL":272}
					{"TRAIT_DLBULLRUSH":273}
					{"TRAIT_DLSPIN":274}
					{"TRAIT_DLSWIM":275}
					{"TRAIT_DROPOPENPASS":276}
					{"TRAIT_FEETINBOUNDS":277}
					{"TRAIT_FIGHTFORYARDS":278}
					{"TRAIT_FORCE_PASS":279}
					{"TRAIT_HIGHMOTOR":280}
					{"TRAIT_HIGHPOINTCATCH":281}
					{"TRAIT_LBSTYLE":282}
					{"TRAIT_PENALTY":283}
					{"TRAIT_PLAY_BALL":284}
					{"TRAIT_POSSESSIONCATCH":285}
					{"TRAIT_PUMPFAKE":286}
					{"TRAIT_QBSTYLE":287}
					{"TRAIT_SENSE_PRESSURE":288}
					{"TRAIT_SENSE_PRESSURE_MAX":289}
					{"TRAIT_STRIPBALL":290}
					{"TRAIT_TACKLELOW":291}
					{"TRAIT_THROWAWAY":292}
					{"TRAIT_TIGHTSPIRAL":293}
					{"TRAIT_TUCK_RUN":294}
					{"TRAIT_YACCATCH":295}
					{"TackleRating":256}
					{"TeamIndex":258}
					{"ThrowAccuracyDeepRating":259}
					{"ThrowAccuracyMidRating":260}
					{"ThrowAccuracyRating":261}
					{"ThrowAccuracyShortRating":262}
					{"ThrowOnTheRunRating":263}
					{"ThrowPowerRating":264}
					{"ThrowUnderPressureRating":265}
					{"TotalInjuryDuration":266}
					{"ToughnessRating":267}
					{"TradeStatus":268}
					{"TraitDevelopment":296}
					{"TraitPredictability":297}
					{"TruckingRating":298}
					{"WasOnPracticeSquadThisYear":299}
					{"WasPreviouslyInjured":300}
					{"WeeklyGoals":301}
					{"Weight":302}
					{"YearDrafted":303}
					{"YearlyAwardCount":304}
					{"YearsPro":305}
					{"ZoneCoverageRating":306}
				}
				{NAME:<BA 1B 65 >:TdfString:"Player"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1541:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=6:
					{"ConfirmationMessage":0}
					{"DisplayHint":1}
					{"DisplayName":2}
					{"Input":3}
					{"IsConfirmationRequired":4}
					{"Name":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"Command"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1542:
			{:<>:TdfStruct:size=:TdfStruct:size=3:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{NAME:<BA 1B 65 >:TdfString:"Command[]"}
				{TYPE:<D3 9C 25 >:TdfInteger:0x1/1}
			}
		}
		{1769:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x41/65}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=6:
					{"Commands":1}
					{"Control":2}
					{"Description":3}
					{"IsSubmittable":0}
					{"Name":4}
					{"Title":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"UIForm"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1857:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x6E9/1769}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=11:
					{"Commands":1}
					{"Control":2}
					{"DataSource":6}
					{"Description":3}
					{"FilterCacheSize":7}
					{"IsFilter1OnDemand":8}
					{"IsFilter2OnDemand":9}
					{"IsFilterRequired":10}
					{"IsSubmittable":0}
					{"Name":4}
					{"Title":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"UIDataForm"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1858:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x741/1857}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=11:
					{"Commands":1}
					{"Control":2}
					{"DataSource":6}
					{"Description":3}
					{"FilterCacheSize":7}
					{"IsFilter1OnDemand":8}
					{"IsFilter2OnDemand":9}
					{"IsFilterRequired":10}
					{"IsSubmittable":0}
					{"Name":4}
					{"Title":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"UISpreadsheetForm"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1859:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x742/1858}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=13:
					{"Commands":1}
					{"Control":2}
					{"DataSource":6}
					{"Description":3}
					{"FilterCacheSize":7}
					{"IsFilter1OnDemand":8}
					{"IsFilter2OnDemand":9}
					{"IsFilterRequired":10}
					{"IsSubmittable":0}
					{"MaxSelectedItems":11}
					{"MinSelectedItems":12}
					{"Name":4}
					{"Title":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"UIListSelectForm"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1874:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x605/1541}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=9:
					{"ConfirmationMessage":0}
					{"DisplayHint":1}
					{"DisplayName":2}
					{"Flow":6}
					{"Input":3}
					{"IsConfirmationRequired":4}
					{"Name":5}
					{"NavigationAction":7}
					{"NavigationString":8}
				}
				{NAME:<BA 1B 65 >:TdfString:"FlowCommand"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{1883:
			{:<>:TdfStruct:size=:TdfStruct:size=3:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{NAME:<BA 1B 65 >:TdfString:"HistoryEntry"}
				{TYPE:<D3 9C 25 >:TdfInteger:0x1/1}
			}
		}
		{2148:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x605/1541}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=6:
					{"ConfirmationMessage":0}
					{"DisplayHint":1}
					{"DisplayName":2}
					{"Input":3}
					{"IsConfirmationRequired":4}
					{"Name":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"ListCommand"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{2811:
			{:<>:TdfStruct:size=:TdfStruct:size=4:
				{BASE:<8A 1C E5 >:TdfInteger:0x741/1857}
				{DICT:<92 98 F4 >:TdfMap:1-0 size=13:
					{"Commands":1}
					{"Control":2}
					{"DataSource":6}
					{"Description":3}
					{"FilterCacheSize":7}
					{"IsFilter1OnDemand":8}
					{"IsFilter2OnDemand":9}
					{"IsFilterRequired":10}
					{"IsSubmittable":0}
					{"MenuItems":11}
					{"Name":4}
					{"SelectedItem":12}
					{"Title":5}
				}
				{NAME:<BA 1B 65 >:TdfString:"UIMenuForm"}
				{TYPE:<D3 9C 25 >:TdfInteger:0}
			}
		}
		{2812:
			{:<>:TdfStruct:size=:TdfStruct:size=3:
				{BASE:<8A 1C E5 >:TdfInteger:0}
				{NAME:<BA 1B 65 >:TdfString:"UIForm[]"}
				{TYPE:<D3 9C 25 >:TdfInteger:0x1/1}
			}
		}
	}
	{RIBC:<CA 98 A3 >:TdfInteger:0x11/17}
	{ROOT:<CA FB F4 >:TdfInteger:0x5280016/86507542}
	{SIBC:<CE 98 A3 >:TdfInteger:0xE/14}
	{TABL:<D2 18 AC >:TdfMap:0-3 size=15:
		{460:
			{:<>:TdfStruct:size=:TdfStruct:size=2:
				{STID:<CF 4A 64 >:TdfInteger:0x752/1874}
				{TABL:<D2 18 AC >:TdfMap:0-7 size=2:
					{108:
						{:TdfIntegerList:size=1:1953265619,}
					}
						{:<>:TdfStruct:size=:TdfStruct:size=1:
							{VECT:<DA 58 F4 >:TdfList:Type:3: size=1:

								{@PD@:<03 01 00 >:TdfStruct:size=:TdfStruct:size=1:
									{[´íÎó]:<2D 56 69 >:TdfString:"ReadTdf.UnknowType:101[ErrPos=2B8A]"}
								}
							}
						}
					}
				}
			}
		}
	}
}
===========================end
`,
	hexdump: "00 02 44 3C 00 00 08 0C 01 D7 00 01 EB 20 00 00 8F 49 74 04 03 2C 92 5C E3 01 01 00 B2 E8 6D 01 10 45 78 70 65 72 74 20 53 63 6F 75 74 69 6E 67 00 CE E8 6D 01 13 54 72 61 69 74 20 45 78 70 65 72 74 20 53 63 6F 75 74 00 00 92 5C E3 01 01 00 B2 E8 6D 01 20 49 6E 63 72 65 61 73 65 64 20 50 6C 61 79 65 72 20 57 65 65 6B 6C 79 20 47 6F 61 6C 20 58 50 00 CE E8 6D 01 1B 54 72 61 69 74 20 49 6E 63 72 65 61 73 65 64 20 45 78 70 65 72 69 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 13 42 65 63 6F 6D 65 20 50 72 65 64 69 63 74 61 62 6C 65 00 CE E8 6D 01 15 54 72 61 69 74 20 50 72 65 64 69 63 74 61 62 69 6C 69 74 79 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 4D 61 73 74 65 72 20 54 72 61 64 65 20 4E 65 67 6F 74 69 61 74 6F 72 00 CE E8 6D 01 16 54 72 61 69 74 20 54 72 61 64 65 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 51 42 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 51 42 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 52 42 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 52 42 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 57 52 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 57 52 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 54 45 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 54 45 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 4F 4C 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 4F 4C 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 44 4C 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 44 4C 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 4C 42 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 4C 42 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 19 44 42 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 44 42 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 4B 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1D 54 72 61 69 74 20 4B 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 50 20 46 72 65 65 20 41 67 65 6E 63 79 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1D 54 72 61 69 74 20 50 20 46 72 65 65 20 41 67 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 51 42 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 51 42 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 52 42 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 52 42 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 57 52 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 57 52 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 54 45 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 54 45 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 4F 4C 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 4F 4C 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 44 4C 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 44 4C 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 4C 42 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 4C 42 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 12 44 42 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1C 54 72 61 69 74 20 44 42 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 11 4B 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1B 54 72 61 69 74 20 4B 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 11 50 20 54 72 61 69 6E 69 6E 67 20 42 6F 6F 73 74 00 CE E8 6D 01 1B 54 72 61 69 74 20 50 20 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 51 42 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 51 42 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 52 42 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 52 42 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 57 52 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 57 52 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 54 45 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 54 45 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 4F 4C 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 4F 4C 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 44 4C 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 44 4C 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 4C 42 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 4C 42 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 15 44 42 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1C 54 72 61 69 74 20 44 42 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 14 4B 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1B 54 72 61 69 74 20 4B 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 14 50 20 52 65 2D 53 69 67 6E 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1B 54 72 61 69 74 20 50 20 43 6F 6E 74 72 61 63 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 51 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 51 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 52 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 52 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 57 52 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 57 52 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 54 45 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 54 45 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 4F 4C 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 4F 4C 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 44 4C 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 44 4C 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 4C 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 4C 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 18 44 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1E 54 72 61 69 74 20 44 42 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 17 4B 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1D 54 72 61 69 74 20 4B 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 92 5C E3 01 01 00 B2 E8 6D 01 17 50 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 CE E8 6D 01 1D 54 72 61 69 74 20 50 20 52 65 74 69 72 65 6D 65 6E 74 20 49 6E 66 6C 75 65 6E 63 65 00 00 9A FC AD 03 92 98 F4 05 00 03 10 3B 8A 1C E5 00 00 92 98 F4 05 01 00 01 11 49 73 55 73 65 72 43 6F 6E 74 72 6F 6C 6C 65 64 00 00 BA 1B 65 01 0B 55 73 65 72 45 6E 74 69 74 79 00 D3 9C 25 00 00 00 81 01 8A 1C E5 00 00 92 98 F4 05 01 00 01 0E 49 73 53 75 62 6D 69 74 74 61 62 6C 65 00 00 BA 1B 65 01 0D 52 65 73 70 6F 6E 73 65 46 6F 72 6D 00 D3 9C 25 00 00 00 92 01 8A 1C E5 00 00 BA 1B 65 01 05 54 65 61 6D 00 D3 9C 25 00 01 00 93 01 8A 1C E5 00 00 BA 1B 65 01 07 54 65 61 6D 5B 5D 00 D3 9C 25 00 01 00 8E 06 8A 1C E5 00 3B 92 98 F4 05 01 00 B3 04 13 41 63 63 65 6C 65 72 61 74 69 6F 6E 52 61 74 69 6E 67 00 01 04 41 67 65 00 02 0E 41 67 69 6C 69 74 79 52 61 74 69 6E 67 00 03 10 41 77 61 72 65 6E 65 73 73 52 61 74 69 6E 67 00 04 0F 42 43 56 69 73 69 6F 6E 52 61 74 69 6E 67 00 06 0B 42 61 63 6B 67 72 6F 75 6E 64 00 05 14 42 6C 6F 63 6B 53 68 65 64 64 69 6E 67 52 61 74 69 6E 67 00 07 10 42 72 65 61 6B 53 61 63 6B 52 61 74 69 6E 67 00 08 12 42 72 65 61 6B 54 61 63 6B 6C 65 52 61 74 69 6E 67 00 09 0E 43 61 70 74 61 69 6E 73 50 61 74 63 68 00 0A 0C 43 61 72 65 65 72 53 74 61 74 73 00 0B 0F 43 61 72 72 79 69 6E 67 52 61 74 69 6E 67 00 0C 15 43 61 74 63 68 49 6E 54 72 61 66 66 69 63 52 61 74 69 6E 67 00 0E 0F 43 61 74 63 68 69 6E 67 52 61 74 69 6E 67 00 0D 08 43 6F 6C 6C 65 67 65 00 0F 11 43 6F 6E 66 69 64 65 6E 63 65 52 61 74 69 6E 67 00 10 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 30 00 11 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 31 00 12 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 32 00 13 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 33 00 14 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 34 00 15 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 35 00 16 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 36 00 17 0F 43 6F 6E 74 72 61 63 74 42 6F 6E 75 73 37 00 18 18 43 6F 6E 74 72 61 63 74 45 78 74 72 61 59 65 61 72 4F 70 74 69 6F 6E 00 19 0F 43 6F 6E 74 72 61 63 74 4C 65 6E 67 74 68 00 1A 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 30 00 1B 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 31 00 1C 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 32 00 1D 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 33 00 1E 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 34 00 1F 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 35 00 20 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 36 00 21 10 43 6F 6E 74 72 61 63 74 53 61 6C 61 72 79 37 00 22 0F 43 6F 6E 74 72 61 63 74 53 74 61 74 75 73 00 23 0D 43 6F 6E 74 72 61 63 74 59 65 61 72 00 24 22 43 75 72 72 65 6E 74 59 65 61 72 53 65 61 73 6F 6E 45 6E 64 69 6E 67 49 6E 6A 75 72 79 57 65 65 6B 00 25 17 44 65 65 70 52 6F 75 74 65 52 75 6E 6E 69 6E 67 52 61 74 69 6E 67 00 26 10 45 4E 44 50 4C 41 59 5F 41 4E 49 4D 56 41 4C 00 28 12 45 6C 75 73 69 76 65 6E 65 73 73 52 61 74 69 6E 67 00 27 11 45 78 70 65 72 69 65 6E 63 65 50 6F 69 6E 74 73 00 29 13 46 69 6E 65 73 73 65 4D 6F 76 65 73 52 61 74 69 6E 67 00 2A 0A 46 69 72 73 74 4E 61 6D 65 00 2B 0A 47 61 6D 65 53 74 61 74 73 00 2C 07 47 65 6E 64 65 72 00 2D 07 48 65 69 67 68 74 00 2E 0F 48 69 74 50 6F 77 65 72 52 61 74 69 6E 67 00 2F 15 49 6D 70 61 63 74 42 6C 6F 63 6B 69 6E 67 52 61 74 69 6E 67 00 30 0D 49 6E 6A 75 72 79 52 61 74 69 6E 67 00 31 0F 49 6E 6A 75 72 79 53 65 76 65 72 69 74 79 00 32 0D 49 6E 6A 75 72 79 53 74 61 74 75 73 00 33 0B 49 6E 6A 75 72 79 54 79 70 65 00 34 0A 49 73 43 72 65 61 74 65 64 00 35 10 49 73 44 65 6D 61 6E 64 52 65 6C 65 61 73 65 00 36 11 49 73 49 6E 6A 75 72 65 64 52 65 73 65 72 76 65 00 37 09 49 73 4C 65 67 65 6E 64 00 38 11 49 73 55 73 65 72 43 6F 6E 74 72 6F 6C 6C 65 64 00 00 0A 4A 65 72 73 65 79 4E 75 6D 00 39 0F 4A 75 6B 65 4D 6F 76 65 52 61 74 69 6E 67 00 3A 0E 4A 75 6D 70 69 6E 67 52 61 74 69 6E 67 00 3B 13 4B 69 63 6B 41 63 63 75 72 61 63 79 52 61 74 69 6E 67 00 3C 10 4B 69 63 6B 50 6F 77 65 72 52 61 74 69 6E 67 00 3D 11 4B 69 63 6B 52 65 74 75 72 6E 52 61 74 69 6E 67 00 3E 09 4C 61 73 74 4E 61 6D 65 00 3F 1D 4C 61 73 74 53 65 61 73 6F 6E 50 72 6F 42 6F 77 6C 41 70 70 65 61 72 61 6E 63 65 73 00 80 01 1F 4C 61 73 74 59 65 61 72 53 65 61 73 6F 6E 45 6E 64 69 6E 67 49 6E 6A 75 72 79 57 65 65 6B 00 81 01 12 4C 61 74 65 73 74 49 6E 6A 75 72 79 53 74 61 67 65 00 82 01 11 4C 61 74 65 73 74 49 6E 6A 75 72 79 57 65 65 6B 00 83 01 11 4C 61 74 65 73 74 49 6E 6A 75 72 79 59 65 61 72 00 84 01 10 4C 65 61 64 42 6C 6F 63 6B 52 61 74 69 6E 67 00 85 01 0C 4C 65 67 61 63 79 53 63 6F 72 65 00 86 01 12 4D 61 6E 43 6F 76 65 72 61 67 65 52 61 74 69 6E 67 00 87 01 12 4D 61 78 49 6E 6A 75 72 79 44 75 72 61 74 69 6F 6E 00 88 01 19 4D 65 64 69 75 6D 52 6F 75 74 65 52 75 6E 6E 69 6E 67 52 61 74 69 6E 67 00 89 01 1A 4D 65 74 61 4D 6F 72 70 68 5F 41 72 6D 73 42 61 72 79 63 65 6E 74 72 69 63 00 8A 01 13 4D 65 74 61 4D 6F 72 70 68 5F 41 72 6D 73 42 61 73 65 00 8B 01 1A 4D 65 74 61 4D 6F 72 70 68 5F 42 61 63 6B 50 6C 61 74 65 41 6D 6F 75 6E 74 00 8C 01 1B 4D 65 74 61 4D 6F 72 70 68 5F 43 61 6C 66 73 42 61 72 79 63 65 6E 74 72 69 63 00 8D 01 14 4D 65 74 61 4D 6F 72 70 68 5F 43 61 6C 66 73 42 61 73 65 00 8E 01 1B 4D 65 74 61 4D 6F 72 70 68 5F 43 68 65 73 74 42 61 72 79 63 65 6E 74 72 69 63 00 8F 01 14 4D 65 74 61 4D 6F 72 70 68 5F 43 68 65 73 74 42 61 73 65 00 90 01 1A 4D 65 74 61 4D 6F 72 70 68 5F 46 65 65 74 42 61 72 79 63 65 6E 74 72 69 63 00 91 01 13 4D 65 74 61 4D 6F 72 70 68 5F 46 65 65 74 42 61 73 65 00 92 01 1B 4D 65 74 61 4D 6F 72 70 68 5F 46 6C 61 6B 4A 61 63 6B 65 74 41 6D 6F 75 6E 74 00 93 01 1B 4D 65 74 61 4D 6F 72 70 68 5F 47 6C 75 74 65 42 61 72 79 63 65 6E 74 72 69 63 00 94 01 14 4D 65 74 61 4D 6F 72 70 68 5F 47 6C 75 74 65 42 61 73 65 00 95 01 19 4D 65 74 61 4D 6F 72 70 68 5F 47 75 74 42 61 72 79 63 65 6E 74 72 69 63 00 96 01 12 4D 65 74 61 4D 6F 72 70 68 5F 47 75 74 42 61 73 65 00 97 01 19 4D 65 74 61 4D 6F 72 70 68 5F 53 68 6F 75 6C 64 65 72 48 65 69 67 68 74 00 98 01 18 4D 65 74 61 4D 6F 72 70 68 5F 53 68 6F 75 6C 64 65 72 57 69 64 74 68 00 99 01 1C 4D 65 74 61 4D 6F 72 70 68 5F 54 68 69 67 68 73 42 61 72 79 63 65 6E 74 72 69 63 00 9A 01 15 4D 65 74 61 4D 6F 72 70 68 5F 54 68 69 67 68 73 42 61 73 65 00 9B 01 12 4D 69 6E 49 6E 6A 75 72 79 44 75 72 61 74 69 6F 6E 00 9C 01 1B 4F 72 69 67 69 6E 61 6C 41 63 63 65 6C 65 72 61 74 69 6F 6E 52 61 74 69 6E 67 00 9D 01 16 4F 72 69 67 69 6E 61 6C 41 67 69 6C 69 74 79 52 61 74 69 6E 67 00 9E 01 18 4F 72 69 67 69 6E 61 6C 41 77 61 72 65 6E 65 73 73 52 61 74 69 6E 67 00 9F 01 17 4F 72 69 67 69 6E 61 6C 42 43 56 69 73 69 6F 6E 52 61 74 69 6E 67 00 A0 01 1C 4F 72 69 67 69 6E 61 6C 42 6C 6F 63 6B 53 68 65 64 64 69 6E 67 52 61 74 69 6E 67 00 A1 01 18 4F 72 69 67 69 6E 61 6C 42 72 65 61 6B 53 61 63 6B 52 61 74 69 6E 67 00 A2 01 1A 4F 72 69 67 69 6E 61 6C 42 72 65 61 6B 54 61 63 6B 6C 65 52 61 74 69 6E 67 00 A3 01 17 4F 72 69 67 69 6E 61 6C 43 61 72 72 79 69 6E 67 52 61 74 69 6E 67 00 A4 01 1D 4F 72 69 67 69 6E 61 6C 43 61 74 63 68 49 6E 54 72 61 66 66 69 63 52 61 74 69 6E 67 00 A6 01 17 4F 72 69 67 69 6E 61 6C 43 61 74 63 68 69 6E 67 52 61 74 69 6E 67 00 A5 01 1F 4F 72 69 67 69 6E 61 6C 44 65 65 70 52 6F 75 74 65 52 75 6E 6E 69 6E 67 52 61 74 69 6E 67 00 A7 01 1A 4F 72 69 67 69 6E 61 6C 45 6C 75 73 69 76 65 6E 65 73 73 52 61 74 69 6E 67 00 A8 01 1B 4F 72 69 67 69 6E 61 6C 46 69 6E 65 73 73 65 4D 6F 76 65 73 52 61 74 69 6E 67 00 A9 01 17 4F 72 69 67 69 6E 61 6C 48 69 74 50 6F 77 65 72 52 61 74 69 6E 67 00 AA 01 1D 4F 72 69 67 69 6E 61 6C 49 6D 70 61 63 74 42 6C 6F 63 6B 69 6E 67 52 61 74 69 6E 67 00 AB 01 15 4F 72 69 67 69 6E 61 6C 49 6E 6A 75 72 79 52 61 74 69 6E 67 00 AC 01 17 4F 72 69 67 69 6E 61 6C 4A 75 6B 65 4D 6F 76 65 52 61 74 69 6E 67 00 AD 01 16 4F 72 69 67 69 6E 61 6C 4A 75 6D 70 69 6E 67 52 61 74 69 6E 67 00 AE 01 1B 4F 72 69 67 69 6E 61 6C 4B 69 63 6B 41 63 63 75 72 61 63 79 52 61 74 69 6E 67 00 AF 01 18 4F 72 69 67 69 6E 61 6C 4B 69 63 6B 50 6F 77 65 72 52 61 74 69 6E 67 00 B0 01 19 4F 72 69 67 69 6E 61 6C 4B 69 63 6B 52 65 74 75 72 6E 52 61 74 69 6E 67 00 B1 01 18 4F 72 69 67 69 6E 61 6C 4C 65 61 64 42 6C 6F 63 6B 52 61 74 69 6E 67 00 B2 01 1A 4F 72 69 67 69 6E 61 6C 4D 61 6E 43 6F 76 65 72 61 67 65 52 61 74 69 6E 67 00 B3 01 21 4F 72 69 67 69 6E 61 6C 4D 65 64 69 75 6D 52 6F 75 74 65 52 75 6E 6E 69 6E 67 52 61 74 69 6E 67 00 B4 01 16 4F 72 69 67 69 6E 61 6C 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 B5 01 1F 4F 72 69 67 69 6E 61 6C 50 61 73 73 42 6C 6F 63 6B 46 69 6E 65 73 73 65 52 61 74 69 6E 67 00 B6 01 1D 4F 72 69 67 69 6E 61 6C 50 61 73 73 42 6C 6F 63 6B 50 6F 77 65 72 52 61 74 69 6E 67 00 B7 01 18 4F 72 69 67 69 6E 61 6C 50 61 73 73 42 6C 6F 63 6B 52 61 74 69 6E 67 00 B8 01 19 4F 72 69 67 69 6E 61 6C 50 6C 61 79 41 63 74 69 6F 6E 52 61 74 69 6E 67 00 B9 01 1E 4F 72 69 67 69 6E 61 6C 50 6C 61 79 52 65 63 6F 67 6E 69 74 69 6F 6E 52 61 74 69 6E 67 00 BA 01 19 4F 72 69 67 69 6E 61 6C 50 6F 77 65 72 4D 6F 76 65 73 52 61 74 69 6E 67 00 BB 01 14 4F 72 69 67 69 6E 61 6C 50 72 65 73 73 52 61 74 69 6E 67 00 BC 01 16 4F 72 69 67 69 6E 61 6C 50 75 72 73 75 69 74 52 61 74 69 6E 67 00 BD 01 16 4F 72 69 67 69 6E 61 6C 52 65 6C 65 61 73 65 52 61 74 69 6E 67 00 BE 01 1E 4F 72 69 67 69 6E 61 6C 52 75 6E 42 6C 6F 63 6B 46 69 6E 65 73 73 65 52 61 74 69 6E 67 00 BF 01 1C 4F 72 69 67 69 6E 61 6C 52 75 6E 42 6C 6F 63 6B 50 6F 77 65 72 52 61 74 69 6E 67 00 80 02 17 4F 72 69 67 69 6E 61 6C 52 75 6E 42 6C 6F 63 6B 52 61 74 69 6E 67 00 81 02 20 4F 72 69 67 69 6E 61 6C 53 68 6F 72 74 52 6F 75 74 65 52 75 6E 6E 69 6E 67 52 61 74 69 6E 67 00 82 02 1F 4F 72 69 67 69 6E 61 6C 53 70 65 63 74 61 63 75 6C 61 72 43 61 74 63 68 52 61 74 69 6E 67 00 83 02 14 4F 72 69 67 69 6E 61 6C 53 70 65 65 64 52 61 74 69 6E 67 00 84 02 17 4F 72 69 67 69 6E 61 6C 53 70 69 6E 4D 6F 76 65 52 61 74 69 6E 67 00 85 02 16 4F 72 69 67 69 6E 61 6C 53 74 61 6D 69 6E 61 52 61 74 69 6E 67 00 86 02 17 4F 72 69 67 69 6E 61 6C 53 74 69 66 66 41 72 6D 52 61 74 69 6E 67 00 87 02 17 4F 72 69 67 69 6E 61 6C 53 74 72 65 6E 67 74 68 52 61 74 69 6E 67 00 88 02 15 4F 72 69 67 69 6E 61 6C 54 61 63 6B 6C 65 52 61 74 69 6E 67 00 89 02 20 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 41 63 63 75 72 61 63 79 44 65 65 70 52 61 74 69 6E 67 00 8A 02 1F 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 41 63 63 75 72 61 63 79 4D 69 64 52 61 74 69 6E 67 00 8B 02 1C 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 41 63 63 75 72 61 63 79 52 61 74 69 6E 67 00 8C 02 21 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 41 63 63 75 72 61 63 79 53 68 6F 72 74 52 61 74 69 6E 67 00 8D 02 1C 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 4F 6E 54 68 65 52 75 6E 52 61 74 69 6E 67 00 8E 02 19 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 50 6F 77 65 72 52 61 74 69 6E 67 00 8F 02 21 4F 72 69 67 69 6E 61 6C 54 68 72 6F 77 55 6E 64 65 72 50 72 65 73 73 75 72 65 52 61 74 69 6E 67 00 90 02 18 4F 72 69 67 69 6E 61 6C 54 6F 75 67 68 6E 65 73 73 52 61 74 69 6E 67 00 91 02 17 4F 72 69 67 69 6E 61 6C 54 72 75 63 6B 69 6E 67 52 61 74 69 6E 67 00 92 02 1B 4F 72 69 67 69 6E 61 6C 5A 6F 6E 65 43 6F 76 65 72 61 67 65 52 61 74 69 6E 67 00 93 02 0E 4F 76 65 72 61 6C 6C 47 72 61 64 65 30 00 94 02 0E 4F 76 65 72 61 6C 6C 47 72 61 64 65 31 00 95 02 0E 4F 76 65 72 61 6C 6C 47 72 61 64 65 32 00 96 02 0E 4F 76 65 72 61 6C 6C 47 72 61 64 65 33 00 97 02 0E 4F 76 65 72 61 6C 6C 47 72 61 64 65 34 00 98 02 0E 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 99 02 0F 50 4C 59 52 5F 41 53 53 45 54 4E 41 4D 45 00 A7 02 0F 50 4C 59 52 5F 42 41 43 4B 50 4C 41 54 45 00 A8 02 0F 50 4C 59 52 5F 42 49 52 54 48 44 41 54 45 00 A9 02 11 50 4C 59 52 5F 42 52 45 41 54 48 45 52 49 54 45 00 AA 02 0F 50 4C 59 52 5F 43 41 50 53 41 4C 41 52 59 00 AB 02 11 50 4C 59 52 5F 43 41 52 45 45 52 50 48 41 53 45 00 AC 02 11 50 4C 59 52 5F 43 45 4C 45 42 52 41 54 49 4F 4E 00 AD 02 0D 50 4C 59 52 5F 43 4F 4D 4D 45 4E 54 00 AE 02 19 50 4C 59 52 5F 43 4F 4E 53 45 43 59 45 41 52 53 57 49 54 48 54 45 41 4D 00 AF 02 0F 50 4C 59 52 5F 44 52 41 46 54 50 49 43 4B 00 B0 02 10 50 4C 59 52 5F 44 52 41 46 54 52 4F 55 4E 44 00 B1 02 0F 50 4C 59 52 5F 44 52 41 46 54 54 45 41 4D 00 B2 02 09 50 4C 59 52 5F 45 47 4F 00 B3 02 0E 50 4C 59 52 5F 45 59 45 50 41 49 4E 54 00 B4 02 0E 50 4C 59 52 5F 46 41 43 45 4D 41 53 4B 00 B5 02 11 50 4C 59 52 5F 46 4C 41 47 50 52 4F 42 4F 57 4C 00 B6 02 10 50 4C 59 52 5F 46 4C 41 4B 4A 41 43 4B 45 54 00 B7 02 11 50 4C 59 52 5F 47 45 4E 45 52 49 43 48 45 41 44 00 B8 02 14 50 4C 59 52 5F 47 52 41 53 53 4C 45 46 54 45 4C 42 4F 57 00 B9 02 13 50 4C 59 52 5F 47 52 41 53 53 4C 45 46 54 48 41 4E 44 00 BA 02 14 50 4C 59 52 5F 47 52 41 53 53 4C 45 46 54 57 52 49 53 54 00 BB 02 15 50 4C 59 52 5F 47 52 41 53 53 52 49 47 48 54 45 4C 42 4F 57 00 BC 02 14 50 4C 59 52 5F 47 52 41 53 53 52 49 47 48 54 48 41 4E 44 00 BD 02 15 50 4C 59 52 5F 47 52 41 53 53 52 49 47 48 54 57 52 49 53 54 00 BE 02 0F 50 4C 59 52 5F 48 41 49 52 43 4F 4C 4F 52 00 BF 02 10 50 4C 59 52 5F 48 41 4E 44 45 44 4E 45 53 53 00 80 03 10 50 4C 59 52 5F 48 41 4E 44 57 41 52 4D 45 52 00 81 03 0C 50 4C 59 52 5F 48 45 4C 4D 45 54 00 82 03 10 50 4C 59 52 5F 48 4F 4D 45 5F 53 54 41 54 45 00 83 03 0F 50 4C 59 52 5F 48 4F 4D 45 5F 54 4F 57 4E 00 84 03 0A 50 4C 59 52 5F 49 43 4F 4E 00 85 03 0F 50 4C 59 52 5F 49 53 43 41 50 54 41 49 4E 00 86 03 12 50 4C 59 52 5F 4A 45 52 53 45 59 53 4C 45 45 56 45 00 87 03 10 50 4C 59 52 5F 4A 45 52 53 45 59 54 59 50 45 00 88 03 15 50 4C 59 52 5F 4C 41 53 54 48 4F 4C 44 4F 55 54 59 45 41 52 00 89 03 13 50 4C 59 52 5F 4C 45 46 54 41 52 4D 53 4C 45 45 56 45 00 8A 03 0E 50 4C 59 52 5F 4C 45 46 54 4B 4E 45 45 00 8B 03 0E 50 4C 59 52 5F 4C 45 46 54 53 48 4F 45 00 8C 03 0E 50 4C 59 52 5F 4C 45 46 54 53 50 41 54 00 8D 03 0F 50 4C 59 52 5F 4C 45 46 54 54 48 49 47 48 00 8E 03 10 50 4C 59 52 5F 4D 4F 55 54 48 50 49 45 43 45 00 8F 03 0D 50 4C 59 52 5F 4E 45 43 4B 50 41 44 00 90 03 0E 50 4C 59 52 5F 4E 45 43 4B 54 59 50 45 00 91 03 12 50 4C 59 52 5F 50 45 52 46 4F 52 4D 4C 45 56 45 4C 00 92 03 0E 50 4C 59 52 5F 50 4F 52 54 52 41 49 54 00 93 03 10 50 4C 59 52 5F 50 52 45 56 54 45 41 4D 49 44 00 94 03 0D 50 4C 59 52 5F 51 42 53 54 59 4C 45 00 95 03 14 50 4C 59 52 5F 52 49 47 48 54 41 52 4D 53 4C 45 45 56 45 00 96 03 0F 50 4C 59 52 5F 52 49 47 48 54 4B 4E 45 45 00 97 03 0F 50 4C 59 52 5F 52 49 47 48 54 53 48 4F 45 00 98 03 0F 50 4C 59 52 5F 52 49 47 48 54 53 50 41 54 00 99 03 10 50 4C 59 52 5F 52 49 47 48 54 54 48 49 47 48 00 9A 03 17 50 4C 59 52 5F 53 49 44 45 4C 49 4E 45 5F 48 45 41 44 47 45 41 52 00 9B 03 0A 50 4C 59 52 5F 53 4B 49 4E 00 9C 03 17 50 4C 59 52 5F 53 4C 45 45 56 45 54 45 4D 50 45 52 41 54 55 52 45 00 9D 03 11 50 4C 59 52 5F 53 4F 43 4B 5F 48 45 49 47 48 54 00 9E 03 0C 50 4C 59 52 5F 53 54 41 4E 43 45 00 9F 03 0B 50 4C 59 52 5F 53 54 59 4C 45 00 A0 03 0E 50 4C 59 52 5F 54 45 4E 44 45 4E 43 59 00 A1 03 0B 50 4C 59 52 5F 54 4F 57 45 4C 00 A2 03 10 50 4C 59 52 5F 55 4E 44 45 52 53 48 49 52 54 00 A3 03 0B 50 4C 59 52 5F 56 49 53 4F 52 00 A4 03 17 50 61 73 73 42 6C 6F 63 6B 46 69 6E 65 73 73 65 52 61 74 69 6E 67 00 9A 02 15 50 61 73 73 42 6C 6F 63 6B 50 6F 77 65 72 52 61 74 69 6E 67 00 9B 02 10 50 61 73 73 42 6C 6F 63 6B 52 61 74 69 6E 67 00 9C 02 12 50 65 72 73 6F 6E 61 6C 69 74 79 52 61 74 69 6E 67 00 9D 02 11 50 6C 61 79 41 63 74 69 6F 6E 52 61 74 69 6E 67 00 9E 02 16 50 6C 61 79 52 65 63 6F 67 6E 69 74 69 6F 6E 52 61 74 69 6E 67 00 A6 02 0D 50 6C 61 79 65 72 42 6F 74 74 6F 6D 00 9F 02 11 50 6C 61 79 65 72 50 65 72 63 65 6E 74 61 67 65 00 A0 02 0A 50 6C 61 79 65 72 54 6F 70 00 A1 02 0B 50 6C 61 79 65 72 54 79 70 65 00 A2 02 16 50 6C 61 79 6F 66 66 43 6F 6E 66 65 72 65 6E 63 65 57 69 6E 73 00 A3 02 14 50 6C 61 79 6F 66 66 44 69 76 69 73 69 6F 6E 57 69 6E 73 00 A4 02 14 50 6C 61 79 6F 66 66 52 6F 75 6E 64 52 65 61 63 68 65 64 00 A5 02 09 50 6F 73 69 74 69 6F 6E 00 A5 03 11 50 6F 77 65 72 4D 6F 76 65 73 52 61 74 69 6E 67 00 A6 03 13 50 72 61 63 74 69 63 65 53 71 75 61 64 59 65 61 72 73 00 A7 03 0F 50 72 65 73 65 6E 74 61 74 69 6F 6E 49 64 00 A8 03 0C 50 72 65 73 73 52 61 74 69 6E 67 00 A9 03 0E 50 72 65 76 54 65 61 6D 49 6E 64 65 78 00 AA 03 13 50 72 6F 42 6F 77 6C 41 70 70 65 61 72 65 6E 63 65 73 00 AB 03 0E 50 75 72 73 75 69 74 52 61 74 69 6E 67 00 AC 03 11 52 65 67 72 65 73 73 69 6F 6E 50 6F 69 6E 74 73 00 AD 03 0E 52 65 6C 65 61 73 65 52 61 74 69 6E 67 00 AE 03 0F 52 65 73 65 72 76 65 64 52 61 74 69 6E 67 00 AF 03 16 52 75 6E 42 6C 6F 63 6B 46 69 6E 65 73 73 65 52 61 74 69 6E 67 00 B0 03 14 52 75 6E 42 6C 6F 63 6B 50 6F 77 65 72 52 61 74 69 6E 67 00 B1 03 0F 52 75 6E 42 6C 6F 63 6B 52 61 74 69 6E 67 00 B2 03 13 52 75 6E 6E 69 6E 67 53 74 79 6C 65 52 61 74 69 6E 67 00 B3 03 07 53 63 68 65 6D 65 00 B4 03 0C 53 65 61 73 6F 6E 53 74 61 74 73 00 B6 03 0D 53 65 61 73 6F 6E 61 6C 47 6F 61 6C 00 B5 03 18 53 68 6F 72 74 52 6F 75 74 65 52 75 6E 6E 69 6E 67 52 61 74 69 6E 67 00 B7 03 0C 53 6B 69 6C 6C 50 6F 69 6E 74 73 00 B8 03 17 53 70 65 63 74 61 63 75 6C 61 72 43 61 74 63 68 52 61 74 69 6E 67 00 B9 03 0C 53 70 65 65 64 52 61 74 69 6E 67 00 BA 03 0F 53 70 69 6E 4D 6F 76 65 52 61 74 69 6E 67 00 BB 03 0E 53 74 61 6D 69 6E 61 52 61 74 69 6E 67 00 BC 03 0F 53 74 69 66 66 41 72 6D 52 61 74 69 6E 67 00 BD 03 0F 53 74 72 65 6E 67 74 68 52 61 74 69 6E 67 00 BE 03 0E 53 75 70 65 72 42 6F 77 6C 57 69 6E 73 00 BF 03 0A 54 45 41 4D 5F 54 59 50 45 00 81 04 10 54 52 41 49 54 5F 42 49 47 48 49 54 54 45 52 00 8D 04 0D 54 52 41 49 54 5F 43 4C 55 54 43 48 00 8E 04 11 54 52 41 49 54 5F 43 4F 56 45 52 5F 42 41 4C 4C 00 8F 04 10 54 52 41 49 54 5F 44 45 45 50 5F 42 41 4C 4C 00 90 04 11 54 52 41 49 54 5F 44 4C 42 55 4C 4C 52 55 53 48 00 91 04 0D 54 52 41 49 54 5F 44 4C 53 50 49 4E 00 92 04 0D 54 52 41 49 54 5F 44 4C 53 57 49 4D 00 93 04 13 54 52 41 49 54 5F 44 52 4F 50 4F 50 45 4E 50 41 53 53 00 94 04 13 54 52 41 49 54 5F 46 45 45 54 49 4E 42 4F 55 4E 44 53 00 95 04 14 54 52 41 49 54 5F 46 49 47 48 54 46 4F 52 59 41 52 44 53 00 96 04 11 54 52 41 49 54 5F 46 4F 52 43 45 5F 50 41 53 53 00 97 04 10 54 52 41 49 54 5F 48 49 47 48 4D 4F 54 4F 52 00 98 04 15 54 52 41 49 54 5F 48 49 47 48 50 4F 49 4E 54 43 41 54 43 48 00 99 04 0E 54 52 41 49 54 5F 4C 42 53 54 59 4C 45 00 9A 04 0E 54 52 41 49 54 5F 50 45 4E 41 4C 54 59 00 9B 04 10 54 52 41 49 54 5F 50 4C 41 59 5F 42 41 4C 4C 00 9C 04 16 54 52 41 49 54 5F 50 4F 53 53 45 53 53 49 4F 4E 43 41 54 43 48 00 9D 04 0F 54 52 41 49 54 5F 50 55 4D 50 46 41 4B 45 00 9E 04 0E 54 52 41 49 54 5F 51 42 53 54 59 4C 45 00 9F 04 15 54 52 41 49 54 5F 53 45 4E 53 45 5F 50 52 45 53 53 55 52 45 00 A0 04 19 54 52 41 49 54 5F 53 45 4E 53 45 5F 50 52 45 53 53 55 52 45 5F 4D 41 58 00 A1 04 10 54 52 41 49 54 5F 53 54 52 49 50 42 41 4C 4C 00 A2 04 10 54 52 41 49 54 5F 54 41 43 4B 4C 45 4C 4F 57 00 A3 04 10 54 52 41 49 54 5F 54 48 52 4F 57 41 57 41 59 00 A4 04 12 54 52 41 49 54 5F 54 49 47 48 54 53 50 49 52 41 4C 00 A5 04 0F 54 52 41 49 54 5F 54 55 43 4B 5F 52 55 4E 00 A6 04 0F 54 52 41 49 54 5F 59 41 43 43 41 54 43 48 00 A7 04 0D 54 61 63 6B 6C 65 52 61 74 69 6E 67 00 80 04 0A 54 65 61 6D 49 6E 64 65 78 00 82 04 18 54 68 72 6F 77 41 63 63 75 72 61 63 79 44 65 65 70 52 61 74 69 6E 67 00 83 04 17 54 68 72 6F 77 41 63 63 75 72 61 63 79 4D 69 64 52 61 74 69 6E 67 00 84 04 14 54 68 72 6F 77 41 63 63 75 72 61 63 79 52 61 74 69 6E 67 00 85 04 19 54 68 72 6F 77 41 63 63 75 72 61 63 79 53 68 6F 72 74 52 61 74 69 6E 67 00 86 04 14 54 68 72 6F 77 4F 6E 54 68 65 52 75 6E 52 61 74 69 6E 67 00 87 04 11 54 68 72 6F 77 50 6F 77 65 72 52 61 74 69 6E 67 00 88 04 19 54 68 72 6F 77 55 6E 64 65 72 50 72 65 73 73 75 72 65 52 61 74 69 6E 67 00 89 04 14 54 6F 74 61 6C 49 6E 6A 75 72 79 44 75 72 61 74 69 6F 6E 00 8A 04 10 54 6F 75 67 68 6E 65 73 73 52 61 74 69 6E 67 00 8B 04 0C 54 72 61 64 65 53 74 61 74 75 73 00 8C 04 11 54 72 61 69 74 44 65 76 65 6C 6F 70 6D 65 6E 74 00 A8 04 14 54 72 61 69 74 50 72 65 64 69 63 74 61 62 69 6C 69 74 79 00 A9 04 0F 54 72 75 63 6B 69 6E 67 52 61 74 69 6E 67 00 AA 04 1B 57 61 73 4F 6E 50 72 61 63 74 69 63 65 53 71 75 61 64 54 68 69 73 59 65 61 72 00 AB 04 15 57 61 73 50 72 65 76 69 6F 75 73 6C 79 49 6E 6A 75 72 65 64 00 AC 04 0C 57 65 65 6B 6C 79 47 6F 61 6C 73 00 AD 04 07 57 65 69 67 68 74 00 AE 04 0C 59 65 61 72 44 72 61 66 74 65 64 00 AF 04 11 59 65 61 72 6C 79 41 77 61 72 64 43 6F 75 6E 74 00 B0 04 09 59 65 61 72 73 50 72 6F 00 B1 04 13 5A 6F 6E 65 43 6F 76 65 72 61 67 65 52 61 74 69 6E 67 00 B2 04 BA 1B 65 01 07 50 6C 61 79 65 72 00 D3 9C 25 00 00 00 85 18 8A 1C E5 00 00 92 98 F4 05 01 00 06 14 43 6F 6E 66 69 72 6D 61 74 69 6F 6E 4D 65 73 73 61 67 65 00 00 0C 44 69 73 70 6C 61 79 48 69 6E 74 00 01 0C 44 69 73 70 6C 61 79 4E 61 6D 65 00 02 06 49 6E 70 75 74 00 03 17 49 73 43 6F 6E 66 69 72 6D 61 74 69 6F 6E 52 65 71 75 69 72 65 64 00 04 05 4E 61 6D 65 00 05 BA 1B 65 01 08 43 6F 6D 6D 61 6E 64 00 D3 9C 25 00 00 00 86 18 8A 1C E5 00 00 BA 1B 65 01 0A 43 6F 6D 6D 61 6E 64 5B 5D 00 D3 9C 25 00 01 00 A9 1B 8A 1C E5 00 81 01 92 98 F4 05 01 00 06 09 43 6F 6D 6D 61 6E 64 73 00 01 08 43 6F 6E 74 72 6F 6C 00 02 0C 44 65 73 63 72 69 70 74 69 6F 6E 00 03 0E 49 73 53 75 62 6D 69 74 74 61 62 6C 65 00 00 05 4E 61 6D 65 00 04 06 54 69 74 6C 65 00 05 BA 1B 65 01 07 55 49 46 6F 72 6D 00 D3 9C 25 00 00 00 81 1D 8A 1C E5 00 A9 1B 92 98 F4 05 01 00 0B 09 43 6F 6D 6D 61 6E 64 73 00 01 08 43 6F 6E 74 72 6F 6C 00 02 0B 44 61 74 61 53 6F 75 72 63 65 00 06 0C 44 65 73 63 72 69 70 74 69 6F 6E 00 03 10 46 69 6C 74 65 72 43 61 63 68 65 53 69 7A 65 00 07 12 49 73 46 69 6C 74 65 72 31 4F 6E 44 65 6D 61 6E 64 00 08 12 49 73 46 69 6C 74 65 72 32 4F 6E 44 65 6D 61 6E 64 00 09 11 49 73 46 69 6C 74 65 72 52 65 71 75 69 72 65 64 00 0A 0E 49 73 53 75 62 6D 69 74 74 61 62 6C 65 00 00 05 4E 61 6D 65 00 04 06 54 69 74 6C 65 00 05 BA 1B 65 01 0B 55 49 44 61 74 61 46 6F 72 6D 00 D3 9C 25 00 00 00 82 1D 8A 1C E5 00 81 1D 92 98 F4 05 01 00 0B 09 43 6F 6D 6D 61 6E 64 73 00 01 08 43 6F 6E 74 72 6F 6C 00 02 0B 44 61 74 61 53 6F 75 72 63 65 00 06 0C 44 65 73 63 72 69 70 74 69 6F 6E 00 03 10 46 69 6C 74 65 72 43 61 63 68 65 53 69 7A 65 00 07 12 49 73 46 69 6C 74 65 72 31 4F 6E 44 65 6D 61 6E 64 00 08 12 49 73 46 69 6C 74 65 72 32 4F 6E 44 65 6D 61 6E 64 00 09 11 49 73 46 69 6C 74 65 72 52 65 71 75 69 72 65 64 00 0A 0E 49 73 53 75 62 6D 69 74 74 61 62 6C 65 00 00 05 4E 61 6D 65 00 04 06 54 69 74 6C 65 00 05 BA 1B 65 01 12 55 49 53 70 72 65 61 64 73 68 65 65 74 46 6F 72 6D 00 D3 9C 25 00 00 00 83 1D 8A 1C E5 00 82 1D 92 98 F4 05 01 00 0D 09 43 6F 6D 6D 61 6E 64 73 00 01 08 43 6F 6E 74 72 6F 6C 00 02 0B 44 61 74 61 53 6F 75 72 63 65 00 06 0C 44 65 73 63 72 69 70 74 69 6F 6E 00 03 10 46 69 6C 74 65 72 43 61 63 68 65 53 69 7A 65 00 07 12 49 73 46 69 6C 74 65 72 31 4F 6E 44 65 6D 61 6E 64 00 08 12 49 73 46 69 6C 74 65 72 32 4F 6E 44 65 6D 61 6E 64 00 09 11 49 73 46 69 6C 74 65 72 52 65 71 75 69 72 65 64 00 0A 0E 49 73 53 75 62 6D 69 74 74 61 62 6C 65 00 00 11 4D 61 78 53 65 6C 65 63 74 65 64 49 74 65 6D 73 00 0B 11 4D 69 6E 53 65 6C 65 63 74 65 64 49 74 65 6D 73 00 0C 05 4E 61 6D 65 00 04 06 54 69 74 6C 65 00 05 BA 1B 65 01 11 55 49 4C 69 73 74 53 65 6C 65 63 74 46 6F 72 6D 00 D3 9C 25 00 00 00 92 1D 8A 1C E5 00 85 18 92 98 F4 05 01 00 09 14 43 6F 6E 66 69 72 6D 61 74 69 6F 6E 4D 65 73 73 61 67 65 00 00 0C 44 69 73 70 6C 61 79 48 69 6E 74 00 01 0C 44 69 73 70 6C 61 79 4E 61 6D 65 00 02 05 46 6C 6F 77 00 06 06 49 6E 70 75 74 00 03 17 49 73 43 6F 6E 66 69 72 6D 61 74 69 6F 6E 52 65 71 75 69 72 65 64 00 04 05 4E 61 6D 65 00 05 11 4E 61 76 69 67 61 74 69 6F 6E 41 63 74 69 6F 6E 00 07 11 4E 61 76 69 67 61 74 69 6F 6E 53 74 72 69 6E 67 00 08 BA 1B 65 01 0C 46 6C 6F 77 43 6F 6D 6D 61 6E 64 00 D3 9C 25 00 00 00 9B 1D 8A 1C E5 00 00 BA 1B 65 01 0D 48 69 73 74 6F 72 79 45 6E 74 72 79 00 D3 9C 25 00 01 00 A4 21 8A 1C E5 00 85 18 92 98 F4 05 01 00 06 14 43 6F 6E 66 69 72 6D 61 74 69 6F 6E 4D 65 73 73 61 67 65 00 00 0C 44 69 73 70 6C 61 79 48 69 6E 74 00 01 0C 44 69 73 70 6C 61 79 4E 61 6D 65 00 02 06 49 6E 70 75 74 00 03 17 49 73 43 6F 6E 66 69 72 6D 61 74 69 6F 6E 52 65 71 75 69 72 65 64 00 04 05 4E 61 6D 65 00 05 BA 1B 65 01 0C 4C 69 73 74 43 6F 6D 6D 61 6E 64 00 D3 9C 25 00 00 00 BB 2B 8A 1C E5 00 81 1D 92 98 F4 05 01 00 0D 09 43 6F 6D 6D 61 6E 64 73 00 01 08 43 6F 6E 74 72 6F 6C 00 02 0B 44 61 74 61 53 6F 75 72 63 65 00 06 0C 44 65 73 63 72 69 70 74 69 6F 6E 00 03 10 46 69 6C 74 65 72 43 61 63 68 65 53 69 7A 65 00 07 12 49 73 46 69 6C 74 65 72 31 4F 6E 44 65 6D 61 6E 64 00 08 12 49 73 46 69 6C 74 65 72 32 4F 6E 44 65 6D 61 6E 64 00 09 11 49 73 46 69 6C 74 65 72 52 65 71 75 69 72 65 64 00 0A 0E 49 73 53 75 62 6D 69 74 74 61 62 6C 65 00 00 0A 4D 65 6E 75 49 74 65 6D 73 00 0B 05 4E 61 6D 65 00 04 0D 53 65 6C 65 63 74 65 64 49 74 65 6D 00 0C 06 54 69 74 6C 65 00 05 BA 1B 65 01 0B 55 49 4D 65 6E 75 46 6F 72 6D 00 D3 9C 25 00 00 00 BC 2B 8A 1C E5 00 00 BA 1B 65 01 09 55 49 46 6F 72 6D 5B 5D 00 D3 9C 25 00 01 00 CA 98 A3 00 11 CA FB F4 00 96 80 C0 52 CE 98 A3 00 0E D2 18 AC 05 00 03 0F 8C 07 CF 4A 64 00 92 1D D2 18 AC 05 00 07 02 AC 01 01 93 DF E3 C6 0E DA 58 F4 04 03 09 03 01 00 03 2D 56 69 65 77 20 64 65 74 61 69 6C 65 64 20 69 6E 66 6F 72 6D 61 74 69 6F 6E 20 61 62 6F 75 74 20 74 68 69 73 20 70 6C 61 79 65 72 2E 00 03 11 56 69 65 77 20 50 6C 61 79 65 72 20 43 61 72 64 00 01 01 00 00 03 13 50 75 73 68 2E 50 6C 61 79 65 72 44 65 74 61 69 6C 73 00 03 0E 50 6C 61 79 65 72 44 65 74 61 69 6C 73 00 01 02 03 12 45 6E 74 65 72 20 50 6C 61 79 65 72 20 43 61 72 64 00 00 00 AD 01 01 93 DF E3 C6 0E DA 58 F4 04 03 09 03 01 00 03 2B 56 69 65 77 20 70 72 6F 67 72 65 73 73 69 6F 6E 20 70 61 63 6B 61 67 65 73 20 66 6F 72 20 74 68 69 73 20 70 6C 61 79 65 72 2E 00 03 1A 56 69 65 77 20 50 72 6F 67 72 65 73 73 69 6F 6E 20 50 61 63 6B 61 67 65 73 00 01 06 00 00 03 19 50 75 73 68 2E 50 72 6F 67 72 65 73 73 69 6F 6E 50 61 63 6B 61 67 65 73 00 03 14 50 72 6F 67 72 65 73 73 69 6F 6E 50 61 63 6B 61 67 65 73 00 01 02 03 13 50 6C 61 79 65 72 20 50 72 6F 67 72 65 73 73 69 6F 6E 00 00 00 00 93 08 CF 4A 64 00 93 2C D2 18 AC 05 00 07 03 94 01 01 A7 AC CC 8D 0B 8E CC F4 04 00 10 B4 86 B0 4C 8C 80 90 47 B3 85 B0 4C AA 8A B0 4C A8 8A B0 4C AE 8A B0 4C AD 8A B0 4C BD 8A B0 4C B0 8A B0 4C A9 8A B0 4C B2 8A B0 4C B5 8A B0 4C B6 8A B0 4C B7 8A B0 4C B8 8A B0 4C BA 8A B0 4C 92 5C E3 01 1C 53 65 6C 65 63 74 20 61 20 70 6C 61 79 65 72 20 74 6F 20 70 72 6F 67 72 65 73 73 00 92 E8 6D 01 01 00 99 1B E4 00 00 99 2B E4 00 00 9A C8 D1 01 10 46 6C 61 67 67 65 64 50 6F 73 69 74 69 6F 6E 00 9A C8 D2 01 01 00 9A CD 11 04 07 17 01 A8 C3 F7 E4 1C 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 04 41 6C 6C 00 B6 1E 36 00 81 01 B6 9B B6 00 00 00 00 01 A8 C3 F7 E4 1C 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 0C 52 65 63 6F 6D 6D 65 6E 64 65 64 00 B6 1E 36 00 81 01 B6 9B B6 00 21 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 51 42 00 DA CC F4 04 00 02 00 21 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 48 42 00 DA CC F4 04 00 02 01 22 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 46 42 00 DA CC F4 04 00 02 02 23 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 57 52 00 DA CC F4 04 00 02 03 24 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 54 45 00 DA CC F4 04 00 02 04 25 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 4C 54 00 DA CC F4 04 00 02 05 26 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 4C 47 00 DA CC F4 04 00 02 06 27 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 02 43 00 DA CC F4 04 00 02 07 28 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 52 47 00 DA CC F4 04 00 02 08 29 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 52 54 00 DA CC F4 04 00 02 09 2A 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 4C 45 00 DA CC F4 04 00 02 0A 2B 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 52 45 00 DA CC F4 04 00 02 0B 2C 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 44 54 00 DA CC F4 04 00 02 0C 2D 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 05 4C 4F 4C 42 00 DA CC F4 04 00 02 0D 2E 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 04 4D 4C 42 00 DA CC F4 04 00 02 0E 2F 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 05 52 4F 4C 42 00 DA CC F4 04 00 02 0F 30 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 43 42 00 DA CC F4 04 00 02 10 31 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 46 53 00 DA CC F4 04 00 02 11 32 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 03 53 53 00 DA CC F4 04 00 02 12 33 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 02 4B 00 DA CC F4 04 00 02 13 34 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 02 50 00 DA CC F4 04 00 02 14 35 00 00 A7 39 6E 00 01 BA 1B 65 01 07 50 6C 61 79 65 72 00 00 00 95 01 01 A7 AC CC 8D 0B 8E CC F4 04 00 06 B4 86 B0 4C 8C 80 90 47 B3 85 B0 4C A8 8A B0 4C AB 8A B0 4C A9 8A B0 4C 92 5C E3 01 14 53 65 6C 65 63 74 20 61 20 74 6F 70 20 65 61 72 6E 65 72 00 92 E8 6D 01 01 00 99 1B E4 00 00 99 2B E4 00 00 9A C8 D1 01 10 46 6C 61 67 67 65 64 50 6F 73 69 74 69 6F 6E 00 9A C8 D2 01 01 00 9A CD 11 04 07 03 01 A8 C3 F7 E4 1C 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 04 41 6C 6C 00 B6 1E 36 00 81 01 B6 9B B6 00 00 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 4F 66 66 65 6E 73 65 00 DA CC F4 04 00 14 00 01 02 03 04 05 06 07 08 09 1A 1B 1C 1D 1E 1F 20 21 22 23 00 00 01 B5 AC 8D D8 05 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 44 65 66 65 6E 73 65 00 DA CC F4 04 00 12 0A 0B 0C 0D 0E 0F 10 11 12 24 25 26 27 28 29 2A 2B 2C 00 00 A7 39 6E 00 01 BA 1B 65 01 0A 54 6F 70 45 61 72 6E 65 72 00 00 00 96 01 01 A7 AC CC 8D 0B 8E CC F4 04 00 0C B4 8A B0 4C BC 8A B0 4C A8 8A B0 4C AC 8A B0 4C AD 8A B0 4C AE 8A B0 4C AF 8A B0 4C B1 8A B0 4C B3 8A B0 4C B9 8A B0 4C BA 8A B0 4C BB 8A B0 4C 92 5C E3 01 2A 56 69 65 77 20 61 20 73 69 6E 67 6C 65 20 70 6C 61 79 65 72 27 73 20 72 65 67 72 65 73 73 69 6F 6E 20 64 65 74 61 69 6C 73 00 92 E8 6D 01 01 00 99 1B E4 00 00 99 2B E4 00 00 9A C8 D1 01 05 54 65 61 6D 00 9A C8 D2 01 09 50 6F 73 69 74 69 6F 6E 00 9A CD 11 04 07 20 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 06 42 65 61 72 73 00 CB 68 6C 00 82 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 42 65 6E 67 61 6C 73 00 CB 68 6C 00 83 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 06 42 69 6C 6C 73 00 CB 68 6C 00 84 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 42 72 6F 6E 63 6F 73 00 CB 68 6C 00 85 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 42 72 6F 77 6E 73 00 CB 68 6C 00 86 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 0B 42 75 63 63 61 6E 65 65 72 73 00 CB 68 6C 00 87 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 0A 43 61 72 64 69 6E 61 6C 73 00 CB 68 6C 00 88 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 43 68 61 72 67 65 72 73 00 CB 68 6C 00 89 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 43 68 69 65 66 73 00 CB 68 6C 00 8A 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 06 43 6F 6C 74 73 00 CB 68 6C 00 8B 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 43 6F 77 62 6F 79 73 00 CB 68 6C 00 8C 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 44 6F 6C 70 68 69 6E 73 00 CB 68 6C 00 8D 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 01 97 88 EC 00 00 B2 18 AC 01 07 45 61 67 6C 65 73 00 CB 68 6C 00 8E 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 46 61 6C 63 6F 6E 73 00 CB 68 6C 00 8F 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 06 34 39 65 72 73 00 CB 68 6C 00 80 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 47 69 61 6E 74 73 00 CB 68 6C 00 91 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 4A 61 67 75 61 72 73 00 CB 68 6C 00 93 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 05 4A 65 74 73 00 CB 68 6C 00 94 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 06 4C 69 6F 6E 73 00 CB 68 6C 00 95 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 50 61 63 6B 65 72 73 00 CB 68 6C 00 98 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 50 61 6E 74 68 65 72 73 00 CB 68 6C 00 99 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 50 61 74 72 69 6F 74 73 00 CB 68 6C 00 9A 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 52 61 69 64 65 72 73 00 CB 68 6C 00 9B 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 05 52 61 6D 73 00 CB 68 6C 00 9C 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 52 61 76 65 6E 73 00 CB 68 6C 00 9D 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 52 65 64 73 6B 69 6E 73 00 CB 68 6C 00 9E 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 53 61 69 6E 74 73 00 CB 68 6C 00 9F 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 53 65 61 68 61 77 6B 73 00 CB 68 6C 00 A0 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 09 53 74 65 65 6C 65 72 73 00 CB 68 6C 00 A1 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 54 65 78 61 6E 73 00 CB 68 6C 00 A2 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 07 54 69 74 61 6E 73 00 CB 68 6C 00 A3 80 80 E6 05 00 00 01 B5 A2 DB FD 15 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 08 56 69 6B 69 6E 67 73 00 CB 68 6C 00 A4 80 80 E6 05 00 00 9A CD 12 04 07 16 01 A8 C3 F7 E4 1C 92 6B 34 00 00 97 88 EC 00 00 B2 18 AC 01 04 41 6C 6C 00 B6 1E 36 00 20 B6 9B B6 00 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 00 B2 18 AC 01 03 51 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 01 B2 18 AC 01 03 48 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 02 B2 18 AC 01 03 46 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 03 B2 18 AC 01 03 57 52 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 04 B2 18 AC 01 03 54 45 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 05 B2 18 AC 01 03 4C 54 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 06 B2 18 AC 01 03 4C 47 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 07 B2 18 AC 01 02 43 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 08 B2 18 AC 01 03 52 47 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 09 B2 18 AC 01 03 52 54 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 0A B2 18 AC 01 03 4C 45 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 0B B2 18 AC 01 03 52 45 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 0C B2 18 AC 01 03 44 54 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 0D B2 18 AC 01 05 4C 4F 4C 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 0E B2 18 AC 01 04 4D 4C 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 0F B2 18 AC 01 05 52 4F 4C 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 10 B2 18 AC 01 03 43 42 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 11 B2 18 AC 01 03 46 53 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 12 B2 18 AC 01 03 53 53 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 13 B2 18 AC 01 02 4B 00 00 00 01 91 96 85 8D 14 92 6B 34 00 00 97 88 EC 00 00 9B 68 6C 00 14 B2 18 AC 01 02 50 00 00 00 A7 39 6E 00 01 BA 1B 65 01 07 50 6C 61 79 65 72 00 00 00 00 95 08 CF 4A 64 00 A4 21 D2 18 AC 05 00 07 05 22 01 93 DF E3 C6 0E DA 58 F4 04 03 06 03 01 00 03 19 41 75 74 6F 50 72 6F 67 72 65 73 73 20 74 68 69 73 20 70 6C 61 79 65 72 00 03 15 41 75 74 6F 20 50 72 6F 67 72 65 73 73 20 50 6C 61 79 65 72 00 01 03 00 00 03 13 41 75 74 6F 50 72 6F 67 72 65 73 73 50 6C 61 79 65 72 00 00 00 23 01 93 DF E3 C6 0E DA 58 F4 04 03 06 03 01 00 03 1A 50 75 72 63 68 61 73 65 20 73 65 6C 65 63 74 65 64 20 70 61 63 6B 61 67 65 00 03 13 50 75 72 63 68 61 73 65 20 50 61 63 6B 61 67 65 20 31 00 01 01 00 00 03 11 50 75 72 63 68 61 73 65 50 61 63 6B 61 67 65 31 00 00 00 24 01 93 DF E3 C6 0E DA 58 F4 04 03 06 03 01 00 03 1A 50 75 72 63 68 61 73 65 20 73 65 6C 65 63 74 65 64 20 70 61 63 6B 61 67 65 00 03 13 50 75 72 63 68 61 73 65 20 50 61 63 6B 61 67 65 20 32 00 01 01 00 00 03 11 50 75 72 63 68 61 73 65 50 61 63 6B 61 67 65 32 00 00 00 25 01 93 DF E3 C6 0E DA 58 F4 04 03 06 03 01 00 03 1A 50 75 72 63 68 61 73 65 20 73 65 6C 65 63 74 65 64 20 70 61 63 6B 61 67 65 00 03 13 50 75 72 63 68 61 73 65 20 50 61 63 6B 61 67 65 20 33 00 01 01 00 00 03 11 50 75 72 63 68 61 73 65 50 61 63 6B 61 67 65 33 00 00 00 26 01 93 DF E3 C6 0E DA 58 F4 04 03 06 03 01 00 03 24 52 65 6D 6F 76 65 20 70 6C 61 79 65 72 20 66 72 6F 6D 20 72 65 63 6F 6D 6D 65 6E 64 65 64 20 6C 69 73 74 00 03 1A 52 65 6D 6F 76 65 20 52 65 63 6F 6D 6D 65 6E 64 65 64 20 50 6C 61 79 65 72 00 01 02 00 00 03 18 52 65 6D 6F 76 65 52 65 63 6F 6D 6D 65 6E 64 65 64 50 6C 61 79 65 72 00 00 00 00 B9 08 CF 4A 64 00 9E 2C D2 18 AC 05 00 07 01 0C 01 86 C8 AE F2 08 92 1D 30 00 03 92 6B 74 00 00 92 8B B4 01 09 50 6F 73 69 74 69 6F 6E 00 92 E8 6D 01 04 50 4F 53 00 96 ED 22 04 01 82 01 03 51 42 00 03 48 42 00 03 46 42 00 03 57 52 00 03 54 45 00 03 4C 54 00 03 4C 47 00 02 43 00 03 52 47 00 03 52 54 00 03 4C 45 00 03 52 45 00 03 44 54 00 05 4C 4F 4C 42 00 04 4D 4C 42 00 05 52 4F 4C 42 00 03 43 42 00 03 46 53 00 03 53 53 00 02 4B 00 02 50 00 03 4B 52 00 03 50 52 00 04 4B 4F 53 00 03 4C 53 00 05 33 44 52 42 00 05 50 57 48 42 00 05 53 4C 57 52 00 04 52 4C 45 00 04 52 52 45 00 04 52 44 54 00 06 53 55 42 4C 42 00 05 53 4C 43 42 00 03 51 42 00 03 48 42 00 03 46 42 00 03 57 52 00 03 54 45 00 03 4C 54 00 03 4C 47 00 02 43 00 03 52 47 00 03 52 54 00 03 4C 45 00 03 52 45 00 03 44 54 00 05 4C 4F 4C 42 00 04 4D 4C 42 00 05 52 4F 4C 42 00 03 43 42 00 03 46 53 00 03 53 53 00 02 4B 00 02 50 00 03 4B 52 00 03 50 52 00 04 4B 4F 53 00 03 4C 53 00 05 33 44 52 42 00 05 50 57 48 42 00 05 53 4C 57 52 00 04 52 4C 45 00 04 52 52 45 00 04 52 44 54 00 06 53 55 42 4C 42 00 05 53 4C 43 42 00 B6 B9 79 00 00 B6 E8 6D 01 10 46 6C 61 67 67 65 64 50 6F 73 69 74 69 6F 6E 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 8B 01 00 00 00 A3 09 CF 4A 64 00 92 2C D2 18 AC 05 00 07 18 B3 05 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 0E 4F 76 65 72 61 6C 6C 20 47 72 61 64 65 00 92 E8 6D 01 04 4F 56 52 00 B6 B9 79 00 00 B6 E8 6D 01 0E 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 B7 0D 28 01 0E 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 37 00 00 B4 06 01 86 C8 AE F2 08 92 1D 30 00 03 92 6B 74 00 00 92 8B B4 01 0C 50 6C 61 79 65 72 20 4E 61 6D 65 00 92 E8 6D 01 05 4E 61 6D 65 00 B6 B9 79 00 00 B6 E8 6D 01 05 4E 61 6D 65 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 09 4C 61 73 74 4E 61 6D 65 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 84 04 00 00 A8 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 04 41 67 65 00 92 E8 6D 01 04 41 67 65 00 B6 B9 79 00 00 B6 E8 6D 01 04 41 67 65 00 B7 0D 28 01 04 41 67 65 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 00 00 00 A9 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 09 42 6F 6E 75 73 20 58 50 00 92 E8 6D 01 09 42 6F 6E 75 73 20 58 50 00 B6 B9 79 00 00 B6 E8 6D 01 08 42 6F 6E 75 73 58 70 00 B7 0D 28 01 01 00 CE 4A 72 00 41 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 AA 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 12 45 78 70 65 72 69 65 6E 63 65 20 50 6F 69 6E 74 73 00 92 E8 6D 01 03 58 50 00 B6 B9 79 00 00 B6 E8 6D 01 11 45 78 70 65 72 69 65 6E 63 65 50 6F 69 6E 74 73 00 B7 0D 28 01 11 45 78 70 65 72 69 65 6E 63 65 50 6F 69 6E 74 73 00 CE 4A 72 00 41 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 00 00 00 AB 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 12 45 78 70 65 72 69 65 6E 63 65 20 50 6F 69 6E 74 73 00 92 E8 6D 01 03 58 50 00 B6 B9 79 00 00 B6 E8 6D 01 11 45 78 70 65 72 69 65 6E 63 65 50 6F 69 6E 74 73 00 B7 0D 28 01 11 45 78 70 65 72 69 65 6E 63 65 50 6F 69 6E 74 73 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 00 00 00 AC 0A 01 86 C8 AE F2 08 92 1D 30 00 03 92 6B 74 00 00 92 8B B4 01 0A 46 69 72 73 74 4E 61 6D 65 00 92 E8 6D 01 0A 46 69 72 73 74 4E 61 6D 65 00 B6 B9 79 00 00 B6 E8 6D 01 0A 46 69 72 73 74 4E 61 6D 65 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 AD 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 07 48 65 69 67 68 74 00 92 E8 6D 01 07 48 65 69 67 68 74 00 B6 B9 79 00 00 B6 E8 6D 01 07 48 65 69 67 68 74 00 B7 0D 28 01 01 00 CE 4A 72 00 41 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 01 DE 99 34 00 00 00 00 AE 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 09 4A 65 72 73 65 79 20 23 00 92 E8 6D 01 09 4A 65 72 73 65 79 20 23 00 B6 B9 79 00 00 B6 E8 6D 01 0A 4A 65 72 73 65 79 4E 75 6D 00 B7 0D 28 01 01 00 CE 4A 72 00 01 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 AF 0A 01 86 C8 AE F2 08 92 1D 30 00 03 92 6B 74 00 00 92 8B B4 01 09 4C 61 73 74 4E 61 6D 65 00 92 E8 6D 01 09 4C 61 73 74 4E 61 6D 65 00 B6 B9 79 00 00 B6 E8 6D 01 09 4C 61 73 74 4E 61 6D 65 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 B0 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 07 4D 61 78 20 58 50 00 92 E8 6D 01 07 4D 61 78 20 58 50 00 B6 B9 79 00 00 B6 E8 6D 01 06 4D 61 78 58 70 00 B7 0D 28 01 01 00 CE 4A 72 00 41 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 B1 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 0E 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 92 E8 6D 01 0E 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 B6 B9 79 00 00 B6 E8 6D 01 0E 4F 76 65 72 61 6C 6C 52 61 74 69 6E 67 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 8B 01 00 00 B2 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 0B 50 6F 72 74 72 61 69 74 49 64 00 92 E8 6D 01 0B 50 6F 72 74 72 61 69 74 49 64 00 B6 B9 79 00 00 B6 E8 6D 01 0B 50 6F 72 74 72 61 69 74 49 64 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 B3 0A 01 86 C8 AE F2 08 92 1D 30 00 01 92 6B 74 00 00 92 8B B4 01 09 50 6F 73 69 74 69 6F 6E 00 92 E8 6D 01 09 50 6F 73 69 74 69 6F 6E 00 B6 B9 79 00 00 B6 E8 6D 01 09 50 6F 73 69 74 69 6F 6E 00 B7 0D 28 01 01 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 B4 0A 01 86 C8 AE F2 08 92 1D 30 00 03 92 6B 74 00 00 92 8B B4 01 12 52 65 67 72 65 73 73 69 6F 6E 53 75 6D 6D 61 72 79 00 92 E8 6D 01 12 52 65 67 72 65 73 73 69 6F 6E 53 75 6D 6D 61 72 79 00 B6 B9 79 00 00 B6 E8 6D 01 12 52 65 67 72 65 73 73 69 6F 6E 53 75 6D 6D 61 72 79 00 B7 0D 28 01 12 52 65 67 72 65 73 73 69 6F 6E 53 75 6D 6D 61 72 79 00 CE 4A 72 00 00 CF 0C B0 01 01 00 CF 48 74 00 00 DA 9C E9 00 00 DE 99 34 00 00 00 00 B5 0A 01 86 C8 AE F2 08 92 1D 30 00 03 92 6B 74 00 00 92 8B B4 01 0C 53 74 61 74 4F 6E 65 4E 61 6D 65 00 92 E8 6D 01 0C 53 74 61 74 4F 6E 65 4E 61 6D 65 00 B6 B9 79 00 00 B6 E8 6D 01 0C 53 74 61 74 4F 6E 65 4E 61 6D 65 00 B7",
}
