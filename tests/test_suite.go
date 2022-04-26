package tests

import "github.com/cjlapao/common-go/helper"

type TestSuite struct {
	TargetSite string
	TestCases  []*TestSuiteCase
	Result     *TestSuiteResult
}

func NewTestSuite() *TestSuite {
	result := TestSuite{
		TestCases: make([]*TestSuiteCase, 0),
	}

	return &result
}

func (ts *TestSuite) AddTestCase(name string) *TestSuiteCase {
	result := TestSuiteCase{
		ID:        helper.RandomString(20),
		TestSuite: ts,
		Name:      name,
		Steps:     make([]*TestSuiteStep, 0),
	}

	ts.TestCases = append(ts.TestCases, &result)
	ts.Result = &TestSuiteResult{}
	ts.Result.Results = make([]*TestSuiteCaseStepResult, 0)

	return &result
}

func (ts *TestSuite) Test() {
	if ts.TestCases != nil && len(ts.TestCases) > 0 {
		for _, testCase := range ts.TestCases {
			results := testCase.Run()
			ts.Result.Results = append(ts.Result.Results, results...)
		}
	}

	ts.Result.Process()
}
