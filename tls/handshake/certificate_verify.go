package handshake

type CertificateVerify struct {
	SignatureScheme uint16
	Signature		[]byte
}