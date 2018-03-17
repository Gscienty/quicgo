package crypto

import (
	"io"
	"hash"
	"encoding/binary"
	"crypto/hmac"
	"bytes"
	"crypto"
	"errors"
)

type HKDF struct {
	expander	hash.Hash
	size		int
	info		[]byte
	counter		byte
	prev		[]byte
	cache		[]byte
}

var _ io.Reader = &HKDF {}

func (this *HKDF) Read(p []byte) (int, error) {
	need := len(p)
	remains := len(this.cache) + int(255 - this.counter + 1) * this.size
	if remains < need {
		return 0, errors.New("hkdf: present")
	}

	n := copy(p, this.cache)
	p = p[n:]

	for len(p) > 0 {
		this.expander.Reset()
		this.expander.Write(this.prev)
		this.expander.Write(this.info)
		this.expander.Write([]byte{ this.counter })
		this.prev = this.expander.Sum(this.prev[:0])
		this.counter++

		this.cache = this.prev
		n = copy(p, this.cache)
		p = p[n:]
	}

	this.cache = this.cache[n:]

	return need, nil
}

func HKDFNew(hash func() hash.Hash, secret, salt, info []byte) io.Reader {
	if salt == nil {
		salt = make([]byte, hash().Size())
	}
	extractor := hmac.New(hash, salt)
	extractor.Write(secret)
	prk := extractor.Sum(nil)

	return &HKDF { hmac.New(hash, prk), extractor.Size(), info, 1, nil, nil }
}

func HKDFExpand(hash crypto.Hash, prk []byte, info []byte, outLen int) []byte {
	out := []byte { }
	T := []byte { }
	i := byte(1)
	for len(out) < outLen {
		block := append(T, info...)
		block = append(block, i)

		h := hmac.New(hash.New, prk)
		h.Write(block)

		T = h.Sum(nil)
		out = append(out, T...)
		i++
	}

	return out[:outLen]
}

func HKDFExtract(hash crypto.Hash, saltIn []byte, data []byte) ([]byte, hash.Hash) {
	salt := saltIn

	if salt == nil {
		salt = bytes.Repeat([]byte { 0 }, hash.Size())
	}

	h := hmac.New(hash.New, salt)
	h.Write(data)
	out := h.Sum(nil)

	return out, h
}

func QuicHKDFExtract(secret []byte, label string, length int) []byte {
	quicLabel := make([]byte, 2 + 1 + 5 + len(label) + 1)
	binary.BigEndian.PutUint16(quicLabel[0:2], uint16(length))
	quicLabel[2] = uint8(5 + len(label))
	copy(quicLabel[3:], []byte("QUIC " + label))
	return HKDFExpand(crypto.SHA256, secret, quicLabel, length)
}