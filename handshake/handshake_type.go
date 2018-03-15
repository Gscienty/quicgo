package handshake

type HandshakeType uint8

const (
	HANDSHAKE_TYPE_CLIENT_HELLO			HandshakeType = 1
	HANDSHAKE_TYPE_SERVER_HELLO			HandshakeType = 2
	HANDSHAKE_TYPE_NEW_SESSION_TICKET	HandshakeType = 4
	HANDSHAKE_TYPE_END_OF_EARLY_DATA	HandshakeType = 5
	HANDSHAKE_TYPE_HELLO_RETRY_REQUEST	HandshakeType = 6
	HANDSHAKE_TYPE_ENCRYPTED_EXTENSIONS	HandshakeType = 8
	HANDSHAKE_TYPE_CERTIFICATE			HandshakeType = 11
	HANDSHAKE_TYPE_CERTIFICATE_REQUEST	HandshakeType = 13
	HANDSHAKE_TYPE_CERTIFICATE_VERIFY	HandshakeType = 15
	HANDSHAKE_TYPE_SERVER_CONFIGURATION	HandshakeType = 17
	HANDSHAKE_TYPE_FINISHED				HandshakeType = 20
	HANDSHAKE_TYPE_KEY_UPDATE			HandshakeType = 24
	HANDSHAKE_TYPE_MESSAGE_HASH			HandshakeType = 254
)