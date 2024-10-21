package tdf

import (
	"unsafe"
)

type BitCast[T any] struct {
	value T
}

func (b *BitCast[T]) Value() T {
	return *(*T)(unsafe.Pointer(b))
}

func (b *BitCast[T]) SetValue(value T) {
	*(*T)(unsafe.Pointer(b)) = value
}
