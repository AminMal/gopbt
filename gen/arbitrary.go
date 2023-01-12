package gen

import "math"

// ------ int types ------
var ArbitraryInt Generator[int] = Between((math.MinInt/2 + 1), (math.MaxInt/2 - 1))

var ArbitraryInt32 Generator[int32] = Between(int32(math.MinInt32)/2+1, int32(math.MaxInt32)/2-1)

var ArbitraryInt64 Generator[int64]  = Between(int64(math.MinInt64)/2+1, int64(math.MaxInt64)/2-1)

// ------ uint types ------
var ArbitraryUint Generator[uint] = Between(uint(0), uint(math.MaxUint))

var ArbitraryUint8 Generator[uint8] = Between(uint8(0), uint8(math.MaxUint8))

var ArbitraryUint16 Generator[uint16] = Between(uint16(0), uint16(math.MaxUint16))

var ArbitraryUint32 Generator[uint32] = Between(uint32(0), uint32(math.MaxUint32))

var ArbitraryUint64 Generator[uint64] = Between(uint64(0), uint64(math.MaxUint64))

// ------ float types ------
var ArbitraryFloat32 Generator[float32] = Between(float32(math.MinInt32)/2+1, float32(math.MaxInt32)/2-1)

var ArbitraryFloat64 Generator[float64] = Between(float64(math.MinInt64)/2+1, float64(math.MaxInt64)/2-1)

// ------ rune ------
var ArbitraryRune Generator[rune] = ArbitraryInt32

// ------ string ------

type stringGen struct {
	alphabet             []rune
	minLength, maxLength int
}

func (s *stringGen) GenerateOne() string {
	strlen := Between(s.minLength, s.maxLength).GenerateOne()
	rs := OneOf(s.alphabet...).GenerateN(uint(strlen))
	return string(rs)
}

func (s *stringGen) GenerateN(n uint) []string {
	res := make([]string, n)
	for i := uint(0); i < n; i++ {
		res[i] = s.GenerateOne()
	}
	return res
}

func StringGen(alphabet string, minLength uint, maxLength uint) Generator[string] {
	actualMin := numericMin(minLength, maxLength)
	actualMax := numericMax(minLength, maxLength)

	return &stringGen{[]rune(alphabet), int(actualMin), int(actualMax)}
}
