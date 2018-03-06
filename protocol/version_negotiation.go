package protocol

import (
	"errors"
	"bytes"
	"encoding/binary"
	"../utils"
)

type VersionNegotiationPacket struct {
	header				Header
	supportedVersions	[]uint32
}

func VersionNegotiationPacketParse (header Header, b *bytes.Reader) (*VersionNegotiationPacket, error) {
	if header.packetType != PACKET_TYPE_VERSON_NEGO {
		return nil, errors.New ("VersionNegotiationPacketParse error: this packet is not a  version negotiation packet")
	}

	versions := []uint32 { }
	versionsCount := b.Len () / 4

	buf := make ([]byte, 4)
	for i := 0; i < versionsCount; i++ {
		b.Read (buf)
		v := binary.BigEndian.Uint32 (buf)
		versions = append (versions, v)
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