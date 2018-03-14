package handshake

type ServerHello struct {
	LegacyVersion			uint16
	Random					[32]byte
	LegacySessionIDEcho		[]byte
	CipherSuite				CipherSuite
	LegacyCompressionMethod	uint8
	Extensions				[]Extension
}