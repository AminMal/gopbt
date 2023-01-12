package gen

type lazyGen[K any] struct {
	genOneFunc func() K
}

func (lg lazyGen[K]) GenerateOne() K { return lg.genOneFunc() }

func (lg lazyGen[K]) GenerateN(n uint) []K {
	res := make([]K, n)
	for i := uint(0); i < n; i++ {
		res[i] = lg.genOneFunc()
	}
	return res
}

func Using[T any, K any](gen Generator[T], compositionAction func(T) K) Generator[K] {
	return lazyGen[K]{genOneFunc: func() K { return compositionAction(gen.GenerateOne()) }}
}

type flattenedLazyGen[K any, T any] struct {
	tGen Generator[T]
	gen  func(T) Generator[K]
}

func (f flattenedLazyGen[K, T]) GenerateOne() K {
	tInstance := f.tGen.GenerateOne()
	return f.gen(tInstance).GenerateOne()
}

func (f flattenedLazyGen[K, T]) GenerateN(n uint) []K {
	res := make([]K, n)
	ts := f.tGen.GenerateN(n)
	for i := uint(0); i < n; i++ {
		res[i] = f.gen(ts[i]).GenerateOne()
	}
	return res
}

func UsingGen[T any, K any](gen Generator[T], flatMapFunc func(T) Generator[K]) Generator[K] {
	return flattenedLazyGen[K, T]{gen, flatMapFunc}
}
