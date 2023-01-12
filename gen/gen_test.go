package gen

import (
	"sort"
	"testing"
	"testing/quick"
)

var globalPropertConf = quick.Config { MaxCount: 1000 }

func checkFailPropery(t *testing.T, err error, unitUnderTestName string, propertyUnderTest string) {
	if err != nil {
		t.Fatalf("%s failed to satisfy the `%s` property: %s", unitUnderTestName, propertyUnderTest, err.Error())
	}
}

type Person struct {
	Name 		string
	Surname 	string
	Age 		int
}

// ------ Between property tests:

func isInBetween[T Numeric](gen Generator[T], value T) bool {
	switch g := gen.(type) {
	case *between[T]:
		return value >= g.min && value < g.max
	default:
		return false
	}
}

func TestBetweenBeingInRange(t *testing.T) {
	betweenBeingInRange := func (min, max, int16, testCaseLength uint16) bool { // refact, maybe use Between as the slice length generator
		r := Between(int(min), int(max))
		for _, v := range r.GenerateN(uint(testCaseLength)) {
			if !isInBetween(r, v) { return false }
		}
		return true
	}

	checkFailPropery(t, quick.Check(betweenBeingInRange, &globalPropertConf), "Between_Numeric", "range satisfaction")
}

func TestBetweenOnly(t *testing.T) {
	isOnly := func(i int, testCaseLength uint8) bool {
		r := Between(i, i)
		for _, v := range r.GenerateN(uint(testCaseLength)) {
			if v != i { return false }
		}
		return true
	}

	checkFailPropery(
		t, quick.Check(isOnly, &globalPropertConf),
		"Between", "between with same min and max is only",
	)
}

func TestMergeBetweens(t *testing.T) {
	twoAdjacnetRangesAreOneBigRange := func (i1, i2, i3 int16, testCaseLength uint8) bool {
		ints := []int{int(i1), int(i2), int(i3)}
		sort.Ints(ints)
		min := ints[0]
		mid := ints[1]
		max := ints[2]

		range1 := Between(min, mid)
		range2 := Between(mid, max)
		wholeRange := Between(min, max)

		for _, num := range wholeRange.GenerateN(uint(testCaseLength)) {
			isEitherInFirstRangeOrSecond := isInBetween(range1, num) || isInBetween(range2, num)
			if !isEitherInFirstRangeOrSecond { return false }
		}

		return true
	}

	checkFailPropery(
		t, quick.Check(twoAdjacnetRangesAreOneBigRange, &globalPropertConf),
		"Between", "two adjacent ranges are one big range",
	)
}

// ------ Only property tests:

func TestOnlyHavingOnlyOneResult(t *testing.T) {
	nameGen := Only("John")
	surnameGen := Only("Doe")
	ageGen := Only(25)

	// when all the fields of a struct can only have one value, then we should only expect one instance as the output
	personGen := UsingGen(nameGen, func (name string) Generator[Person] {
		return UsingGen(surnameGen, func (surname string) Generator[Person] {
			return Using(ageGen, func (age int) Person {
				return Person { name, surname, age }
			})
		})
	})

	expectedPerson := Person { Name: nameGen.GenerateOne(), Surname: surnameGen.GenerateOne(), Age: ageGen.GenerateOne() }

	onlyBeingOnly := func (amount uint16) bool { // todo, maybe use Between as the amount generator
		for _, p := range personGen.GenerateN(uint(amount)) {
			if p != expectedPerson { return false }
		}
		return true
	}

	checkFailPropery(t, quick.Check(onlyBeingOnly, &globalPropertConf), "Only", "only generating only expected")
}

// ------ Composition property tests:

type aggregatorGenerator[T any] struct {
	baseGen		Generator[T]
	aggCount 	int
}

func (a *aggregatorGenerator[T]) GenerateOne() T {
	value := a.baseGen.GenerateOne()
	a.aggCount += 1
	return value
}

func (a *aggregatorGenerator[T]) GenerateN(n uint) []T {
	values := a.baseGen.GenerateN(n)
	a.aggCount += len(values)
	return values
}

func aggregator[T any](g Generator[T]) Generator[T] {
	return &aggregatorGenerator[T]{g, 0}
}

func stackLen[T any](agg Generator[T]) int {
	switch g := agg.(type) {
	case *aggregatorGenerator[T]:
		return g.aggCount
	default:
		return 0
	}
}

func reset[T any](agg Generator[T]) {
	switch g := agg.(type) {
	case *aggregatorGenerator[T]:
		g.aggCount = 0
	default:
		// no action
	}
}

func TestGeneratorsBeingEffects(t *testing.T) {
	nameGen := aggregator(OneOf("John", "Bob"))
	ageGen := aggregator(Between(0, 80))
	surnameGen := aggregator(OneOf("Watson", "Marly"))

	numSubGenerators := 3 // nameGen, surnameGen, ageGen

	personGen := UsingGen(nameGen, func (name string) Generator[Person] {
		return UsingGen(surnameGen, func(surname string) Generator[Person] {
			return Using(ageGen, func (age int) Person {
				return Person {name, surname, age}
			})
		})
	})

	// make sure that no value has been constructed

	valuesGenerated := stackLen(nameGen) + stackLen(surnameGen) + stackLen(ageGen)
	if valuesGenerated != 0 {
		t.Fatalf("composed generators violate the Generators being Effects property.")
	}

	invokingComposedGenCausesEvaluation := func (count uint8) bool {
		personGen.GenerateN(uint(count))
		valuesGenerated = stackLen(nameGen) + stackLen(surnameGen) + stackLen(ageGen)
		reset(nameGen)
		reset(surnameGen)
		reset(ageGen)

		return valuesGenerated == int(count) * numSubGenerators
	}

	checkFailPropery(
		t, quick.Check(invokingComposedGenCausesEvaluation, &globalPropertConf),
		"Generator", "generators being effects",
	)
}
