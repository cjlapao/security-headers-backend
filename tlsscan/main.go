package tlsscan

import "net/url"

type TlsScanResult struct {
	Version                    string                 `json:"version" yaml:"version"`
	HandshakeComplete          bool                   `json:"handshakeComplete" yaml:"handshakeComplete"`
	DidResume                  bool                   `json:"didResume" yaml:"didResume"`
	CipherSuite                uint16                 `json:"cipherSuite" yaml:"cipherSuite"`
	NegotiatedProtocol         string                 `json:"negotiatedProtocol" yaml:"negotiatedProtocol"`
	NegotiatedProtocolIsMutual bool                   `json:"negotiatedProtocolIsMutual" yaml:"negotiatedProtocolIsMutual"`
	ServerName                 string                 `json:"serverName" yaml:"serverName"`
	Certificates               []TlsCertificateResult `json:"certificates" yaml:"certificates"`
}

type TlsCertificateResult struct {
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

type TlsScanService struct {
	Url    *url.URL      `json:"" yaml:"" bson:"url"`
	Result TlsScanResult `json:"result" yaml:"result" bson:"result"`
}

func New() *TlsScanService {
	svc := TlsScanService{}

	return &svc
}

func (svc *TlsScanService) Scan(urlToScan string) (*TlsScanResult, error) {
	parsedUrl, err := url.Parse(urlToScan)
	result := TlsScanResult{}

	if err != nil {
		return nil, err
	}

	svc.Url = parsedUrl

	return &result, nil
}
