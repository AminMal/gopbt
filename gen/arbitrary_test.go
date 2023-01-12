package gen

import (
	"testing"
	"testing/quick"
)

var globalPropertConf = quick.Config { MaxCount: 1000 }

func checkFailPropery(t *testing.T, err error, unitUnderTestName string, propertyUnderTest string) {
	if err != nil {
		t.Fatalf("%s failed to satisfy the `%s` property: %s", unitUnderTestName, propertyUnderTest, err.Error())
	}
}

// ------ Between property tests:

func TestBetweenBeingInRange(t *testing.T) {
	betweenBeingInRange := func (min, max, int, testCaseLength uint16) bool { // refact, maybe use Between as the slice length generator
		actualMin := numericMin(min, max)
		actualMax := numericMax(min, max)
		for _, v := range Between(actualMin, actualMax).GenerateN(uint(testCaseLength)) {
			if !(v >= actualMin && v <= actualMax) { return false }
		}
		return true
	}

	checkFailPropery(t, quick.Check(betweenBeingInRange, &globalPropertConf), "Between_Numeric", "range satisfaction")
}

// ------ Only property tests:

func TestOnlyHavingOnlyOneResult(t *testing.T) {
	nameGen := Only("John")
	surnameGen := Only("Doe")
	ageGen := Only(25)

	type Person struct {
		Name 		string
		Surname 	string
		Age 		int
	}

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


