package gopbt

import (
	"fmt"
	"reflect"
	"testing/quick"

	"github.com/AminMal/gopbt/gen"
)

type Session struct {
	mapping *typeGenMapping

	// SupportAdhocGenerators means that the program runtime can create generators for types that don't have any generators available
	SupportAdhocGenerators bool
}

func (s *Session) getGeneratorFor(t reflect.Type) (anyGen, bool) {
	g, ok := s.mapping.generatorMapping[t.Name()]
	return g, ok
}

func NewSessionWithPrimitives() *Session {
	return &Session{mapping: &typeGenMapping{primitiveGenerators}}
}

func NewSession() *Session {
	return &Session{mapping: &typeGenMapping{make(map[string]anyGen)}}
}

func SetGen[T any](s *Session, g gen.Generator[T]) {
	typeName := reflect.TypeOf(*new(T)).Name()
	s.mapping.setGenerator(typeName, wrap(g))
}

func functionAndType(f any) (v reflect.Value, t reflect.Type, ok bool) {
	v = reflect.ValueOf(f)
	ok = v.Kind() == reflect.Func
	if !ok {
		return
	}
	t = v.Type()
	return
}

func toInterfaces(values []reflect.Value) []any {
	// Copy-paste from testing/quick
	ret := make([]any, len(values))
	for i, v := range values {
		ret[i] = v.Interface()
	}
	return ret
}

func (s *Session) arbitraryValues(args []reflect.Value, f reflect.Type, config *quick.Config) (err error) {
	for j := 0; j < len(args); j++ {
		// todo, check if it's slice, then we can either lookup, or generate based on the base type
		correspondingArgType := f.In(j)
		if gen, ok := s.mapping.generatorMapping[correspondingArgType.Name()]; ok {
			args[j] = gen.GenerateOne()
		} else if s.SupportAdhocGenerators {
			g, canGenerateGenerator := s.adhocValueGenerator(correspondingArgType, complexSize)
			if !canGenerateGenerator {
				err = quick.SetupError(fmt.Sprintf("cannot generate gen.Generator[%s] (argument order: %d)", correspondingArgType, j))
				return
			}
			s.mapping.setGenerator(correspondingArgType.Name(), g)
			args[j] = g.GenerateOne()
		} else {
			err = quick.SetupError(fmt.Sprintf("no generator found for type %s (argument order: %d)", correspondingArgType, j))
			return
		}
	}

	return
}

func getMaxCount(c *quick.Config) (maxCount int) {
	maxCount = c.MaxCount
	if maxCount == 0 {
		if c.MaxCountScale != 0 {
			maxCount = int(c.MaxCountScale * float64(*defaultMaxCount))
		} else {
			maxCount = *defaultMaxCount
		}
	}

	return
}

func validateFunctionType(fType reflect.Type) error {
	if fType.NumOut() != 1 {
		return quick.SetupError("function does not return one value")
	}
	if fType.Out(0).Kind() != reflect.Bool {
		return quick.SetupError("function does not return a bool")
	}
	return nil
}

func (s *Session) Check(f any, conf *quick.Config) error {
	if conf == nil {
		conf = &defaultConfig
	}

	fVal, fType, ok := functionAndType(f)
	if !ok {
		return quick.SetupError("argument is not a function")
	}

	if functionValidationErr := validateFunctionType(fType); functionValidationErr != nil {
		return functionValidationErr
	}

	arguments := make([]reflect.Value, fType.NumIn())
	maxCount := getMaxCount(conf)

	for i := 0; i < maxCount; i++ {
		err := s.arbitraryValues(arguments, fType, conf)
		if err != nil {
			return err
		}

		if !fVal.Call(arguments)[0].Bool() {
			return &quick.CheckError{Count: i + 1, In: toInterfaces(arguments)}
		}
	}

	return nil
}
