package tests

import "strings"

type TestSuiteAssertionOperation int

const (
	ShouldBeEqual TestSuiteAssertionOperation = iota
	ShouldNotBeEqual
	ShouldBeNil
	ShouldNotBeNil
	ShouldContainSubstring
	ShouldNotContainSubstring
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
	case "shouldcontainsubstring":
		result = ShouldContainSubstring
	case "shouldnotcontainsubstring":
		result = ShouldNotContainSubstring
	}

	return result
}
