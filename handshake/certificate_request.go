package handshake

type CertificateRequest struct {
	CertificateRequestContext	[]byte
	Extensions					[]Extension
}