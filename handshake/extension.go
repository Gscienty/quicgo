package handshake

import (
	"errors"
	"bytes"
	"../utils"
)

type ExtensionType uint16

const (
	EXTENSION_TYPE_SERVER_NAME								ExtensionType = 0
	EXTENSION_TYPE_MAX_FRAGMENT_LENGTH						ExtensionType = 1
	EXTENSION_TYPE_STATUS_REQUEST							ExtensionType = 5
	EXTENSION_TYPE_SUPPORTED_GROUPS							ExtensionType = 10
	EXTENSION_TYPE_SIGNATURE_ALGORITHMS						ExtensionType = 13
	EXTENSION_TYPE_USE_SRTP									ExtensionType = 14
	EXTENSION_TYPE_HEARTBEAT								ExtensionType = 15
	EXTENSION_TYPE_APPLICATION_LAYER_PROTOCOL_NEGOTIATION	ExtensionType = 16
	EXTENSION_TYPE_SIGNED_CERTIFICATE_TIMESTAMP				ExtensionType = 18
	EXTENSION_TYPE_CLIENT_CERTIFICATE_TYPE					ExtensionType = 19
	EXTENSION_TYPE_SERVER_CERTIFICATE_TYPE					ExtensionType = 20
	EXTENSION_TYPE_PADDING									ExtensionType = 21
	EXTENSION_TYPE_PRE_SHARED_KEY							ExtensionType = 41
	EXTENSION_TYPE_EARLY_DATA								ExtensionType = 42
	EXTENSION_TYPE_SUPPORTED_VERSIONS						ExtensionType = 43
	EXTENSION_TYPE_COOKIE									ExtensionType = 44
	EXTENSION_TYPE_PSK_KEY_EXCHANGE_MODES					ExtensionType = 45
	EXTENSION_TYPE_CERTIFICATE_AUTHORITIES					ExtensionType = 47
	EXTENSION_TYPE_OID_FILTERS								ExtensionType = 48
	EXTENSION_TYPE_POST_HANDSHAKE_AUTH						ExtensionType = 49
	EXTENSION_TYPE_SIGNATURE_ALGORITHMS_CERT				ExtensionType = 50
	EXTENSION_TYPE_KEY_SHARE								ExtensionType = 51

	EXTENSION_TYPE_QUIC_GSCI_TLS							ExtensionType = 66
)

type IExtension interface {
	Type() ExtensionType
	Serialize(b *bytes.Buffer) error
	Parse(b *bytes.Reader) (int, error)
}

type Extension struct {
	Type	ExtensionType
	Data	[]byte
}

func (this *Extension) Serialize(b *bytes.Buffer) error {
	utils.BigEndian.WriteUInt(b, uint64(this.Type), 2)
	writedLen, err := b.Write(this.Data)
	if err != nil {
		return err
	}
	if writedLen != len(this.Data) {
		return errors.New("wrong length")
	}
	return nil
}

type Extensions []Extension

func (this Extensions) Serialize(b *bytes.Buffer) error {
	utils.BigEndian.WriteUInt(b, uint64(len(this)), 2)
	for _, v := range this {
		err := v.Serialize(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Extensions) Add(e IExtension) error {
	b := bytes.NewBuffer([]byte { })
	err := e.Serialize(b)
	if err != nil {
		return err
	}
	for i := range *this {
		if (*this)[i].Type == e.Type() {
			(*this)[i].Data = b.Bytes()
			return nil
		}
	}

	*this = append(*this, Extension { e.Type(), b.Bytes() })
	return nil
}

func (this Extensions) Find(dst IExtension) (bool, error) {
	for _, ext := range this {
		if ext.Type == dst.Type() {
			_, err := dst.Parse(bytes.NewReader(ext.Data))
			if err != nil {
				return true, err
			}
			return true, nil
		}
	}
	return false, nil
}