package gen

import (
	"fmt"
	"math/rand"
)

type Generator[T any] interface {
	GenerateOne() T
	GenerateN(n uint) []T
}

// ------ exact generator ------

type only[T any] struct {
	value T
}

func (o *only[T]) GenerateOne() T { return o.value }

func (o *only[T]) GenerateN(n uint) []T {
	res := make([]T, n)
	for i := uint(0); i < n; i++ {
		res[i] = o.value
	}
	return res
}

func Only[T any](t T) Generator[T] { return &only[T]{value: t} }

// ------ random selector ------

type oneOf[T any] struct {
	values []T
}

func (o *oneOf[T]) genOne(length int) T {
	return o.values[rand.Intn(length)]
}

func (o *oneOf[T]) GenerateOne() T {
	return o.genOne(len(o.values))
}

func (o *oneOf[T]) GenerateN(n uint) []T {
	res := make([]T, n)
	length := len(o.values)
	for i := uint(0); i < n; i++ {
		res[i] = o.genOne(length)
	}
	return res
}

func OneOf[T any](values ...T) Generator[T] {
	return &oneOf[T]{values: values}
}

// todo, add support for complex numbers
type Numeric interface {
	uint8 | uint16 | uint32 | uint64 | uint | int8 | int16 | int32 | int64 | int | float32 | float64
}

func numericMin[T Numeric](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func numericMax[T Numeric](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// ------ range selector ------

type between[T Numeric] struct {
	min, max T
}

func (r *between[T]) GenerateOne() T {
	switch diff := any(r.max - r.min).(type) {
	case uint8:
		return any(randUint8(diff)).(T) + r.min
	case uint16:
		return any(randUint16(diff)).(T) + r.min
	case uint32:
		return any(randUint32(diff)).(T) + r.min
	case uint64:
		return any(randUint64(diff)).(T) + r.min
	case uint:
		return any(randUint(diff)).(T) + r.min
	case int8:
		return any(randInt8(diff)).(T) + r.min
	case int16:
		return any(randInt16(diff)).(T) + r.min
	case int32:
		return any(randInt32(diff)).(T) + r.min
	case int64:
		return any(randInt64(diff)).(T) + r.min
	case int:
		return any(randInt(diff)).(T) + r.min
	case float32:
		return any(randFloat32(diff)).(T) + r.min
	case float64:
		return any(randFloat64(diff)).(T) + r.min
	default:
		panic(fmt.Errorf("match error: unrecognized Numeric type %t", diff))
	}
}

func (r *between[T]) GenerateN(n uint) []T {
	res := make([]T, n)
	for i := uint(0); i < n; i++ {
		res[i] = r.GenerateOne()
	}
	return res
}

func Between[T Numeric](min, max T) Generator[T] {
	actualMin := numericMin(min, max)
	actualMax := numericMax(min, max)

	if actualMin == actualMax {
		return Only(min)
	}
	return &between[T]{actualMin, actualMax}
}
