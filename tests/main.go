package tests

import "github.com/cjlapao/common-go/log"

var logger = log.Get()

type TestSuiteCaseStepResult struct {
	ID               string
	Passed           bool
	Weight           int
	AssertionResults []*TestSuiteCaseStepAssertionResult
}

type TestSuiteCaseStepAssertionResult struct {
	Passed        bool
	Assertion     string
	ExpectedValue string
	FoundValue    string
	Error         string
}
