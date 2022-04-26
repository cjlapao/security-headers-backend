package tests

import "github.com/cjlapao/common-go/helper"

type TestSuiteCase struct {
	ID        string
	TestSuite *TestSuite
	Name      string `json:"name" yaml:"name"`
	Steps     []*TestSuiteStep
}

func (testCase *TestSuiteCase) AddStep(stepType string) *TestSuiteStep {
	result := TestSuiteStep{
		ID:            helper.RandomString(20),
		TestSuiteCase: testCase,
		Type:          "http",
		Method:        "GET",
		Weight:        10,
		Assertions:    make([]*TestSuiteAssertion, 0),
	}

	if testCase.TestSuite != nil && testCase.TestSuite.TargetSite != "" {
		result.Url = testCase.TestSuite.TargetSite
	}

	if stepType != "" {
		result.Type = stepType
	}

	testCase.Steps = append(testCase.Steps, &result)

	return &result
}

func (testCase TestSuiteCase) Run() []*TestSuiteCaseStepResult {
	stepResults := make([]*TestSuiteCaseStepResult, 0)

	if testCase.Steps != nil && len(testCase.Steps) > 0 {
		for _, step := range testCase.Steps {
			result := step.Run()
			stepResults = append(stepResults, result)
		}
	}

	return stepResults
}
