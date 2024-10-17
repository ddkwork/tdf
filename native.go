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
