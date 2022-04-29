package tests

import "fmt"

func (assertion *TestSuiteAssertion) AssertBool(name string, valueToCompare bool, compareTo bool) *TestSuiteCaseStepAssertionResult {
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
		assertionResult.Error = "Booleans cannot be nil"
	case ShouldNotBeNil:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot be nil"
	case ShouldBeEmpty:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot be empty"
	case ShouldNotBeEmpty:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot be empty"
	case ShouldContainSubstring:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot contain string"
	case ShouldNotContainSubstring:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot contain string"
	case ShouldBeGreaterThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldNotBeGreaterThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldBeLessThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldNotBeLessThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldBeGreaterOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldNotBeGreaterOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldBeLessOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldNotBeLessOrEqualThan:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have greater operator"
	case ShouldHaveCountOf:
		assertionResult.Passed = false
		assertionResult.Error = "Booleans cannot have count"
	case ShouldBeTrue:
		if valueToCompare {
			assertionResult.Passed = true
			assertionResult.ExpectedValue = "true"
		} else {
			assertionResult.Passed = false
			assertionResult.ExpectedValue = "true"
			assertionResult.Error = fmt.Sprintf("Expected %v to be true", name)
		}
	case ShouldBeFalse:
		if !valueToCompare {
			assertionResult.Passed = true
			assertionResult.ExpectedValue = "false"
		} else {
			assertionResult.Passed = false
			assertionResult.ExpectedValue = "false"
			assertionResult.Error = fmt.Sprintf("Expected %v to be false", name)
		}

	default:
		assertionResult.Passed = false
		assertionResult.Error = "Assertion not implemented"
	}

	return &assertionResult
}
