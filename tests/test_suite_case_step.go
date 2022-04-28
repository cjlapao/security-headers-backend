package tests

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/helper/strhelper"
)

type TestSuiteStep struct {
	TestSuiteCase    *TestSuiteCase        `json:"" yaml:""`
	ID               string                `json:"" yaml:""`
	Type             string                `json:"type" yaml:"type"`
	Weight           int                   `json:"weight" yaml:"weight"`
	Assertions       []string              `json:"assertions" yaml:"assertions"`
	ParsedAssertions []*TestSuiteAssertion `json:"" yaml:""`
}

func (testSuitStep *TestSuiteStep) Run() *TestSuiteCaseStepResult {
	stepResult := TestSuiteCaseStepResult{
		ID:               testSuitStep.ID,
		Passed:           true,
		Weight:           testSuitStep.Weight,
		AssertionResults: make([]*TestSuiteCaseStepAssertionResult, 0),
	}

	if testSuitStep.ParsedAssertions != nil && len(testSuitStep.ParsedAssertions) > 0 {
		for _, assertion := range testSuitStep.ParsedAssertions {
			assertionResult := testSuitStep.NewAssert(testSuitStep.TestSuiteCase.Object, assertion)
			stepResult.AssertionResults = append(stepResult.AssertionResults, assertionResult)
		}
	}

	for _, assertionResult := range stepResult.AssertionResults {
		if !assertionResult.Passed {
			stepResult.Passed = false
			break
		}
	}

	return &stepResult
}

func (testSuitStep *TestSuiteStep) AddAssertion(assertion string) error {
	if testSuitStep.ParsedAssertions == nil {
		testSuitStep.ParsedAssertions = make([]*TestSuiteAssertion, 0)
	}

	newAssertion := NewTestSuiteAssertion(assertion)

	if newAssertion == nil {
		return errors.New("failed to parse assertion")
	}

	testSuitStep.ParsedAssertions = append(testSuitStep.ParsedAssertions, newAssertion)
	return nil
}

func (testSuitStep *TestSuiteStep) TestHeaders(response *http.Response) error {
	return nil
}

func (testSuitStep *TestSuiteStep) Assert(objectToAssert interface{}, assertion *TestSuiteAssertion, field ...string) *TestSuiteCaseStepAssertionResult {
	assertkey := assertion.Field
	if field != nil && len(field) == 1 {
		assertkey = field[0]
	}
	t := reflect.TypeOf(objectToAssert)
	v := reflect.ValueOf(objectToAssert)
	switch t.Kind() {
	case reflect.Map:
		return testSuitStep.AssertMap(assertion.Field, v, assertion)
	default:
		for i := 0; i < v.NumField(); i = i + 1 {
			fv := v.Field(i)
			ft := t.Field(i)
			name := fmt.Sprintf("%v", ft.Name)
			if strings.EqualFold(name, assertkey) {
				switch fv.Kind() {
				case reflect.Bool:
					compareTo := strhelper.ToBoolean(assertion.ExpectedResult)
					return testSuitStep.AssertBool(ft.Name, fv.Bool(), compareTo, assertion)
				case reflect.String:
					return testSuitStep.AssertString(ft.Name, fv.String(), assertion.ExpectedResult, assertion)
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					compareTo, err := strconv.Atoi(assertion.ExpectedResult)
					if err != nil {
						return &TestSuiteCaseStepAssertionResult{
							Passed:        false,
							Assertion:     assertion.Assertion,
							ExpectedValue: assertion.ExpectedResult,
							Error:         err.Error(),
						}
					}
					return testSuitStep.AssertInt(ft.Name, fv.Int(), int64(compareTo), assertion)
				case reflect.Map:
					return testSuitStep.AssertMap(ft.Name, fv, assertion)
				case reflect.Interface:
					return testSuitStep.Assert(fv.Interface(), assertion, assertion.Property)
				}
				break
			}
		}
	}
	return nil
}

func (testSuitStep *TestSuiteStep) NewAssert(objectToAssert interface{}, assertion *TestSuiteAssertion) *TestSuiteCaseStepAssertionResult {
	assertkey := ""
	if len(assertion.FieldTree) > 0 {
		assertkey = assertion.FieldTree[0]
		// Popping out the initial field
		if len(assertion.FieldTree) > 1 {
			assertion.FieldTree = assertion.FieldTree[1:]
		}
	}

	if assertkey == "" {
		return nil
	}

	t := reflect.TypeOf(objectToAssert)
	v := reflect.ValueOf(objectToAssert)
	switch t.Kind() {
	case reflect.Map:
		return testSuitStep.AssertMap(assertion.Field, v, assertion)
	default:
		for i := 0; i < v.NumField(); i = i + 1 {
			fv := v.Field(i)
			ft := t.Field(i)
			name := fmt.Sprintf("%v", ft.Name)
			if strings.EqualFold(name, assertkey) {
				switch fv.Kind() {
				case reflect.Bool:
					compareTo := strhelper.ToBoolean(assertion.ExpectedResult)
					return testSuitStep.AssertBool(ft.Name, fv.Bool(), compareTo, assertion)
				case reflect.String:
					return testSuitStep.AssertString(ft.Name, fv.String(), assertion.ExpectedResult, assertion)
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					compareTo, err := strconv.Atoi(assertion.ExpectedResult)
					if err != nil {
						return &TestSuiteCaseStepAssertionResult{
							Passed:        false,
							Assertion:     assertion.Assertion,
							ExpectedValue: assertion.ExpectedResult,
							Error:         err.Error(),
						}
					}
					return testSuitStep.AssertInt(ft.Name, fv.Int(), int64(compareTo), assertion)
				case reflect.Map:
					return testSuitStep.AssertMap(ft.Name, fv, assertion)
				case reflect.Interface:
					return testSuitStep.NewAssert(fv.Interface(), assertion)
				}
				break
			}
		}
	}
	return nil
}
func (testSuitStep *TestSuiteStep) AssertInt(name string, valueToCompare int64, compareTo int64, assertion *TestSuiteAssertion) *TestSuiteCaseStepAssertionResult {
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

func (testSuitStep *TestSuiteStep) AssertString(name string, valueToCompare string, compareTo string, assertion *TestSuiteAssertion) *TestSuiteCaseStepAssertionResult {
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

func (testSuitStep *TestSuiteStep) AssertBool(name string, valueToCompare bool, compareTo bool, assertion *TestSuiteAssertion) *TestSuiteCaseStepAssertionResult {
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

func (testSuitStep *TestSuiteStep) AssertMap(name string, mapValue reflect.Value, assertion *TestSuiteAssertion) *TestSuiteCaseStepAssertionResult {
	assertkey := ""
	if len(assertion.FieldTree) > 0 {
		assertkey = assertion.FieldTree[0]
		// Popping out the initial field
		if len(assertion.FieldTree) > 1 {
			assertion.FieldTree = assertion.FieldTree[1:]
		}
	}

	if assertkey == "" {
		return nil
	}

	found := false

	for _, e := range mapValue.MapKeys() {
		name := fmt.Sprintf("%v", e)
		if strings.EqualFold(name, assertkey) {
			found = true
			val := mapValue.MapIndex(e)
			switch valType := val.Interface().(type) {
			case bool:
				compareTo := strhelper.ToBoolean(assertion.ExpectedResult)
				return testSuitStep.AssertBool(name, valType, compareTo, assertion)
			case string:
				return testSuitStep.AssertString(name, valType, assertion.ExpectedResult, assertion)
			case []string:
				value := ""
				for _, str := range valType {
					value = fmt.Sprintf("%v;%v", value, str)
				}
				value = strings.TrimLeft(value, ";")
				return testSuitStep.AssertString(name, value, assertion.ExpectedResult, assertion)
			case int:
				compareTo, err := strconv.Atoi(assertion.ExpectedResult)
				if err != nil {
					return &TestSuiteCaseStepAssertionResult{
						Passed:        false,
						Assertion:     assertion.Assertion,
						ExpectedValue: assertion.ExpectedResult,
						Error:         err.Error(),
					}
				}

				return testSuitStep.AssertInt(name, int64(valType), int64(compareTo), assertion)
			case interface{}:
				return testSuitStep.NewAssert(val.Interface(), assertion)
			default:
				isMap := reflect.ValueOf(valType).Kind() == reflect.Map
				if isMap {
					return testSuitStep.AssertMap(name, val, assertion)
				}
				fmt.Printf("%v", valType)
			}
		}
	}

	if !found {
		return &TestSuiteCaseStepAssertionResult{
			Passed:        false,
			Assertion:     assertion.Assertion,
			ExpectedValue: assertion.ExpectedResult,
			Error:         fmt.Sprintf("%v was not found", assertion.Field),
		}
	}

	return nil
}
