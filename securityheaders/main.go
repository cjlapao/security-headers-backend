package securityheaders

import (
	"github.com/cjlapao/common-go/guard"
	"github.com/cjlapao/security-headers-backend/tests"
)

type SecurityHeadersService struct {
	TargetSite string
	TestSuite  *tests.TestSuite
	Result     *SecurityHeadersResult
}

func New(targetSite string) *SecurityHeadersService {
	if err := guard.EmptyOrNil(targetSite, "Target Site"); err != nil {
		return nil
	}

	svc := SecurityHeadersService{
		TargetSite: targetSite,
	}

	return &svc
}

func (s *SecurityHeadersService) Calculate() *SecurityHeadersResult {
	s.TestSuite = tests.NewTestSuite()
	s.TestSuite.TargetSite = s.TargetSite
	testCase := s.TestSuite.AddTestCase("strict-transport-security")
	cookiesTestCase := s.TestSuite.AddTestCase("Cookies")
	cookiesStep := cookiesTestCase.AddStep()
	cookiesStep.AddAssertion("result.cookies.Identity::External.secure ShouldBeTrue")

	// s.addStrictTransportSecurityStep(testCase)
	s.addContent(testCase)

	s.TestSuite.Test()

	result := SecurityHeadersResult{
		Points:     s.TestSuite.Result.TotalPassedWeight,
		Score:      s.TestSuite.Result.Score,
		TargetSite: s.TargetSite,
		Headers:    make([]SecurityHeadersResultHeader, 0),
		Cookies:    make([]SecurityHeadersResultCookie, 0),
	}

	s.Result = &result
	return s.Result
}

func (s *SecurityHeadersService) addStrictTransportSecurityStep(testCase *tests.TestSuiteCase) {
	testCaseStep := testCase.AddStep()
	testCaseStep.Weight = 20
	testCaseStep.AddAssertion("result.statuscode ShouldBeEqual 200")
	testCaseStep.AddAssertion("result.headers.strict-transport-security ShouldNotBeEmpty")
	testCaseStep.AddAssertion("result.headers.strict-transport-security ShouldContainSubstring max-age=")
	testCaseStep.AddAssertion("result.headers.strict-transport-security ShouldContainSubstring  includeSubDomains")
}

func (s *SecurityHeadersService) addContent(testCase *tests.TestSuiteCase) {
	testCaseStep := testCase.AddStep()
	testCaseStep.Weight = 20
	testCaseStep.AddAssertion("result.statuscode ShouldBeEqual 200")
	testCaseStep.AddAssertion("result.headers.content-type ShouldNotBeEmpty")
}
