package crypto

import (
	"encoding/binary"
	"errors"
	"crypto/aes"
	"crypto/cipher"
	"../protocol"
)

type AeadAsmGcm struct {
	selfIV		[]byte
	otherIV		[]byte
	encrypter	cipher.AEAD
	decrypter	cipher.AEAD
}

const IV_LEN = 12

var _ AEAD = &AeadAsmGcm { }

func AeadAsmGcmNew(selfKey []byte, otherKey []byte, selfIV []byte, otherIV []byte) (AEAD, error) {
	if len(selfIV) != IV_LEN || len(otherIV) != IV_LEN {
		return nil, errors.New("AES-GCM: expected 12 byte IVs")
	}

	encrypterCipher, err := aes.NewCipher(selfKey)
	if err != nil {
		return nil, err
	}
	encrypter, err := cipher.NewGCM(encrypterCipher)
	if err != nil {
		return nil, err
	}
	decrypterCipher, err := aes.NewCipher(otherKey)
	if err != nil {
		return nil, err
	}
	decrypter, err := cipher.NewGCM(decrypterCipher)
	if err != nil {
		return nil, err
	}

	return &AeadAsmGcm {
		selfIV:		selfIV,
		otherIV:	otherIV,
		encrypter:	encrypter,
		decrypter:	decrypter,
	}, nil
}

func (this *AeadAsmGcm) makeNonce(iv []byte, packetNumber protocol.PacketNumber) []byte {
	nonce := make([]byte, IV_LEN)

	binary.BigEndian.PutUint64(nonce[IV_LEN - 8:], uint64(packetNumber))
	for i := 0; i < IV_LEN; i++ {
		nonce[i] ^= iv[i]
	}
	return nonce
}

func (this *AeadAsmGcm) Open(dst []byte, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, error) {
	return this.decrypter.Open(dst, this.makeNonce(this.otherIV, packetNumber), src, associatedData)
}

func (this *AeadAsmGcm) Seal(dst []byte, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) []byte {
	return this.encrypter.Seal(dst, this.makeNonce(this.selfIV, packetNumber), src, associatedData)
}

func (this *AeadAsmGcm) Overhead() int {
	return this.encrypter.Overhead()
}