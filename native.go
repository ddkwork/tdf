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
	reflect.TypeFor[int8]():            IntegerType,
	reflect.TypeFor[int16]():           IntegerType,
	reflect.TypeFor[int32]():           IntegerType,
	reflect.TypeFor[int64]():           IntegerType,
	reflect.TypeFor[uint8]():           IntegerType,
	reflect.TypeFor[uint16]():          IntegerType,
	reflect.TypeFor[uint32]():          IntegerType,
	reflect.TypeFor[uint64]():          IntegerType,
	reflect.TypeFor[string]():          StructType,
	reflect.TypeFor[[]byte]():          StructType, // blob
	reflect.TypeFor[struct{}]():        StructType,
	reflect.TypeFor[[]any]():           ListType,
	reflect.TypeFor[map[any]any]():     MapType,
	reflect.TypeFor[Union]():           UnionType,
	reflect.TypeFor[Variable]():        VariableType,
	reflect.TypeFor[float32]():         FloatType,
	reflect.TypeFor[float64]():         FloatType,
	reflect.TypeFor[BlazeObjectType](): BlazeObjectTypeType,
	reflect.TypeFor[BlazeObjectId]():   BlazeObjectIdType,
	reflect.TypeFor[time.Time]():       TimeValueType,
}
