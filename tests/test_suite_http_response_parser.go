package tests

import (
	"net/http"
	"strings"
	"time"
)

type TestSuiteHttpResponse struct {
	Body            interface{}       `json:"body" yaml:"body"`
	ContentLength   int64             `json:"contentLength" yaml:"contentLength"`
	Headers         map[string]string `json:"headers" yaml:"headers"`
	Cookies         map[string]TestSuiteHttpResponseCookie
	Protocol        string                   `json:"protocol" yaml:"protocol"`
	ProtocolMajor   int                      `json:"protocolMajor" yaml:"protocolMajor"`
	ProtocolMinor   int                      `json:"protocolMinor" yaml:"protocolMinor"`
	Status          string                   `json:"status" yaml:"status"`
	StatusCode      int                      `json:"statusCode" yaml:"statusCode"`
	Tls             TestSuiteHttpResponseTls `json:"tls" yaml:"tls"`
	TlsCertificates []TestSuiteTlsCertificate
}

type TestSuiteHttpResponseTls struct {
	Version                    string `json:"version" yaml:"version"`
	HandshakeComplete          bool   `json:"handshakeComplete" yaml:"handshakeComplete"`
	DidResume                  bool   `json:"didResume" yaml:"didResume"`
	CipherSuite                uint16 `json:"cipherSuite" yaml:"cipherSuite"`
	NegotiatedProtocol         string `json:"negotiatedProtocol" yaml:"negotiatedProtocol"`
	NegotiatedProtocolIsMutual bool   `json:"negotiatedProtocolIsMutual" yaml:"negotiatedProtocolIsMutual"`
	ServerName                 string `json:"serverName" yaml:"serverName"`
}

type TestSuiteHttpResponseCookie struct {
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

type TestSuiteTlsCertificate struct {
	Subject            string
	CommonNames        string
	AlternativeNames   string
	SerialNumber       string
	ValidFrom          string
	ValidUntil         string
	Key                string
	Issuer             string
	SignatureAlgorithm string
	IsCaa              bool
	Trusted            bool
}

func (httpResponse TestSuiteHttpResponse) Parse(response *http.Response) *TestSuiteHttpResponse {
	result := TestSuiteHttpResponse{
		ContentLength: response.ContentLength,
		Protocol:      response.Proto,
		ProtocolMajor: response.ProtoMajor,
		ProtocolMinor: response.ProtoMinor,
		Status:        response.Status,
		StatusCode:    response.StatusCode,
		Headers:       make(map[string]string),
		Cookies:       make(map[string]TestSuiteHttpResponseCookie),
		Tls:           TestSuiteHttpResponseTls{},
	}

	// Extracting all keys
	for key, values := range response.Header {
		result.Headers[strings.ToLower(key)] = strings.Join(values, ";")
	}

	// Extracting all cookies
	for _, cookie := range response.Cookies() {
		resultCookie := TestSuiteHttpResponseCookie{
			Domain:   cookie.Domain,
			Expires:  cookie.Expires,
			HttpOnly: cookie.HttpOnly,
			MaxAge:   cookie.MaxAge,
			Name:     cookie.Name,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			Value:    cookie.Value,
			SameSite: parseCookieSameSite(cookie.SameSite),
		}

		result.Cookies[cookie.Name] = resultCookie
	}

	// extracting body
	// TODO: implement body extraction strategy

	// Extracting TLS
	if response.TLS != nil {
		result.Tls.Version = parseTlsVersion(response.TLS.Version)
		result.Tls.DidResume = response.TLS.DidResume
		result.Tls.HandshakeComplete = response.TLS.HandshakeComplete
		result.Tls.NegotiatedProtocol = response.TLS.NegotiatedProtocol
		result.Tls.NegotiatedProtocolIsMutual = response.TLS.NegotiatedProtocolIsMutual
		result.Tls.ServerName = response.TLS.ServerName
		result.Tls.CipherSuite = response.TLS.CipherSuite
		for _, cert := range response.TLS.PeerCertificates {
			tlsCert := TestSuiteTlsCertificate{}
			tlsCert.AlternativeNames = strings.Join(cert.DNSNames, ",")
			tlsCert.IsCaa = cert.IsCA
			tlsCert.SerialNumber = cert.Issuer.SerialNumber
			tlsCert.Issuer = cert.Issuer.CommonName
			tlsCert.Key = cert.PublicKeyAlgorithm.String()
			tlsCert.SignatureAlgorithm = cert.SignatureAlgorithm.String()
			tlsCert.Subject = cert.Subject.CommonName
			tlsCert.ValidFrom = cert.NotBefore.Format(time.RFC3339)
			tlsCert.ValidUntil = cert.NotAfter.Format(time.RFC3339)

			result.TlsCertificates = append(result.TlsCertificates, tlsCert)
		}
	}

	return &result
}

func parseCookieSameSite(val http.SameSite) string {
	switch val {
	case http.SameSiteDefaultMode:
		return "default"
	case http.SameSiteLaxMode:
		return "lax"
	case http.SameSiteStrictMode:
		return "strict"
	case http.SameSiteNoneMode:
		return "none"
	default:
		return "unknown"
	}
}

func parseTlsVersion(val uint16) string {
	switch val {
	case 769:
		return "tls1.0"
	case 770:
		return "tls1.1"
	case 771:
		return "tls1.2"
	case 772:
		return "tls1.3"
	default:
		return "unknown"
	}
}
