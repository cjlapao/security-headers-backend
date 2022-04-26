package tests

import (
	"errors"
	"strings"
)

type TestSuiteAssertion struct {
	Assertion      string
	Field          string
	Property       string
	Operation      TestSuiteAssertionOperation
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
	if len(parts) < 2 || len(parts) > 3 {
		return nil, errors.New("not enough operation")
	}

	fieldParts := strings.Split(parts[0], ".")

	if len(fieldParts) < 2 || !strings.HasPrefix(parts[0], "result.") {
		return nil, errors.New("not a result query")
	}

	result := TestSuiteAssertion{}
	result.Assertion = str
	result.Operation = result.Operation.Parse(parts[1])

	if len(fieldParts) >= 2 {
		result.Field = fieldParts[1]
	}

	if len(fieldParts) >= 3 {
		result.Property = fieldParts[2]
	}

	if len(parts) == 3 {
		result.ExpectedResult = parts[2]
	}

	return &result, nil
}
