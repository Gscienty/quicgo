package handshake

import (
	"encoding/asn1"
	"net"
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
	info	[]byte
	secret	[]byte
}

const cookieSecretSize = 32
const cookieNonceSize = 32

const (
	cookiePrefixIP = iota
	cookiePrefixString
)

func CookieProtectorNew(info []byte) (*CookieProtector, error) {
	secret := make([]byte, cookieSecretSize)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	return &CookieProtector { secret, info }, nil
}

func encodeRemoteAddr(remoteAddr net.Addr) []byte {
	if udpAddr, ok := remoteAddr.(*net.UDPAddr); ok {
		return append([]byte { cookiePrefixIP }, udpAddr.IP...)
	}
	return append([]byte { cookiePrefixString }, []byte(remoteAddr.String())...)
}

func decodeRemoteAddr(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	if data[0] == cookiePrefixIP {
		return net.IP(data[1:]).String()
	}
	return string(data[1:])
}

func (this *CookieProtector) NewToken(remoteAddr net.Addr) ([]byte, error) {
	data, err := asn1.Marshal(token {
		Data:		encodeRemoteAddr(remoteAddr),
		Timestamp:	time.Now().Unix(),
	});

	if err != nil {
		return nil, err
	}
	return this.newToken(data)
}

func (this *CookieProtector) DecodeToken(data []byte) (*Cookie, error) {
	if len(data) == 0 {
		return nil, nil
	}

	data, err := this.decodeToken(data)
	if err != nil {
		return nil, err
	}
	t := &token { }
	rest, err := asn1.Unmarshal(data, t)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("rest when unpacking token")
	}
	return &Cookie {
		RemoteAddr:	decodeRemoteAddr(t.Data),
		SentTime:	time.Unix(t.Timestamp, 0),
	}, nil
}

func (this *CookieProtector) newToken(data []byte) ([]byte, error) {
	nonce := make([]byte, cookieNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	aead, aeadNonce, err := this.createAEAD(nonce)
	if err != nil {
		return nil, err
	}
	return append(nonce, aead.Seal(nil, aeadNonce, data, nil)...), nil
}

func (this *CookieProtector) decodeToken(p []byte) ([]byte, error) {
	if len(p) < cookieNonceSize {
		return nil, errors.New("Token too short")
	}
	nonce := p[:cookieNonceSize]
	aead, aeadNonce, err := this.createAEAD(nonce)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, aeadNonce, p[cookieNonceSize:], nil)
}

func (this *CookieProtector) createAEAD(nonce []byte) (cipher.AEAD, []byte, error) {
	h := crypto.HKDFNew(sha256.New, this.secret, nonce, this.info)
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