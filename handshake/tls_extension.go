package handshake

import (
	"errors"
	"bytes"
)

type tlsExtension struct {
	data []byte
}

var _ IExtension = &tlsExtension { }

func (this *tlsExtension) Type() ExtensionType {
	return EXTENSION_TYPE_QUIC_GSCI_TLS
}

func (this *tlsExtension) Serialize(b *bytes.Buffer) error {
	writedLen, err := b.Write(this.data)
	if err != nil {
		return err
	}
	if writedLen != len(this.data) {
		return errors.New("wrong length")
	}
	return nil
}

func (this *tlsExtension) Parse(b *bytes.Reader) (int, error) {
	length := b.Len()

	buf := make([]byte, length)
	readedLen, err := b.Read(buf)
	if err != nil {
		return 0, err
	}
	if readedLen != length {
		return 0, errors.New("internal wrong")
	}

	return length, nil
}