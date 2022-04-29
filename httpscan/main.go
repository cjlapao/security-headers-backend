package httpscan

import (
	"net/http"
	"net/url"
	"time"

	"github.com/cjlapao/security-headers-backend/tlsscan"
)

type HttpScanResult struct {
	Body          interface{}       `json:"body" yaml:"body"`
	ContentLength int64             `json:"contentLength" yaml:"contentLength"`
	Headers       map[string]string `json:"headers" yaml:"headers"`
	Cookies       map[string]HttpScanResultCookies
	Protocol      string                       `json:"protocol" yaml:"protocol"`
	ProtocolMajor int                          `json:"protocolMajor" yaml:"protocolMajor"`
	ProtocolMinor int                          `json:"protocolMinor" yaml:"protocolMinor"`
	Status        string                       `json:"status" yaml:"status"`
	StatusCode    int                          `json:"statusCode" yaml:"statusCode"`
	Tls           tlsscan.TlsScanResult        `json:"tls" yaml:"tls"`
	Certificates  tlsscan.TlsCertificateResult `json:"certificates" yaml:"certificates"`
}

type HttpScanResultCookies struct {
	Domain   string    `json:"domain" yaml:"domain"`
	Expires  time.Time `json:"expires" yaml:"expires"`
	HttpOnly bool      `json:"httpOnly" yaml:"httpOnly"`
	MaxAge   int       `json:"maxAge" yaml:"maxAge"`
	Name     string    `json:"name" yaml:"name"`
	Path     string    `json:"path" yaml:"path"`
	Secure   bool      `json:"secure" yaml:"secure"`
	SameSite string    `json:"sameSite" yaml:"sameSite"`
	Value    string    `json:"value" yaml:"value"`
}

type HttpScanService struct {
	Url      *url.URL       `json:"url" yaml:"url" bson:"url"`
	Request  *http.Request  `json:"" yaml:"" bson:""`
	Response *http.Response `json:"" yaml:"" bson:""`
	Result   HttpScanResult `json:"result" yaml:"result" bson:"result"`
}

func New() *HttpScanService {
	svc := HttpScanService{}

	return &svc
}

func (svc HttpScanService) Scan(urlToScan string) (*HttpScanResult, error) {
	result := HttpScanResult{}

	return &result, nil
}
