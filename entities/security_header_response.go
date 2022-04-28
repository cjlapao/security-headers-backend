package entities

type SecurityHeaderResponse struct {
	TargetSite string                                `json:"targetSite" yaml:"targetSite"`
	Score      string                                `json:"score" yaml:"score"`
	Headers    []SecurityHeaderServiceResponseHeader `json:"headers" yaml:"headers"`
}

type SecurityHeaderServiceResponseHeader struct {
	Passed bool   `json:"passed" yaml:"passed"`
	Header string `json:"header" yaml:"header"`
}
