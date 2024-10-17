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
var NativeTypeBind = map[reflect.Kind]BaseKind{
	reflect.Int8:                       IntegerKind,
	reflect.Int16:                      IntegerKind,
	reflect.Int32:                      IntegerKind,
	reflect.Int64:                      IntegerKind,
	reflect.Uint8:                      IntegerKind,
	reflect.Uint16:                     IntegerKind,
	reflect.Uint32:                     IntegerKind,
	reflect.Uint64:                     IntegerKind,
	reflect.String:                     StringKind,
	reflect.Slice:                      BinaryKind, //[]byte blob BINARY list
	reflect.Struct:                     StructKind,
	reflect.Array:                      ListKind,
	reflect.Map:                        MapKind,
	reflect.TypeFor[Union]().Kind():    UnionKind,
	reflect.TypeFor[Variable]().Kind(): VariableKind,
	reflect.Float32:                    FloatKind,
	reflect.Float64:                    FloatKind,
	reflect.TypeFor[BlazeObjectType]().Kind(): BlazeObjectTypeKind,
	reflect.Interface:                         BlazeObjectIdKind,
	reflect.TypeFor[time.Time]().Kind():       TimeValueKind,
}
