package protocol

import "errors"
import "encoding/binary"

const (
	PUBLIC_FLAG_VERSION			= 0x01
	PUBLIC_FLAG_RESET			= 0x02

	PUBLIC_FLAG_CONNID_8BYTE	= 0x0C
	PUBLIC_FLAG_CONNID_4BYTE	= 0x08
	PUBLIC_FLAG_CONNID_1BYTE	= 0x04
	PUBLIC_FLAG_CONNID_0BYTE	= 0x00

	PUBLIC_FLAG_PACKNUM_6BYTE	= 0x30
	PUBLIC_FLAG_PACKNUM_4BYTE	= 0x20
	PUBLIC_FLAG_PACKNUM_2BYTE	= 0x10
	PUBLIC_FLAG_PACKNUM_1BYTE	= 0x00

	PUBLIC_FLAG_MULTIPATH		= 0x40
	PUBLIC_FLAG_UNUSED			= 0x80
)

var _parseConnectionIdSize	= [] int { 0, 1, 4, 8 }
var _parsePackectNumberSize	= [] int { 1, 2 ,4, 6 }

type PublicHeader struct {
	flags			uint8
	connectionID 	uint64
	version 		uint32
	packetNumber	uint64
}

func (this *PublicHeader) Initialize () {
	this.flags = 0
	this.connectionID =  0
	this.version = 0
	this.packetNumber = 0
}

func (this *PublicHeader) Parse (data []byte) (size int, err error) {
	dataLen := len (data)
	if dataLen < 2 {
		err = errors.New ("QUIC PublicHeader parse error: data too small")
		return
	}

	this.flags = data[0]
	size = 1
	switch this.GetConnectionSize () {
	case 0:
		this.connectionID = 0
		break
	case 1:
		this.connectionID = uint64 (data[size])
		size += 1
		break
	case 4:
		this.connectionID = uint64 (binary.LittleEndian.Uint32 (data[size:]))
		size += 4
		break
	case 8:
		this.connectionID = uint64 (binary.LittleEndian.Uint64 (data[size:]))
		size += 8
		break
	}

	if this.ExistVersion () {
		this.version = uint32 (binary.LittleEndian.Uint32 (data[size:]))
		size += 4
	}

	switch this.GetPacketNumberSize () {
	case 1:
		this.packetNumber = uint64 (data[size])
		size += 1
		break
	case 2:
		this.packetNumber = uint64 (binary.LittleEndian.Uint16 (data[size:]))
		size += 2
		break
	case 4:
		this.packetNumber = uint64 (binary.LittleEndian.Uint32 (data[size:]))
		size += 4
		break
	case 6:
		this.packetNumber = uint64 (binary.LittleEndian.Uint32 (data[size:])) +
			(uint64 (binary.LittleEndian.Uint16 (data[size + 4:])) << 32)
		size += 6
		break
	}

	if size != this.GetSerializedSize () {
		err = errors.New ("QUIC PublicHeader parse error: internal error, calculate size different")
	}
	
	return
}

func (this *PublicHeader) GetSerializedData (data []byte, isServer bool) (size int, err error) {
	if this.IsReset () {
		// special packet: reset packet
		if len (data) < 9 {
			err = errors.New ("QUIC PublicHeader serialize: data too small")
			return
		}

		data[0] = PUBLIC_FLAG_RESET | PUBLIC_FLAG_CONNID_8BYTE
		binary.LittleEndian.PutUint64 (data[1:], uint64 (this.connectionID))
		size = 9
		return
	} else if this.ExistVersion () && isServer {
		// special packet: version negotiation packet
		if len (data) < 9 {
			err = errors.New ("QUIC PublicHeader serialize: data too small")
			return
		}

		data[0] = PUBLIC_FLAG_VERSION | PUBLIC_FLAG_CONNID_8BYTE
		binary.LittleEndian.PutUint64 (data[1:], uint64 (this.connectionID))
		size = 9
		return
	}

	// regular packet
	if this.GetSerializedSize () > len (data) {
		err = errors.New ("QUIC PublicHeader serialize: data too small")
		return
	}
	size = 1
	data[0] = this.flags
	switch this.GetConnectionSize () {
	case 0:
		break
	case 1:
		data[1] = uint8 (this.connectionID)
		size += 1
		break
	case 4:
		binary.LittleEndian.PutUint32 (data[1:], uint32 (this.connectionID))
		size += 4
		break
	case 8:
		binary.LittleEndian.PutUint64 (data[1:], uint64 (this.connectionID))
		size += 8
		break
	}

	if this.ExistVersion () {
		binary.LittleEndian.PutUint32 (data[size:], uint32 (this.version))
		size += 4
	}

	switch this.GetPacketNumberSize () {
	case 1:
		data[size] = uint8 (this.packetNumber)
		size += 1
		break
	case 2:
		binary.LittleEndian.PutUint16 (data[size:], uint16 (this.packetNumber))
		size += 2
		break
	case 4:
		binary.LittleEndian.PutUint32 (data[size:], uint32 (this.packetNumber))
		size += 4
		break
	case 6:
		binary.LittleEndian.PutUint32 (data[size:], uint32 (this.packetNumber))
		binary.LittleEndian.PutUint16 (data[size + 4:], uint16 (this.packetNumber >> 32))
		size += 6
		break
	}

	if size != this.GetSerializedSize () {
		err = errors.New ("QUIC PublicHeader serialize: internal error, calculate size different")
	}
	return
}

func (this *PublicHeader) GetSerializedSize () (size int) {
	size = (1 + _parseConnectionIdSize[this.GetConnectionSize ()] + _parsePackectNumberSize[this.GetPacketNumberSize ()])
	if this.flags & PUBLIC_FLAG_VERSION != 0 {
		size += 4
	}
	return
}

func (this *PublicHeader) GetConnectionSize () (size int) {
	size = int ((this.flags >> 2) & 0x03)
	return
}

func (this *PublicHeader) SetConnectionSize (size int) (err error) {
	switch size {
	case 0:
		this.flags = (this.flags & 0xF3) | PUBLIC_FLAG_CONNID_0BYTE
		return
	case 1:
		this.flags = (this.flags & 0xF3) | PUBLIC_FLAG_CONNID_1BYTE
		return
	case 4:
		this.flags = (this.flags & 0xF3) | PUBLIC_FLAG_CONNID_4BYTE
		return
	case 8:
		this.flags = (this.flags & 0xF3) | PUBLIC_FLAG_CONNID_8BYTE
		return
	}

	err = errors.New ("Set Connection Size: size error")
	return
}

func (this *PublicHeader) GetPacketNumberSize () (size int) {
	size = int ((this.flags >> 4) & 0x03)
	return
}

func (this *PublicHeader) SetPacketNumberSize (size int) (err error) {
	switch size {
	case 1:
		this.flags = (this.flags & 0xCF) | PUBLIC_FLAG_PACKNUM_1BYTE
		return
	case 2:
		this.flags = (this.flags & 0xCF) | PUBLIC_FLAG_PACKNUM_2BYTE
		return
	case 4:
		this.flags = (this.flags & 0xCF) | PUBLIC_FLAG_PACKNUM_4BYTE
		return
	case 6:
		this.flags = (this.flags & 0xCF) | PUBLIC_FLAG_PACKNUM_6BYTE
		return
	}

	err = errors.New ("Set PacketNumber Size: size error")
	return
}

func (this *PublicHeader) ExistVersion () (exist bool) {
	exist = (this.flags & PUBLIC_FLAG_VERSION) != 0
	return
}

func (this *PublicHeader) GetVersion () (version uint32) {
	version = this.version
	return
}

func (this *PublicHeader) SetVersion (version uint32) {
	this.flags |= PUBLIC_FLAG_VERSION
	this.version = version
	return
}

func (this *PublicHeader) IsReset () (isReset bool) {
	isReset = (this.flags & PUBLIC_FLAG_RESET) != 0
	return
}

func (this *PublicHeader) SetReset (reset bool) {
	if reset {
		this.flags |= PUBLIC_FLAG_RESET
	} else {
		this.flags &= 0xFF ^ PUBLIC_FLAG_RESET
	}
	return
}