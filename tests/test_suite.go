package tests

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cjlapao/common-go/helper"
	"gopkg.in/yaml.v3"
)

type TestSuite struct {
	TargetSite string           `json:"targetSite" yaml:"targetSite"`
	TestCases  []*TestSuiteCase `json:"testCases" yaml:"testCases"`
	Result     *TestSuiteResult `json:"result" yaml:"result"`
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

func LoadYamlFromFile(filePath string, dest *TestSuite) error {
	if !helper.FileExists(filePath) {
		return errors.New(fmt.Sprintf("file %v was not found", filePath))
	}

	fileContent, err := helper.ReadFromFile(filePath)
	if err != nil {
		return err
	}

	UnmarshalYaml(string(fileContent), dest)
	return nil
}

func UnmarshalJson(jsonStr string, dest *TestSuite) {
	var tmpObj TestSuite
	json.Unmarshal([]byte(jsonStr), &tmpObj)

	getObject(tmpObj, dest)
}

func UnmarshalYaml(yamlStr string, dest *TestSuite) {
	var tmpObj TestSuite
	yaml.Unmarshal([]byte(yamlStr), &tmpObj)

	getObject(tmpObj, dest)
}

func getObject(tempObject TestSuite, dest *TestSuite) {
	dest.TargetSite = tempObject.TargetSite
	for _, tmpCase := range tempObject.TestCases {
		testCase := dest.AddTestCase(tmpCase.Name)
		for _, tmpStep := range tmpCase.Steps {
			step := testCase.AddStep(tmpStep.Type)
			if tmpStep.Method != "" {
				step.Method = tmpStep.Method
			}

			if tmpStep.Timeout != "" {
				step.Timeout = tmpStep.Timeout
			}

			if tmpStep.Url != "" {
				step.Url = tmpStep.Url
			}

			if tmpStep.Weight > 0 {
				step.Weight = tmpStep.Weight
			}

			for _, tmpAssertion := range tmpStep.Assertions {
				step.AddAssertion(tmpAssertion)
			}
		}
	}
}
