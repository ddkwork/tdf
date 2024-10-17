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
}

var CodecMap = map[reflect.Type]Codec{
	reflect.TypeFor[bool]():            {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[int8]():            {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[int16]():           {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[int32]():           {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[int64]():           {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[uint8]():           {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[uint16]():          {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[uint32]():          {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[uint64]():          {baseType: IntegerType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[string]():          {baseType: StringType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[[]byte]():          {baseType: StructType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[struct{}]():        {baseType: StructType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[[]any]():           {baseType: ListType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[map[any]any]():     {baseType: MapType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[Union]():           {baseType: UnionType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[Variable]():        {baseType: VariableType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[float32]():         {baseType: FloatType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[float64]():         {baseType: FloatType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[BlazeObjectType](): {baseType: BlazeObjectTypeType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[BlazeObjectId]():   {baseType: BlazeObjectIdType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
	reflect.TypeFor[time.Time]():       {baseType: TimeValueType, encode: func() {}, decode: func() {}, equal: func() bool { return false }},
}
