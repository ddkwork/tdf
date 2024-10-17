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

// NativeTypeBind  for struct filed
var NativeTypeBind = map[reflect.Type]BaseKind{
	reflect.TypeFor[int8]():            IntegerKind,
	reflect.TypeFor[int16]():           IntegerKind,
	reflect.TypeFor[int32]():           IntegerKind,
	reflect.TypeFor[int64]():           IntegerKind,
	reflect.TypeFor[uint8]():           IntegerKind,
	reflect.TypeFor[uint16]():          IntegerKind,
	reflect.TypeFor[uint32]():          IntegerKind,
	reflect.TypeFor[uint64]():          IntegerKind,
	reflect.TypeFor[string]():          StringKind,
	reflect.TypeFor[[]byte]():          BinaryKind, // blob
	reflect.TypeFor[struct{}]():        StructKind,
	reflect.TypeFor[[]any]():           ListKind,
	reflect.TypeFor[map[any]any]():     MapKind,
	reflect.TypeFor[Union]():           UnionKind,
	reflect.TypeFor[Variable]():        VariableKind,
	reflect.TypeFor[float32]():         FloatKind,
	reflect.TypeFor[float64]():         FloatKind,
	reflect.TypeFor[BlazeObjectType](): BlazeObjectTypeKind,
	reflect.TypeFor[BlazeObjectId]():   BlazeObjectIdKind,
	reflect.TypeFor[time.Time]():       TimeValueKind,
}
