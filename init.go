package gopbt

import (
	"flag"
	"reflect"
	"testing/quick"

	"github.com/AminMal/gopbt/gen"
)

// todo, add these to init
var defaultMaxCount *int = flag.Int("gopbtchecks", 100, "The default number of iterations for each check")
var defaultConfig quick.Config

// todo, add this to init
var primitiveGenerators map[string]anyGen

var defaultAlphabet string
var defaultStringGen gen.Generator[string]

func init() {
	defaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_=+?/`~\"\\:;"
	defaultStringGen = gen.StringGen(defaultAlphabet, uint(0), uint(complexSize))

	primitiveGenerators = map[string]anyGen{
		reflect.TypeOf(0).Name():          wrap(gen.ArbitraryInt),
		reflect.TypeOf(int32(0)).Name():   wrap(gen.ArbitraryInt32),
		reflect.TypeOf(int64(0)).Name():   wrap(gen.ArbitraryInt64),
		reflect.TypeOf(uint(0)).Name():    wrap(gen.ArbitraryUint),
		reflect.TypeOf(uint16(0)).Name():  wrap(gen.ArbitraryUint16),
		reflect.TypeOf(uint32(0)).Name():  wrap(gen.ArbitraryUint32),
		reflect.TypeOf(uint64(0)).Name():  wrap(gen.ArbitraryUint64),
		reflect.TypeOf(float32(0)).Name(): wrap(gen.ArbitraryFloat32),
		reflect.TypeOf(float64(0)).Name(): wrap(gen.ArbitraryFloat64),
		reflect.TypeOf('r').Name():        wrap(gen.ArbitraryRune),
		reflect.TypeOf("").Name():         wrap(defaultStringGen),
	}
}
