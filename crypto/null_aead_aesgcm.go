package crypto

import (
	"encoding/binary"
	"crypto"
	"../protocol"
)

var quicVersion1Salt = []byte{0xaf, 0xc8, 0x24, 0xec, 0x5f, 0xc7, 0x7e, 0xca, 0x1e, 0x9d, 0x36, 0xf3, 0x7f, 0xb2, 0xd4, 0x65, 0x18, 0xc3, 0x66, 0x39}

func nullAeadAesGcmNew(connectionID protocol.ConnectionID, isClient bool) (AEAD, error) {
	cs, ss := computeSecrets(connectionID)

	clientKey, clientIV := computeNullAeadParameters(cs)
	serverKey, serverIV := computeNullAeadParameters(ss)

	if isClient {
		return AeadAsmGcmNew(clientKey, serverKey, clientIV, serverIV)
	} else {
		return AeadAsmGcmNew(serverKey, clientKey, clientIV, serverIV)
	}
}

func NullAeadAesGcmClientNew(connectionID protocol.ConnectionID) (AEAD, error) {
	return nullAeadAesGcmNew(connectionID, true)
}

func NullAeadAesGcmServerNew(connectionID protocol.ConnectionID) (AEAD, error) {
	return nullAeadAesGcmNew(connectionID, false)
}

func computeSecrets(connectionID protocol.ConnectionID) ([]byte, []byte) {
	connID := make([]byte, 8)
	binary.BigEndian.PutUint64(connID, uint64(connectionID))
	handshakeSecret, _ := HKDFExtract(crypto.SHA256, quicVersion1Salt, connID)
	clientSecret := QuicHKDFExtract(handshakeSecret, "client hs", crypto.SHA256.Size())
	serverSecret := QuicHKDFExtract(handshakeSecret, "server hs", crypto.SHA256.Size())
	return clientSecret, serverSecret
}

func computeNullAeadParameters(secret []byte) ([]byte, []byte) {
	key := QuicHKDFExtract(secret, "key", 16)
	iv := QuicHKDFExtract(secret, "iv", 12)
	return key, iv
}