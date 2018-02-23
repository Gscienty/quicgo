package packet

import "errors"
import "encoding/binary"

type Header struct {
	packetType uint8
	connId uint64
	version uint32
	packetNumber uint32

	isLongHeader bool
	keyPhase bool
	connIdOmitted bool
}

func (this *Header) Parse(data []byte) (size int, err error) {
	this.isLongHeader = ((data[0] & 0x80) != 0)
	if this.isLongHeader {
		// Long Header
		if len (data) < 17 {
			err = errors.New ("Header.Parse error: data too small")
			return
		}
		this.packetType = data[size] & 0x7F
		size += 1
		this.connId = binary.BigEndian.Uint64 (data[size:])
		size += 8
		this.version = binary.BigEndian.Uint32 (data[size:])
		size += 4
		if this.packetType != 0 {
			this.packetNumber = binary.BigEndian.Uint32 (data[size:])
			size += 4
		}

		this.keyPhase = false
		this.connIdOmitted = false
	} else {
		// Short Header
		if len (data) < 10 {
			err = errors.New ("Header.Parse error, data too small")
			return
		}
		this.connIdOmitted = (data[0] & 0x40) == 0
		this.keyPhase = (data[0] & 0x020) == 0
		this.packetType = data[0] & 0x1F
		size += 1
		this.connId = binary.BigEndian.Uint64 (data[size:])
		size += 8
		switch this.packetType {
		case 0x1F:
			this.packetNumber = uint32 (data[size])
			size += 1
			break
		case 0x1E:
			this.packetNumber = uint32 (binary.BigEndian.Uint16 (data[size:]))
			size += 2
			break
		case 0x1D:
			this.packetNumber = uint32 (binary.BigEndian.Uint32 (data[size:]))
			size += 4
			break
		}
	}
	return
}

func (this *Header) GetSerializedSize () int {
	if this.isLongHeader {
		return 17
	} else {
		switch this.packetType {
		case 0x1F:
			return 10
		case 0x1E:
			return 11
		case 0x1D:
			return 13
		default:
			return 0
		}
	}
}

func (this *Header) Serialize (data []byte) (size int, err error) {
	size = 0
	if len (data) < this.GetSerializedSize () {
		err = errors.New ("Header.Serialize error: data too small")
		return
	}

	if this.isLongHeader {
		data[size] = 0x80 | this.packetType
		size += 1
		binary.BigEndian.PutUint64 (data[size:], this.connId)
		size += 8
		binary.BigEndian.PutUint32 (data[size:], this.version)
		size += 4
		binary.BigEndian.PutUint32 (data[size:], this.packetNumber)
		size += 4

	} else {
		data[size] = 0x00
		if this.connIdOmitted == false {
			data[size] |= 0x40
		}
		if this.keyPhase == false {
			data[size] |= 0x20
		}
		data[size] |= this.packetType
		size += 1
		binary.BigEndian.PutUint64 (data[size:], this.connId)
		size += 8
		switch this.packetType {
		case 0x1F:
			data[size] = uint8 (this.packetNumber)
			size += 1
			break
		case 0x1E:
			binary.BigEndian.PutUint16 (data[size:], uint16 (this.packetNumber))
			size += 2
			break
		case 0x1D:
			binary.BigEndian.PutUint32 (data[size:], uint32 (this.packetNumber))
			size += 4
			break
		}
	}
	return
}
