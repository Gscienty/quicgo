package handshake

type CertificateRequest struct {
	Context		[]byte
	Extensions	[]Extension
}