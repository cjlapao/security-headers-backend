package tests

import (
	"fmt"
	"strconv"
	"strings"
)

func (assertion *TestSuiteAssertion) AssertString(name string, valueToCompare string, compareTo string) *TestSuiteCaseStepAssertionResult {
	assertionResult := TestSuiteCaseStepAssertionResult{
		Assertion:     assertion.Assertion,
		Passed:        false,
		ExpectedValue: fmt.Sprintf("%v", compareTo),
		FoundValue:    fmt.Sprintf("%v", valueToCompare),
	}

	switch assertion.Operation {
	case ShouldBeEqual:
		if strings.EqualFold(compareTo, valueToCompare) {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldNotBeEqual:
		if !strings.EqualFold(compareTo, valueToCompare) {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldBeNil:
		if valueToCompare == "" {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be empty, found %v", name, valueToCompare)
		}
	case ShouldNotBeNil:
		if valueToCompare != "" {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be empty", name)
		}
	case ShouldBeEmpty:
		if valueToCompare == "" {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be empty, found %v", name, valueToCompare)
		}
	case ShouldNotBeEmpty:
		if valueToCompare != "" {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be empty", name)
		}
	case ShouldContainSubstring:
		if strings.ContainsAny(compareTo, valueToCompare) {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to contain %v", name, compareTo)
		}
	case ShouldNotContainSubstring:
		if !strings.ContainsAny(compareTo, valueToCompare) {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to contain %v", name, compareTo)
		}
	case ShouldBeGreaterThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldNotBeGreaterThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldBeLessThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldNotBeLessThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldBeGreaterOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldNotBeGreaterOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldBeLessOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldNotBeLessOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot have greater operator"
	case ShouldHaveCountOf:
		valueLen, err := strconv.Atoi(compareTo)
		if err != nil {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v count to be a number, found %v", name, compareTo)
		}
		if len(valueToCompare) == valueLen {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to have count of %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldBeTrue:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot be booleans"
	case ShouldBeFalse:
		assertionResult.Passed = false
		assertionResult.Error = "Strings cannot be booleans"
	default:
		assertionResult.Passed = false
		assertionResult.Error = "Assertion not implemented"
	}

	return &assertionResult
}
