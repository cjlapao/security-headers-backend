package tests

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TestSuiteStep struct {
	TestSuiteCase *TestSuiteCase
	ID            string
	Type          string
	Method        string
	Url           string
	Timeout       string
	Weight        int
	Assertions    []*TestSuiteAssertion
}

func (testSuitStep *TestSuiteStep) Run() *TestSuiteCaseStepResult {
	stepResult := TestSuiteCaseStepResult{
		ID:               testSuitStep.ID,
		Passed:           true,
		Weight:           testSuitStep.Weight,
		AssertionResults: make([]*TestSuiteCaseStepAssertionResult, 0),
	}

	if testSuitStep.Assertions != nil && len(testSuitStep.Assertions) > 0 {
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

		for _, assertion := range testSuitStep.Assertions {
			switch strings.ToLower(assertion.Field) {
			case "statuscode":
				result, err := testSuitStep.TestStatusCode(response, assertion)
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
	if testSuitStep.Assertions == nil {
		testSuitStep.Assertions = make([]*TestSuiteAssertion, 0)
	}

	newAssertion := NewTestSuiteAssertion(assertion)

	if newAssertion == nil {
		return errors.New("failed to parse assertion")
	}

	testSuitStep.Assertions = append(testSuitStep.Assertions, newAssertion)
	return nil
}

func (testSuitStep *TestSuiteStep) TestHeaders(response *http.Response) error {
	return nil
}

func (testSuitStep *TestSuiteStep) TestStatusCode(response *http.Response, assertion *TestSuiteAssertion) (bool, error) {
	switch assertion.Operation {
	case ShouldBeEqual:
		if strings.EqualFold(fmt.Sprint(response.StatusCode), assertion.ExpectedResult) {
			return true, nil
		}
		return false, errors.New(fmt.Sprintf("expected status code of %v, found %v", response.StatusCode, assertion.ExpectedResult))
	default:
		return false, errors.New("not implemented")
	}
}
