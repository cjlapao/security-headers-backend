package tests

type TestSuiteResult struct {
	PassedAll         bool
	TotalTests        int
	TotalFailedWeight int
	TotalPassedWeight int
	TotalWeight       int
	Score             string
	Results           []*TestSuiteCaseStepResult
}

func (a *TestSuiteResult) Process() {
	a.PassedAll = true
	a.TotalTests = 0
	a.TotalFailedWeight = 0
	a.TotalPassedWeight = 0
	for _, stepsResult := range a.Results {
		a.TotalTests += 1
		a.TotalWeight += stepsResult.Weight
		if !stepsResult.Passed {
			a.PassedAll = false
			a.TotalFailedWeight += stepsResult.Weight
		} else {
			a.TotalPassedWeight += stepsResult.Weight
		}
	}

	if a.TotalFailedWeight == 0 && a.TotalPassedWeight == 0 {
		a.Score = "A+"
		return
	}

	if a.TotalFailedWeight < 0 || a.TotalPassedWeight < 0 {
		a.Score = "R"
		return
	}

	switch score := a.calculatePercentage(); {
	case score >= 100:
		a.Score = "A+"
	case score >= 90 && score < 100:
		a.Score = "A"
	case score >= 80 && score < 90:
		a.Score = "B"
	case score >= 70 && score < 80:
		a.Score = "C"
	case score >= 60 && score < 70:
		a.Score = "D"
	case score >= 50 && score < 60:
		a.Score = "E"
	case score >= 0 && score < 50:
		a.Score = "F"
	default:
		a.Score = "R"
	}
}

func (r *TestSuiteResult) calculatePercentage() int {
	if r.TotalWeight == 0 {
		return 100
	}

	result := r.TotalPassedWeight * 100 / r.TotalWeight
	return result
}
