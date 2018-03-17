package handshake

import (
	"errors"
	"io"
	"crypto/aes"
	"crypto/sha256"
	"crypto/cipher"
	"time"
	"crypto/rand"
	"../crypto"
)

type Cookie struct {
	RemoteAddr	string
	SentTime	time.Time
}

type token struct {
	Data		[]byte
	Timestamp	int64
}

type CookieProtector struct {
	secret	[]byte
}

const cookieSecretSize = 32
const cookieNonceSize = 32

func DefaultCookieProtectorNew() (*CookieProtector, error) {
	secret := make([]byte, cookieSecretSize)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	return &CookieProtector { secret }, nil
}

func (this *CookieProtector) NewToken(data []byte, info []byte) ([]byte, error) {
	nonce := make([]byte, cookieNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	aead, aeadNonce, err := this.createAEAD(nonce, info)
	if err != nil {
		return nil, err
	}
	return append(nonce, aead.Seal(nil, aeadNonce, data, nil)...), nil
}

func (this *CookieProtector) DecodeToken(p []byte, info []byte) ([]byte, error) {
	if len(p) < cookieNonceSize {
		return nil, errors.New("Token too short")
	}
	nonce := p[:cookieNonceSize]
	aead, aeadNonce, err := this.createAEAD(nonce, info)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, aeadNonce, p[cookieNonceSize:], nil)
}

func (this *CookieProtector) createAEAD(nonce []byte, info []byte) (cipher.AEAD, []byte, error) {
	h := crypto.HKDFNew(sha256.New, this.secret, nonce, info)
	key := make([]byte, 32)
	if _, err := io.ReadFull(h, key); err != nil {
		return nil, nil, err
	}
	aeadNonce := make([]byte, 12)
	if _, err := io.ReadFull(h, aeadNonce); err != nil {
		return nil, nil, err
	}

	cp, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aead, err := cipher.NewGCM(cp)
	if err != nil {
		return nil, nil, err
	}

	return aead, aeadNonce, nil
}