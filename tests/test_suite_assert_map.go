package tests

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/helper/strhelper"
)

func (assertion *TestSuiteAssertion) AssertMap(name string, mapValue reflect.Value) *TestSuiteCaseStepAssertionResult {
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
				return assertion.AssertBool(name, valType, compareTo)
			case string:
				return assertion.AssertString(name, valType, assertion.ExpectedResult)
			case []string:
				value := ""
				for _, str := range valType {
					value = fmt.Sprintf("%v;%v", value, str)
				}
				value = strings.TrimLeft(value, ";")
				return assertion.AssertString(name, value, assertion.ExpectedResult)
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

				return assertion.AssertInt(name, int64(valType), int64(compareTo))
			case interface{}:
				return assertion.Assert(val.Interface())
			default:
				isMap := reflect.ValueOf(valType).Kind() == reflect.Map
				if isMap {
					return assertion.AssertMap(name, val)
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
			Error:         fmt.Sprintf("%v was not found", assertkey),
		}
	}

	return nil
}
