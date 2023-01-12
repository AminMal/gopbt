package gen

import "math"

// ------ int types ------
func ArbitraryInt() Generator[int] { return Between((math.MinInt/2 + 1), (math.MaxInt/2 - 1)) }

func ArbitraryInt32() Generator[int32] {
	return Between(int32(math.MinInt32)/2+1, int32(math.MaxInt32)/2-1)
}

func ArbitraryInt64() Generator[int64] {
	return Between(int64(math.MinInt64)/2+1, int64(math.MaxInt64)/2-1)
}

// ------ uint types ------
func ArbitraryUint() Generator[uint] { return Between(uint(0), uint(math.MaxUint)) }

func ArbitraryUint8() Generator[uint8] { return Between(uint8(0), uint8(math.MaxUint8)) }

func ArbitraryUint16() Generator[uint16] { return Between(uint16(0), uint16(math.MaxUint16)) }

func ArbitraryUint32() Generator[uint32] { return Between(uint32(0), uint32(math.MaxUint32)) }

func ArbitraryUint64() Generator[uint64] { return Between(uint64(0), uint64(math.MaxUint64)) }

// ------ float types ------
func ArbitraryFloat32() Generator[float32] {
	return Between(float32(math.MinInt32)/2+1, float32(math.MaxInt32)/2-1)
}

func ArbitraryFloat64() Generator[float64] {
	return Between(float64(math.MinInt64)/2+1, float64(math.MaxInt64)/2-1)
}

// ------ rune ------
func ArbitraryRune() Generator[rune] { return ArbitraryInt32() }

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
