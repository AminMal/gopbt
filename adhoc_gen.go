package gopbt

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/AminMal/gopbt/gen"
)

const complexSize = 50

type structAdhocGenerator struct {
	s                *Session
	t                reflect.Type
	structFieldTypes []reflect.Type
}

func (sag *structAdhocGenerator) GenerateOne() reflect.Value {
	v := reflect.New(sag.t).Elem()

	for i, ft := range sag.structFieldTypes {
		if ft.Kind() != reflect.Struct {
			fieldValue, ok := sag.s.sizedValue(ft, complexSize)
			if !ok {
				panic(fmt.Errorf("cannot generate value of type `%s`", ft.Name()))
			}
			v.Field(i).Set(fieldValue)
			continue
		}
		g, fieldValue, ok := sag.s.generateSizedGeneratorAndValue(ft, complexSize)
		if !ok {
			panic(fmt.Errorf("cannot generate value of type `%s`", ft.Name()))
		}
		sag.s.mapping.setGenerator(ft.Name(), g)
		v.Field(i).Set(fieldValue)
	}
	return v
}

func (sag *structAdhocGenerator) GenerateN(n uint) []reflect.Value {
	values := make([]reflect.Value, n, n)
	for i := uint(0); i < n; i++ {
		values[i] = sag.GenerateOne()
	}
	return values
}

func (s *Session) adhocValueGenerator(t reflect.Type, size int) (anyGen, bool) {
	if _, ok := s.sizedValue(t, size); !ok {
		return nil, false // if we cannot instantiate now, we cannot also create generators
	} else {
		// we're sure that we can create instances now, we can safely ignore the `ok` in adhocGenerator.Generate functions
		fieldTypes := make([]reflect.Type, t.NumField(), t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fieldTypes[i] = t.Field(i).Type
		}
		return &structAdhocGenerator{s, t, fieldTypes}, true
	}
}

func (s *Session) generateSizedGeneratorAndValue(t reflect.Type, size int) (gen anyGen, value reflect.Value, ok bool) {
	if _, ok2 := s.sizedValue(t, size); !ok2 {
		return
	} else {
		fieldTypes := make([]reflect.Type, t.NumField(), t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fieldTypes[i] = t.Field(i).Type
		}
		gen = &structAdhocGenerator{s, t, fieldTypes}
		value = gen.GenerateOne()
		ok = true
		return
	}
}

// sizedValue is almost the same as sizedValue in testing/quick

func (s *Session) sizedValue(t reflect.Type, size int) (value reflect.Value, ok bool) {
	v := reflect.New(t).Elem()

	if g, alreadySupports := s.getGeneratorFor(t); alreadySupports {
		return g.GenerateOne(), true
	}

	switch concrete := t; concrete.Kind() {
	case reflect.Bool:
		v.SetBool(rand.Int()&1 == 0)
	case reflect.Float32:
		v.SetFloat(float64(gen.ArbitraryFloat32.GenerateOne()))
	case reflect.Float64:
		v.SetFloat(gen.ArbitraryFloat64.GenerateOne())
	case reflect.Complex64:
		v.SetComplex(
			complex(float64(gen.ArbitraryFloat32.GenerateOne()), float64(gen.ArbitraryFloat32.GenerateOne())),
		)
	case reflect.Complex128:
		v.SetComplex(
			complex(gen.ArbitraryFloat64.GenerateOne(), gen.ArbitraryFloat64.GenerateOne()),
		)
	case reflect.Int16:
		v.SetInt(int64(gen.ArbitraryInt64.GenerateOne()))
	case reflect.Int32:
		v.SetInt(int64(gen.ArbitraryInt64.GenerateOne()))
	case reflect.Int64:
		v.SetInt(int64(gen.ArbitraryInt64.GenerateOne()))
	case reflect.Int8:
		v.SetInt(int64(gen.ArbitraryInt64.GenerateOne()))
	case reflect.Int:
		v.SetInt(int64(gen.ArbitraryInt64.GenerateOne()))
	case reflect.Uint16:
		v.SetUint(gen.ArbitraryUint64.GenerateOne())
	case reflect.Uint32:
		v.SetUint(gen.ArbitraryUint64.GenerateOne())
	case reflect.Uint64:
		v.SetUint(gen.ArbitraryUint64.GenerateOne())
	case reflect.Uint8:
		v.SetUint(gen.ArbitraryUint64.GenerateOne())
	case reflect.Uint:
		v.SetUint(gen.ArbitraryUint64.GenerateOne())
	case reflect.Uintptr:
		v.SetUint(gen.ArbitraryUint64.GenerateOne())
	case reflect.Map:
		numElems := gen.Between(0, size).GenerateOne()
		v.Set(reflect.MakeMap(concrete))
		for i := 0; i < numElems; i++ {
			key, ok1 := s.sizedValue(concrete.Key(), size)
			value, ok2 := s.sizedValue(concrete.Elem(), size)
			if !ok1 || !ok2 {
				return reflect.Value{}, false
			}
			v.SetMapIndex(key, value)
		}
	case reflect.Pointer:
		if gen.Between(0, size).GenerateOne() == 0 {
			v.Set(reflect.Zero(concrete)) // Generate nil pointer.
		} else {
			elem, ok := s.sizedValue(concrete.Elem(), size)
			if !ok {
				return reflect.Value{}, false
			}
			v.Set(reflect.New(concrete.Elem()))
			v.Elem().Set(elem)
		}
	case reflect.Slice:
		numElems := gen.Between(0, size).GenerateOne()
		sizeLeft := size - numElems
		v.Set(reflect.MakeSlice(concrete, numElems, numElems))
		for i := 0; i < numElems; i++ {
			elem, ok := s.sizedValue(concrete.Elem(), sizeLeft)
			if !ok {
				return reflect.Value{}, false
			}
			v.Index(i).Set(elem)
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem, ok := s.sizedValue(concrete.Elem(), size)
			if !ok {
				return reflect.Value{}, false
			}
			v.Index(i).Set(elem)
		}
	case reflect.String:
		v.SetString(defaultStringGen.GenerateOne())
	case reflect.Struct:
		n := v.NumField()
		// Divide sizeLeft evenly among the struct fields.
		sizeLeft := size
		if n > sizeLeft {
			sizeLeft = 1
		} else if n > 0 {
			sizeLeft /= n
		}
		for i := 0; i < n; i++ {
			elem, ok := s.sizedValue(concrete.Field(i).Type, sizeLeft)
			if !ok {
				return reflect.Value{}, false
			}
			v.Field(i).Set(elem)
		}
	case reflect.Func:
		// todo: add function implementation support!
		panic("todo: add function implementation support!")
	default:
		return reflect.Value{}, false
	}

	return v, true
}
