package gen

import "math/rand"

func randUint8(n uint8) uint8 {
	return uint8(rand.Int31n(int32(n)))
}

func randUint16(n uint16) uint16 {
	return uint16(rand.Int31n(int32(n)))
}

func randUint32(n uint32) uint32 {
	return uint32(rand.Int31n(int32(n)))
}

func randUint64(n uint64) uint64 {
	return uint64(rand.Int63n(int64(n)))
}

func randUint(n uint) uint {
	if n <= 1<<32-1 {
		return uint(randUint32(uint32(n)))
	} else {
		return uint(randUint64(uint64(n)))
	}
}

func randInt8(n int8) int8 {
	return int8(rand.Int31n(int32(n)))
}

func randInt16(n int16) int16 {
	return int16(rand.Int31n(int32(n)))
}

func randInt32(n int32) int32 {
	return rand.Int31n(n)
}

func randInt64(n int64) int64 {
	return rand.Int63n(n)
}

func randInt(n int) int {
	return rand.Intn(n)
}

func randFloat32(n float32) float32 {
	return rand.Float32() * n
}

func randFloat64(n float64) float64 {
	return rand.Float64() * n
}
