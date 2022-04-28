package securityheaders

import "time"

type SecurityHeadersResult struct {
	TargetSite string                        `json:"targetSite" yaml:"targetSite"`
	Points     int                           `json:"points" yaml:"points"`
	Score      string                        `json:"score" yaml:"score"`
	Headers    []SecurityHeadersResultHeader `json:"headers" yaml:"headers"`
	Cookies    []SecurityHeadersResultCookie `json:"cookies" yaml:"cookies"`
}

type SecurityHeadersResultHeader struct {
	Passed bool   `json:"passed" yaml:"passed"`
	Header string `json:"header" yaml:"header"`
	Value  string `json:"value" yaml:"value"`
}

type SecurityHeadersResultCookie struct {
	Domain   string    `json:"domain" yaml:"domain"`
	Expires  time.Time `json:"expires" yaml:"expires"`
	HttpOnly bool      `json:"httpOnly" yaml:"httpOnly"`
	MaxAge   int       `json:"maxAge" yaml:"maxAge"`
	Name     string    `json:"name" yaml:"name"`
	Path     string    `json:"path" yaml:"path"`
	Secure   bool      `json:"secure" yaml:"secure"`
	Value    string    `json:"value" yaml:"value"`
}
