package handshake

import (
	"errors"
	"bytes"
	"../protocol"
)

type transportParameterServerHandler struct {
	selfParameters		*TransportParameters
	parametersChan		chan TransportParameters

	version				protocol.Version
	supportedVersions	[]protocol.Version
}

func transportParameterServerHandlerNew(
	selfParameters		*TransportParameters,
	supportedVersions	[]protocol.Version,
	version				protocol.Version,
) *transportParameterServerHandler {
	parametersChan := make(chan TransportParameters, 1)
	return &transportParameterServerHandler {
		selfParameters:		selfParameters,
		parametersChan:		parametersChan,
		supportedVersions:	supportedVersions,
		version:			version,
	}
}

func (this *transportParameterServerHandler) Send(handshakeType HandshakeType, extensions *Extensions) error {
	if handshakeType != HANDSHAKE_TYPE_ENCRYPTED_EXTENSIONS {
		return nil
	}

	transportParameters := append(
		this.selfParameters.transferTransportParameter(),
		TransportParameter { statelessResetTokenParameterID, bytes.Repeat([]byte { 42 }, 16) },
	)
	buf := bytes.NewBuffer([]byte { })
	err := encryptedExtension { this.version, this.supportedVersions, transportParameters }.Serialize(buf)
	if err != nil {
		return err
	}

	return extensions.Add(&tlsExtension { buf.Bytes() })
}

func (this *transportParameterServerHandler) Receive(handshakeType HandshakeType, extensions *Extensions) error {
	ext := &tlsExtension { }
	founded, err := extensions.Find(ext)
	if err != nil {
		return err
	}

	if handshakeType != HANDSHAKE_TYPE_CLIENT_HELLO {
		if founded {
			return errors.New("Unexpected QUIC extension in handshake message")
		}
		return nil
	}

	if founded == false {
		return errors.New("ClientHello didn't contain a QUIC extension")
	}

	chtp := &clientHelloTransportParameters { }
	if _, err := chtp.Parse(bytes.NewReader(ext.data)); err != nil {
		return err
	}

	initialVersion := chtp.InitialVersion
	if initialVersion != this.version && initialVersion.SupportedVersions(this.supportedVersions) {
		return errors.New("Client should have used the initial version")
	}

	for _, p := range chtp.Parameters {
		if p.Parameter == statelessResetTokenParameterID {
			return errors.New("client sent a stateless reset token")
		}
	}

	parameters, err := transportParameterToTransportParameters(chtp.Parameters)
	if err != nil {
		return err
	}
	this.parametersChan <- *parameters
	return nil
}

func (this *transportParameterServerHandler) GetPeerParams() <-chan TransportParameters {
	return this.parametersChan
}