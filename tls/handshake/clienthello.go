package handshake

type ClientHello struct {
	LegacyVersion 	uint16
	Random			[32]byte
	LegacySessionID	[]byte
	CipherSuites	[]CipherSuite
	Extensions		[]Extension
}