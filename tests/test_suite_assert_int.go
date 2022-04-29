package tests

import "fmt"

func (assertion *TestSuiteAssertion) AssertInt(name string, valueToCompare int64, compareTo int64) *TestSuiteCaseStepAssertionResult {
	assertionResult := TestSuiteCaseStepAssertionResult{
		Assertion:     assertion.Assertion,
		Passed:        false,
		ExpectedValue: fmt.Sprintf("%v", compareTo),
		FoundValue:    fmt.Sprintf("%v", valueToCompare),
	}

	switch assertion.Operation {
	case ShouldBeEqual:
		if valueToCompare == compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldNotBeEqual:
		if valueToCompare != compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldBeNil:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be nil"
	case ShouldNotBeNil:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be nil"
	case ShouldBeEmpty:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be nil"
	case ShouldNotBeEmpty:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be nil"
	case ShouldContainSubstring:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be nil"
	case ShouldNotContainSubstring:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be nil"
	case ShouldBeGreaterThan:
		if valueToCompare > compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be greater than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldNotBeGreaterThan:
		if valueToCompare < compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be greater than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldBeLessThan:
		if valueToCompare < compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be less than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldNotBeLessThan:
		if valueToCompare > compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be less than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldBeGreaterOrEqualThan:
		if valueToCompare >= compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be greater or equal than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldNotBeGreaterOrEqualThan:
		if valueToCompare <= compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be greater or equal than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldBeLessOrEqualThan:
		if valueToCompare <= compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v to be less or equal than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldNotBeLessOrEqualThan:
		if valueToCompare >= compareTo {
			assertionResult.Passed = true
		} else {
			assertionResult.Passed = false
			assertionResult.Error = fmt.Sprintf("Expected %v not to be less or equal than %v, found %v", name, compareTo, valueToCompare)
		}
	case ShouldHaveCountOf:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot have count"
	case ShouldBeTrue:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be booleans"
	case ShouldBeFalse:
		assertionResult.Passed = false
		assertionResult.Error = "Integers cannot be booleans"
	default:
		assertionResult.Passed = false
		assertionResult.Error = "Assertion not implemented"
	}

	return &assertionResult
}
