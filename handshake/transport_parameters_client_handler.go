package handshake

import (
	"errors"
	"bytes"
	"../protocol"
	"../utils"
)

type transportParameterClientHandler struct {
	selfParameters		*TransportParameters
	parametersChan		chan TransportParameters

	initialVersion		protocol.Version
	supportedVersions	[]protocol.Version
	version				protocol.Version
}

type clientHelloTransportParameters struct {
	InitialVersion	protocol.Version
	Parameters		[]TransportParameter
}

func (this clientHelloTransportParameters) Serialize(b *bytes.Buffer) error {
	utils.BigEndian.WriteUInt(b, uint64(this.InitialVersion), 4)
	parametersBuffer := bytes.NewBuffer([]byte { })
	for _, p := range this.Parameters {
		err := p.Serialize(parametersBuffer)
		if err != nil {
			return err
		}
	}
	utils.BigEndian.WriteUInt(b, uint64(parametersBuffer.Len()), 2)
	_, err := b.Write(parametersBuffer.Bytes())
	return err
}

func transportParameterClientHandlerNew(
	selfParameters		*TransportParameters,
	initialVersion		protocol.Version,
	supportedVersions	[]protocol.Version,
	version				protocol.Version,
) *transportParameterClientHandler {
	parametersChan := make(chan TransportParameters)
	return &transportParameterClientHandler {
		selfParameters:		selfParameters,
		parametersChan:		parametersChan,
		initialVersion:		initialVersion,
		supportedVersions:	supportedVersions,
		version:			version,
	}
}

func (this *transportParameterClientHandler) Send(handshakeType HandshakeType, extensions *Extensions) error {
	if handshakeType != HANDSHAKE_TYPE_CLIENT_HELLO {
		return nil
	}

	buf := bytes.NewBuffer([]byte { })
	err := clientHelloTransportParameters { this.initialVersion, this.selfParameters.transferTransportParameter() }.Serialize(buf)
	if err != nil {
		return err
	}
	return extensions.Add(&tlsExtension { buf.Bytes() })
}

func (this *transportParameterClientHandler) Receive(handshakeType HandshakeType, extensions *Extensions) error {
	ext := &tlsExtension { }
	founded, err := extensions.Find(ext)
	if err != nil {
		return err
	}

	if handshakeType != HANDSHAKE_TYPE_ENCRYPTED_EXTENSIONS {
		if founded {
			return errors.New("Unexpected QUIC extension in handshake message")
		}
		return nil
	}
	if founded == false {
		return errors.New("EncryptedExtensions message didn't contain a QUIC extension")
	}

	eetp := &encryptedExtension { }
	if _, err = eetp.Parse(bytes.NewReader(ext.data)); err != nil {
		return err
	}

	return nil
}