package securityheaders

import "net/http"

type SecurityHeadersService struct {
	Headers http.Header
	Result  SecurityHeadersResult
}

func New(headers http.Header) *SecurityHeadersService {
	svc := SecurityHeadersService{
		Headers: headers,
		Result: SecurityHeadersResult{
			Points: 0,
		},
	}

	return &svc
}

func (s *SecurityHeadersService) Calculate() SecurityHeadersResult {

}
