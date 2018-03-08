package packet

import (
	"errors"
	"bytes"
	"encoding/binary"
	"../protocol"
	"../utils"
)

type SupportVersion uint32

const (
	QUIC_VERSION_PICOTEST	= SupportVersion (0x50435130)
	QUIC_VERSION_GOOGLE_1	= SupportVersion (0x5130303)
	QUIC_VERSION_GOOGLE_2	= SupportVersion (0x5130313)
	QUIC_VERSION_GOOGLE_3	= SupportVersion (0x5130323)
	QUIC_VERSION_GOOGLE_4	= SupportVersion (0x5130333)
	QUIC_VERSION_GOOGLE_5	= SupportVersion (0x5130343)
	QUIC_VERSION_WINQUIC	= SupportVersion (0xABCD000)
	QUIC_VERSION_MOZILLA	= SupportVersion (0xF123F0C)
	QUIC_VERSION_FACEBOOK	= SupportVersion (0xFACEB00)
	QUIC_VERSION_IETF_1		= SupportVersion (0xFF000001)
	QUIC_VERSION_IETF_2		= SupportVersion (0xFF000002)
	QUIC_VERSION_IETF_3		= SupportVersion (0xFF000003)
	QUIC_VERSION_IETF_4		= SupportVersion (0xFF000004)
	QUIC_VERSION_IETF_5		= SupportVersion (0xFF000005)
	QUIC_VERSION_IETF_6		= SupportVersion (0xFF000006)
	QUIC_VERSION_IETF_7		= SupportVersion (0xFF000007)
	QUIC_VERSION_IETF_8		= SupportVersion (0xFF000008)
	QUIC_VERSION_ETH		= SupportVersion (0xf0f0f0f)
)

var supportedVersions = []SupportVersion { QUIC_VERSION_IETF_8 }

type VersionNegotiationPacket struct {
	header				protocol.Header
	supportedVersions	[]SupportVersion
}

func VersionNegotiationPacketParse (header protocol.Header, b *bytes.Reader) (*VersionNegotiationPacket, error) {
	if header.GetPacketType () != protocol.PACKET_TYPE_VERSON_NEGO {
		return nil, errors.New ("VersionNegotiationPacketParse error: this packet is not a  version negotiation packet")
	}

	versions := []SupportVersion { }
	versionsCount := b.Len () / 4

	buf := make ([]byte, 4)
	for i := 0; i < versionsCount; i++ {
		b.Read (buf)
		v := binary.BigEndian.Uint32 (buf)
		versions = append (versions, SupportVersion (v))
	}

	return &VersionNegotiationPacket {
		header: header,
		supportedVersions: versions,
	}, nil
}

func (this *VersionNegotiationPacket) Serialize (b *bytes.Buffer) error {
	this.header.Serialize (b)
	for _, v := range this.supportedVersions {
		utils.BigEndian.WriteUInt (b, uint64 (v), 4)
	}
	return nil
}

func IsSupportedVersion (supported []SupportVersion, version SupportVersion) bool {
	for _, v := range supported {
		if version == v {
			return true
		}
	}
	return false
}

func (this *VersionNegotiationPacket) ChooseSupportedVersion (supportedVersion []SupportVersion) (SupportVersion, bool) {
	for _, ourVersion := range supportedVersion {
		for _, theirVersion := range this.supportedVersions {
			if ourVersion == theirVersion {
				return ourVersion, true
			}
		}
	}

	return 0, false
}

