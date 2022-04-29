package tests

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/helper/reflect_helper"
	"github.com/cjlapao/common-go/helper/strhelper"
)

type TestSuiteAssertion struct {
	Object         interface{}
	Assertion      string
	Operation      TestSuiteAssertionOperation
	FieldTree      []string
	ExpectedResult string
}

func NewTestSuiteAssertion(assertion string) *TestSuiteAssertion {
	result, err := TestSuiteAssertion{}.Parse(assertion)

	if err != nil {
		return nil
	}

	return result
}

func (assert TestSuiteAssertion) Parse(str string) (*TestSuiteAssertion, error) {
	parts := strings.Split(str, " ")
	if len(parts) < 2 {
		return nil, errors.New("not enough operation")
	}

	fieldParts := strings.Split(parts[0], ".")

	if len(fieldParts) < 2 || !strings.HasPrefix(parts[0], "result.") {
		return nil, errors.New("not a result query")
	}

	result := TestSuiteAssertion{
		FieldTree: make([]string, 0),
	}
	result.Assertion = str
	result.Operation = result.Operation.Parse(parts[1])

	// Creating the field tree for traversing during the step check
	for _, fieldPart := range fieldParts {
		if strings.ToLower(fieldPart) != "result" {
			parsedFieldPart := strings.ReplaceAll(fieldPart, "::", ".")
			result.FieldTree = append(result.FieldTree, strings.ToLower(parsedFieldPart))
		}
	}

	if len(parts) >= 3 {
		for i := 2; i < len(parts); i = i + 1 {
			if len(result.ExpectedResult) > 0 {
				result.ExpectedResult += " "
			}
			result.ExpectedResult += parts[i]
		}
	}

	result.ExpectedResult = strings.TrimLeft(result.ExpectedResult, "\"")
	result.ExpectedResult = strings.TrimRight(result.ExpectedResult, "\"")
	result.ExpectedResult = strings.TrimLeft(result.ExpectedResult, "'")
	result.ExpectedResult = strings.TrimRight(result.ExpectedResult, "'")

	return &result, nil
}

func (assertion *TestSuiteAssertion) Assert(objectToAssert interface{}) *TestSuiteCaseStepAssertionResult {
	// Setting the object on first assert
	if reflect_helper.IsNilOrEmpty(assertion.Object) {
		assertion.Object = objectToAssert
	}

	// queueing the next field in the field tree of the assertion
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
		return assertion.AssertMap(assertkey, v)
	default:
		for i := 0; i < v.NumField(); i = i + 1 {
			fv := v.Field(i)
			ft := t.Field(i)
			name := fmt.Sprintf("%v", ft.Name)
			if strings.EqualFold(name, assertkey) {
				switch fv.Kind() {
				case reflect.Bool:
					compareTo := strhelper.ToBoolean(assertion.ExpectedResult)
					return assertion.AssertBool(ft.Name, fv.Bool(), compareTo)
				case reflect.String:
					return assertion.AssertString(ft.Name, fv.String(), assertion.ExpectedResult)
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
					return assertion.AssertInt(ft.Name, fv.Int(), int64(compareTo))
				case reflect.Map:
					return assertion.AssertMap(ft.Name, fv)
				case reflect.Interface:
					return assertion.Assert(fv.Interface())
				}
				break
			}
		}
	}
	return nil
}
