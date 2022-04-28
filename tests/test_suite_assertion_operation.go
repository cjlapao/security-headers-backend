package tests

import "strings"

type TestSuiteAssertionOperation int

const (
	ShouldBeEqual TestSuiteAssertionOperation = iota
	ShouldNotBeEqual
	ShouldBeNil
	ShouldNotBeNil
	ShouldBeEmpty
	ShouldNotBeEmpty
	ShouldContainSubstring
	ShouldNotContainSubstring
	ShouldBeGreaterThan
	ShouldNotBeGreaterThan
	ShouldBeLessThan
	ShouldNotBeLessThan
	ShouldBeGreaterOrEqualThan
	ShouldNotBeGreaterOrEqualThan
	ShouldBeLessOrEqualThan
	ShouldNotBeLessOrEqualThan
	ShouldHaveCountOf
	ShouldBeTrue
	ShouldBeFalse
)

func (op TestSuiteAssertionOperation) Parse(str string) TestSuiteAssertionOperation {
	var result TestSuiteAssertionOperation
	switch strings.ToLower(str) {
	case "shouldbeequal":
		result = ShouldBeEqual
	case "shouldnotbeequal":
		result = ShouldNotBeEqual
	case "shouldbenil":
		result = ShouldBeNil
	case "shouldnotbenil":
		result = ShouldNotBeNil
	case "shouldbeempty":
		result = ShouldBeEmpty
	case "shouldnotbeempty":
		result = ShouldNotBeEmpty
	case "shouldcontainsubstring":
		result = ShouldContainSubstring
	case "shouldnotcontainsubstring":
		result = ShouldNotContainSubstring
	case "shouldbegreaterthan":
		result = ShouldBeGreaterThan
	case "shouldnotbegreaterthan":
		result = ShouldNotBeGreaterThan
	case "shouldbelessthan":
		result = ShouldBeLessThan
	case "shouldbegreaterorequalthan":
		result = ShouldBeGreaterOrEqualThan
	case "shouldnotbegreaterorequalthan":
		result = ShouldNotBeGreaterOrEqualThan
	case "shouldbelessorequalthan":
		result = ShouldBeLessOrEqualThan
	case "shouldnotbelessorequalthan":
		result = ShouldNotBeLessOrEqualThan
	case "shouldhavecountof":
		result = ShouldHaveCountOf
	case "shouldbetrue":
		result = ShouldBeTrue
	case "shouldbefalse":
		result = ShouldBeFalse
	}

	return result
}
