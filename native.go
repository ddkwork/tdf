package tdf

import (
	"reflect"
	"time"
)

type (
	Union           byte
	Variable        int
	BlazeObjectType byte
	BlazeObjectId   byte
	TimeValue       time.Time
)

func (b BaseType) Valid() bool { return b >= IntegerType && b < MaxType }

// NativeTypeBind  for struct filed
var NativeTypeBind = map[reflect.Type]BaseType{
	reflect.TypeFor[bool]():            IntegerType,
	reflect.TypeFor[int8]():            IntegerType,
	reflect.TypeFor[int16]():           IntegerType,
	reflect.TypeFor[int32]():           IntegerType,
	reflect.TypeFor[int64]():           IntegerType,
	reflect.TypeFor[uint8]():           IntegerType,
	reflect.TypeFor[uint16]():          IntegerType,
	reflect.TypeFor[uint32]():          IntegerType,
	reflect.TypeFor[uint64]():          IntegerType,
	reflect.TypeFor[string]():          StructType,
	reflect.TypeFor[[]byte]():          StructType,   // blob
	reflect.TypeFor[struct{}]():        StructType,   //java显示tdf就是结构体类型,ID_TERM 0结尾
	reflect.TypeFor[[]any]():           ListType,     //java显示方法名称为Vector,类型，大小，遍历
	reflect.TypeFor[map[any]any]():     MapType,      //k v type,size,遍历
	reflect.TypeFor[Union]():           UnionType,    //todo mock enum
	reflect.TypeFor[Variable]():        VariableType, //VariableTdfContainer ID_TERM 0结尾
	reflect.TypeFor[float32]():         FloatType,    //floatToIntBits
	reflect.TypeFor[float64]():         FloatType,
	reflect.TypeFor[BlazeObjectType](): BlazeObjectTypeType, //getComponentId and getTypeId  整型编解码,难道是树形的层级下标和孩子节点下标?
	reflect.TypeFor[BlazeObjectId]():   BlazeObjectIdType,   //getComponentId , getTypeId and getEntityId 整型编解码
	reflect.TypeFor[time.Time]():       TimeValueType,
}

// 类型 encoderHelper 算法说明，解码也一样,使用函数对称结构，自动验证算法的正确性
type Codec struct {
	baseType BaseType
	encode   func()
	decode   func()
	equal    func() bool
	note     string //算法说明
}

var codecMap = map[BaseType]Codec{
	IntegerType:         {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	StringType:          {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	StructType:          {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	ListType:            {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	MapType:             {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	UnionType:           {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	VariableType:        {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	FloatType:           {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	BlazeObjectTypeType: {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	BlazeObjectIdType:   {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	TimeValueType:       {encode: func() {}, decode: func() {}, equal: func() bool { return false }},
}
var nativeCodecMap = map[reflect.Type]Codec{
	reflect.TypeFor[bool]():            {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[int8]():            {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[int16]():           {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[int32]():           {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[int64]():           {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[uint8]():           {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[uint16]():          {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[uint32]():          {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[uint64]():          {baseType: IntegerType, encode: codecMap[IntegerType].encode, decode: codecMap[IntegerType].decode, equal: codecMap[IntegerType].equal},
	reflect.TypeFor[string]():          {baseType: StringType, encode: codecMap[StringType].encode, decode: codecMap[StringType].decode, equal: codecMap[StringType].equal},
	reflect.TypeFor[[]byte]():          {baseType: StructType, encode: codecMap[StructType].encode, decode: codecMap[StructType].decode, equal: codecMap[StructType].equal},
	reflect.TypeFor[struct{}]():        {baseType: StructType, encode: codecMap[StructType].encode, decode: codecMap[StructType].decode, equal: codecMap[StructType].equal},
	reflect.TypeFor[[]any]():           {baseType: ListType, encode: codecMap[ListType].encode, decode: codecMap[ListType].decode, equal: codecMap[ListType].equal},
	reflect.TypeFor[map[any]any]():     {baseType: MapType, encode: codecMap[MapType].encode, decode: codecMap[MapType].decode, equal: codecMap[MapType].equal},
	reflect.TypeFor[Union]():           {baseType: UnionType, encode: codecMap[UnionType].encode, decode: codecMap[UnionType].decode, equal: codecMap[UnionType].equal},
	reflect.TypeFor[Variable]():        {baseType: VariableType, encode: codecMap[VariableType].encode, decode: codecMap[VariableType].decode, equal: codecMap[VariableType].equal},
	reflect.TypeFor[float32]():         {baseType: FloatType, encode: codecMap[FloatType].encode, decode: codecMap[FloatType].decode, equal: codecMap[FloatType].equal},
	reflect.TypeFor[float64]():         {baseType: FloatType, encode: codecMap[FloatType].encode, decode: codecMap[FloatType].decode, equal: codecMap[FloatType].equal},
	reflect.TypeFor[BlazeObjectType](): {baseType: BlazeObjectTypeType, encode: codecMap[BlazeObjectTypeType].encode, decode: codecMap[BlazeObjectTypeType].decode, equal: codecMap[BlazeObjectTypeType].equal},
	reflect.TypeFor[BlazeObjectId]():   {baseType: BlazeObjectIdType, encode: codecMap[BlazeObjectIdType].encode, decode: codecMap[BlazeObjectIdType].decode, equal: codecMap[BlazeObjectIdType].equal},
	reflect.TypeFor[time.Time]():       {baseType: TimeValueType, encode: codecMap[TimeValueType].encode, decode: codecMap[TimeValueType].decode, equal: codecMap[TimeValueType].equal},
}

func GetCodec(t reflect.Type) (Codec, bool) {
	if c, ok := nativeCodecMap[t]; ok {
		return c, true
	}
	return Codec{}, false
}

func (c Codec) BaseType() BaseType {
	return c.baseType
}

func (c Codec) Encode() {
	c.encode()
}

func (c Codec) Decode() {
	c.decode()
}

func (c Codec) Equal() bool {
	return c.equal()
}

func (c Codec) Note() string {
	return c.note
}

func (c Codec) SetNote(note string) {
	c.note = note
}

// 类型 encoderHelper 算法说明，解码也一样,使用函数对称结构，自动验证算法的正确性
type encoderHelper struct {
	codec Codec
	value interface{}
}

func (e encoderHelper) Encode() {
	e.codec.Encode()
}

func (e encoderHelper) Decode() {
	e.codec.Decode()
}

func (e encoderHelper) Equal() bool {
	return e.codec.Equal()
}

func (e encoderHelper) Note() string {
	return e.codec.Note()
}

func (e encoderHelper) SetNote(note string) {
	e.codec.SetNote(note)
}

func (e encoderHelper) BaseType() BaseType {
	return e.codec.BaseType()
}

func (e encoderHelper) Value() interface{} {
	return e.value
}

func (e *encoderHelper) SetValue(value interface{}) {
	e.value = value
}

func (e encoderHelper) IsStructType() bool {
	return e.codec.BaseType() == StructType
}
