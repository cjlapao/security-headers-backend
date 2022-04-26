package controllers

import (
	"net/http"

	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/security-headers-backend/common"
	"github.com/cjlapao/security-headers-backend/entities"
	"github.com/cjlapao/security-headers-backend/tests"
)

func ValidateHeadersController(w http.ResponseWriter, r *http.Request) {
	var body entities.SecurityHeaderRequest

	err := http_helper.MapRequestBody(r, &body)

	if err != nil {
		common.WriteError(w, err)
		return
	}

	suites := tests.TestSuite{
		TargetSite: "https://devprod-sfc.ivanticlouddev.com/api/disco/probe",
	}

	testCase := suites.AddTestCase("x1")
	testCaseStep := testCase.AddStep("http")
	testCaseStep.Weight = 20
	testCaseStep.AddAssertion("result.statuscode ShouldBeEqual 200")

	testCase1 := suites.AddTestCase("x2")
	testCaseStep1 := testCase1.AddStep("http")
	testCaseStep1.Url = "https://devprod-sfc.ivanticlouddev.com/api/disco/probe123"
	testCaseStep1.Weight = 5
	testCaseStep1.AddAssertion("result.statuscode ShouldBeEqual 200")

	suites.Test()
}
