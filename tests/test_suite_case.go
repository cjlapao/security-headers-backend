package tests

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cjlapao/common-go/helper"
)

type TestSuiteCase struct {
	ID           string                `json:"" yaml:""`
	TestSuite    *TestSuite            `json:"" yaml:""`
	Name         string                `json:"name" yaml:"name"`
	Type         string                `json:"type" yaml:"type"`
	Method       string                `json:"method" yaml:"method"`
	Url          string                `json:"url" yaml:"url"`
	Timeout      int                   `json:"timeout" yaml:"timeout"`
	Steps        []*TestSuiteStep      `json:"steps" yaml:"steps"`
	HttpResponse TestSuiteHttpResponse `json:"httpResponse" yaml:"httpResponse"`
	Response     *http.Response        `json:"" yaml:""`
	Object       interface{}           `json:"" yaml:""`
}

func (testCase *TestSuiteCase) AddStep() *TestSuiteStep {
	result := TestSuiteStep{
		ID:               helper.RandomString(20),
		TestSuiteCase:    testCase,
		Type:             testCase.Type,
		Weight:           10,
		ParsedAssertions: make([]*TestSuiteAssertion, 0),
	}

	testCase.Steps = append(testCase.Steps, &result)

	return &result
}

func (testCase *TestSuiteCase) WithObject(object interface{}) *TestSuiteCase {
	testCase.Object = object
	testCase.Type = "object"

	return testCase
}

func (testCase *TestSuiteCase) WithUrl(url string) *TestSuiteCase {
	testCase.Type = "http"
	testCase.Url = url

	return testCase
}
func (testCase *TestSuiteCase) Run() []*TestSuiteCaseStepResult {
	stepResults := make([]*TestSuiteCaseStepResult, 0)

	if testCase.Steps != nil && len(testCase.Steps) > 0 {
		if testCase.Type == "http" {
			// Parsing the url
			url, err := url.Parse(testCase.Url)

			if err != nil {
				stepResult := TestSuiteCaseStepResult{
					Passed: false,
					Weight: -1,
				}

				stepResult.AssertionResults = append(stepResult.AssertionResults, &TestSuiteCaseStepAssertionResult{
					Error: err.Error(),
				})

				stepResults = append(stepResults, &stepResult)

				return stepResults
			}

			client := http.Client{}

			request := http.Request{
				Method: strings.ToUpper(testCase.Method),
				URL:    url,
			}

			testCase.Response, err = client.Do(&request)
			if err != nil {
				stepResult := TestSuiteCaseStepResult{
					Passed: false,
					Weight: -1,
				}
				stepResult.AssertionResults = append(stepResult.AssertionResults, &TestSuiteCaseStepAssertionResult{
					Error: err.Error(),
				})
				stepResults = append(stepResults, &stepResult)

				return stepResults
			}

			testCase.HttpResponse = *testCase.HttpResponse.Parse(testCase.Response)
			testCase.Object = testCase.HttpResponse
		}

		for _, step := range testCase.Steps {
			result := step.Run()
			stepResults = append(stepResults, result)
		}
	}

	return stepResults
}
