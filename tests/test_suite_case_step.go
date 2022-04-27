package tests

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/helper/strhelper"
)

type TestSuiteStep struct {
	TestSuiteCase    *TestSuiteCase        `json:"" yaml:""`
	ID               string                `json:"" yaml:""`
	Type             string                `json:"type" yaml:"type"`
	Method           string                `json:"method" yaml:"method"`
	Url              string                `json:"url" yaml:"url"`
	Timeout          string                `json:"timeout" yaml:"timeout"`
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
		url, err := url.Parse(testSuitStep.Url)

		if err != nil {
			stepResult.Passed = false
			stepResult.AssertionResults = append(stepResult.AssertionResults, &TestSuiteCaseStepAssertionResult{
				Error: err.Error(),
			})
			return &stepResult
		}

		client := http.Client{}

		request := http.Request{
			Method: strings.ToUpper(testSuitStep.Method),
			URL:    url,
		}

		response, err := client.Do(&request)

		if err != nil {
			stepResult.Passed = false
			stepResult.AssertionResults = append(stepResult.AssertionResults, &TestSuiteCaseStepAssertionResult{
				Error: err.Error(),
			})
			return &stepResult
		}

		for _, assertion := range testSuitStep.ParsedAssertions {
			// we may use headers in the test case so converting to the right field name
			if strings.EqualFold(assertion.Field, "headers") {
				assertion.Field = "header"
			}
			result, err := testSuitStep.Assert(response, assertion)
			assertionResult := TestSuiteCaseStepAssertionResult{
				Passed:    result,
				Assertion: assertion.Assertion,
			}

			if err != nil {
				assertionResult.Error = err.Error()
			}

			stepResult.AssertionResults = append(stepResult.AssertionResults, &assertionResult)
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

func (testSuitStep *TestSuiteStep) Assert(response *http.Response, assertion *TestSuiteAssertion) (bool, error) {
	t := reflect.TypeOf(*response)
	v := reflect.ValueOf(*response)
	for i := 0; i < v.NumField(); i = i + 1 {
		fv := v.Field(i)
		ft := t.Field(i)
		if strings.EqualFold(ft.Name, assertion.Field) {
			switch fv.Kind() {
			case reflect.Bool:
				compareTo := strhelper.ToBoolean(assertion.ExpectedResult)
				return testSuitStep.AssertBool(ft.Name, fv.Bool(), compareTo, assertion.Operation)
			case reflect.String:
				return testSuitStep.AssertString(ft.Name, fv.String(), assertion.ExpectedResult, assertion.Operation)
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				compareTo, err := strconv.Atoi(assertion.ExpectedResult)
				if err != nil {
					return false, err
				}
				return testSuitStep.AssertInt(ft.Name, fv.Int(), int64(compareTo), assertion.Operation)
			case reflect.Map:
				if assertion.Property == "" {
					return false, errors.New("field is a map and no property to compare to was defined")
				}
				for _, e := range fv.MapKeys() {
					name := fmt.Sprintf("%v", e)
					if strings.EqualFold(name, assertion.Property) {
						val := fv.MapIndex(e)
						switch valType := val.Interface().(type) {
						case bool:
							compareTo := strhelper.ToBoolean(assertion.ExpectedResult)
							return testSuitStep.AssertBool(ft.Name, valType, compareTo, assertion.Operation)
						case string:
							return testSuitStep.AssertString(ft.Name, valType, assertion.ExpectedResult, assertion.Operation)
						case []string:
							value := ""
							for _, str := range valType {
								value = fmt.Sprintf("%v;%v", value, str)
							}
							value = strings.TrimLeft(value, ";")
							return testSuitStep.AssertString(ft.Name, value, assertion.ExpectedResult, assertion.Operation)
						case int:
							compareTo, err := strconv.Atoi(assertion.ExpectedResult)
							if err != nil {
								return false, err
							}

							return testSuitStep.AssertInt(ft.Name, int64(valType), int64(compareTo), assertion.Operation)
						}
					}
				}
			}
			break
		}
	}
	return true, nil
}

func (testSuitStep *TestSuiteStep) AssertInt(name string, valueToCompare int64, compareTo int64, operation TestSuiteAssertionOperation) (bool, error) {
	switch operation {
	case ShouldBeEqual:
		if valueToCompare == compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeEqual:
		if valueToCompare != compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be %v, found %v", name, compareTo, valueToCompare))
	case ShouldBeNil:
		return false, errors.New("integers cannot be nil")
	case ShouldNotBeNil:
		return false, errors.New("integers cannot be nil")
	case ShouldBeEmpty:
		return false, errors.New("integers cannot be empty")
	case ShouldNotBeEmpty:
		return false, errors.New("integers cannot be empty")
	case ShouldContainSubstring:
		return false, errors.New("integers cannot contain strings")
	case ShouldNotContainSubstring:
		return false, errors.New("integers cannot contain strings")
	case ShouldBeGreaterThan:
		if valueToCompare > compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be greater than %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeGreaterThan:
		if valueToCompare < compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be greater than %v, found %v", name, compareTo, valueToCompare))
	case ShouldBeLessThan:
		if valueToCompare < compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be less than %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeLessThan:
		if valueToCompare > compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be less than %v, found %v", name, compareTo, valueToCompare))
	case ShouldBeGreaterOrEqualThan:
		if valueToCompare >= compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be greater or equal than %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeGreaterOrEqualThan:
		if valueToCompare <= compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be greater or equal than %v, found %v", name, compareTo, valueToCompare))
	case ShouldBeLessOrEqualThan:
		if valueToCompare <= compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be less or equal than %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeLessOrEqualThan:
		if valueToCompare >= compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be less or equal than %v, found %v", name, compareTo, valueToCompare))
	case ShouldHaveCountOf:
		return false, errors.New("integers cannot have count")
	default:
		return false, errors.New("not implemented")
	}
}

func (testSuitStep *TestSuiteStep) AssertString(name string, valueToCompare string, compareTo string, operation TestSuiteAssertionOperation) (bool, error) {
	switch operation {
	case ShouldBeEqual:
		if strings.EqualFold(compareTo, valueToCompare) {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeEqual:
		if !strings.EqualFold(compareTo, valueToCompare) {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be %v, found %v", name, compareTo, valueToCompare))
	case ShouldBeNil:
		if valueToCompare == "" {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be empty, found %v", name, valueToCompare))
	case ShouldNotBeNil:
		if valueToCompare != "" {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be empty", name))
	case ShouldBeEmpty:
		if valueToCompare == "" {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be empty, found %v", name, valueToCompare))
	case ShouldNotBeEmpty:
		if valueToCompare != "" {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be empty", name))
	case ShouldContainSubstring:
		if strings.ContainsAny(compareTo, valueToCompare) {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to contain %v", name, compareTo))
	case ShouldNotContainSubstring:
		if !strings.ContainsAny(compareTo, valueToCompare) {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to contain %v", name, compareTo))
	case ShouldBeGreaterThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldNotBeGreaterThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldBeLessThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldNotBeLessThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldBeGreaterOrEqualThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldNotBeGreaterOrEqualThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldBeLessOrEqualThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldNotBeLessOrEqualThan:
		return false, errors.New("strings cannot have greater operator")
	case ShouldHaveCountOf:
		valueLen, err := strconv.Atoi(compareTo)
		if err != nil {
			return false, errors.New(fmt.Sprintf("expected %v count to be a number, found %v", name, compareTo))
		}
		if len(valueToCompare) == valueLen {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to have count of %v, found %v", name, compareTo, valueToCompare))
	default:
		return false, errors.New("not implemented")
	}
}

func (testSuitStep *TestSuiteStep) AssertBool(name string, valueToCompare bool, compareTo bool, operation TestSuiteAssertionOperation) (bool, error) {
	switch operation {
	case ShouldBeEqual:
		if valueToCompare == compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v to be %v, found %v", name, compareTo, valueToCompare))
	case ShouldNotBeEqual:
		if valueToCompare != compareTo {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected %v not to be %v, found %v", name, compareTo, valueToCompare))
	case ShouldBeNil:
		return false, errors.New("booleans cannot be nil")
	case ShouldNotBeNil:
		return false, errors.New("booleans cannot be nil")
	case ShouldBeEmpty:
		return false, errors.New("booleans cannot be empty")
	case ShouldNotBeEmpty:
		return false, errors.New("booleans cannot be empty")
	case ShouldContainSubstring:
		return false, errors.New("booleans cannot contain string")
	case ShouldNotContainSubstring:
		return false, errors.New("booleans cannot contain string")
	case ShouldBeGreaterThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldNotBeGreaterThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldBeLessThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldNotBeLessThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldBeGreaterOrEqualThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldNotBeGreaterOrEqualThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldBeLessOrEqualThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldNotBeLessOrEqualThan:
		return false, errors.New("booleans cannot have greater operator")
	case ShouldHaveCountOf:
		return false, errors.New("booleans cannot have count")
	default:
		return false, errors.New("not implemented")
	}
}
