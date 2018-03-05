package protocol

import (
	"errors"
	"bytes"
	"../utils"
)

type PacketType		uint8
type ConnectionID	uint64
type Version		uint32
type PacketNumber	uint32

const (
	PACKET_TYPE_INITIAL			= 0x7F
	PACKET_TYPE_RETRY			= 0x7E
	PACKET_TYPE_HANDSHAKE		= 0x7D
	PACKET_TYPE_0RTT_PROTECTED	= 0x7C
	PACKET_TYPE_SHORT_1_OCTET	= 0x1F
	PACKET_TYPE_SHORT_2_OCTET	= 0x1E
	PACKET_TYPE_SHORT_4_OCTET	= 0x1D
)

type Header struct {
	isLongHeader			bool
	omitConnectionIDFlag	bool
	keyPhaseBit				bool

	packetType				PacketType
	connectionID			ConnectionID
	version					Version
	packetNumber			PacketNumber
}

func (Header) Parse (b *bytes.Reader) (*Header, error) {
	firstByte, err := b.ReadByte ()
	if err != nil {
		return nil, err
	}
	b.UnreadByte ()

	if firstByte & 0x80 != 0 {
		return parseLongHeader (b)
	} else {
		return parseShortHeader (b)
	}
}

func (this *Header) Serialize (b *bytes.Buffer) error {
	if this.isLongHeader {
		return this.serializeLongHeader (b)
	} else {
		return this.serializeShortHeader (b)
	}
}

func (this *Header) SerializedLength () uint8 {
	if this.isLongHeader {
		return 1 + 8 + 4 + 4
	} else {
		var v uint8 = 1
		if this.omitConnectionIDFlag {
			v += 8
		}
		switch this.packetType {
		case PACKET_TYPE_SHORT_1_OCTET:
			v += 1
		case PACKET_TYPE_SHORT_2_OCTET:
			v += 2
		case PACKET_TYPE_SHORT_4_OCTET:
			v += 4
		}

		return v
	}
}

func parseLongHeader (b *bytes.Reader) (*Header, error) {
	headByte, _ := b.ReadByte ()
	connID, err := utils.BigEndian.ReadUInt (b, 8)
	if err != nil {
		return nil, err
	}

	version, err := utils.BigEndian.ReadUInt (b, 4)
	if err != nil {
		return nil, err
	}
	packetNumber, err := utils.BigEndian.ReadUInt (b, 4)
	if err != nil {
		return nil, err
	}
	ret := &Header {
		isLongHeader:	true,
		packetType:		PacketType (headByte & 0x7F),
		connectionID:	ConnectionID (connID),
		version: 		Version (version),
		packetNumber:	PacketNumber (packetNumber),
	}

	return ret, nil
}

func parseShortHeader (b *bytes.Reader) (*Header, error) {
	headByte, _ := b.ReadByte ()
	ret := &Header { }
	if headByte & 0x40 != 0 {
		ret.omitConnectionIDFlag = true
	} else {
		ret.omitConnectionIDFlag = false
	}
	if headByte & 0x20 != 0 {
		ret.keyPhaseBit = true
	} else {
		ret.keyPhaseBit = false
	}
	ret.packetType = PacketType (headByte & 0x1F)
	if ret.omitConnectionIDFlag == false {
		v, err := utils.BigEndian.ReadUInt (b, 8)
		if err != nil {
			return nil, err
		}
		ret.connectionID = ConnectionID (v)
	}

	switch ret.packetType {
	case PACKET_TYPE_SHORT_1_OCTET:
		v, err := b.ReadByte ()
		if err != nil {
			return nil, err
		}
		ret.packetNumber = PacketNumber (v)
	case PACKET_TYPE_SHORT_2_OCTET:
		v, err := utils.BigEndian.ReadUInt (b, 2)
		if err != nil {
			return nil, err
		}
		ret.packetNumber = PacketNumber (v)
	case PACKET_TYPE_SHORT_4_OCTET:
		v, err := utils.BigEndian.ReadUInt (b, 4)
		if err != nil {
			return nil, err
		}
		ret.packetNumber = PacketNumber (v)
	default:
		return nil, errors.New ("Header.parseShortHeader error: packet type error")
	}
	return ret, nil
}

func (this *Header) serializeLongHeader (b *bytes.Buffer) error {
	b.WriteByte (0x80 | uint8 (this.packetType))
	utils.BigEndian.WriteUInt (b, uint64 (this.connectionID), 8)
	utils.BigEndian.WriteUInt (b, uint64 (this.version), 4)
	utils.BigEndian.WriteUInt (b, uint64 (this.packetNumber), 4)
	return nil
}

func (this *Header) serializeShortHeader (b *bytes.Buffer) error {
	var flags uint8 = 0x00

	if this.omitConnectionIDFlag {
		flags |= 0x40
	}
	if this.keyPhaseBit {
		flags |= 0x20
	}
	
	b.WriteByte (flags | uint8 (this.packetType))

	if this.omitConnectionIDFlag == false {
		utils.BigEndian.WriteUInt (b, uint64 (this.connectionID), 8)
	}

	switch this.packetType {
	case PACKET_TYPE_SHORT_1_OCTET:
		b.WriteByte (uint8 (this.packetNumber))
	case PACKET_TYPE_SHORT_2_OCTET:
		utils.BigEndian.WriteUInt (b, uint64 (this.packetNumber), 2)
	case PACKET_TYPE_SHORT_4_OCTET:
		utils.BigEndian.WriteUInt (b, uint64 (this.packetNumber), 4)
	}

	return nil
}