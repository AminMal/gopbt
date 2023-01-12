package gopbt

import (
	"reflect"

	"github.com/AminMal/gopbt/gen"
)

type anyGen interface {
	gen.Generator[reflect.Value]
}

func wrap[T any](g gen.Generator[T]) anyGen {
	return &generatorWrapper[T]{g}
}

type generatorWrapper[T any] struct {
	g gen.Generator[T]
}

func (g *generatorWrapper[T]) GenerateOne() reflect.Value {
	return reflect.ValueOf(g.g.GenerateOne())
}

func (g *generatorWrapper[T]) GenerateN(n uint) []reflect.Value {
	values := make([]reflect.Value, n)
	for i, v := range g.g.GenerateN(n) {
		values[i] = reflect.ValueOf(v)
	}
	return values
}
