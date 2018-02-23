package packet

import "errors"
import "encoding/binary"

type Header struct {
	packetType uint8
	connId uint64
	version uint32
	packetNumber uint32

	keyPhase bool
	connIdOmitted bool
}

func (this *Header) Parse(data []byte) (size int, err error) {
	if (data[0] & 0x80) != 0 {
		// Long Header
		if len (data) < 136 {
			err = errors.New ("Header.Parse error: data too small")
			return
		}
		this.packetType = data[size] & 0x7F
		size += 1
		this.connId = binary.BigEndian.Uint64 (data[size:])
		size += 8
		this.version = binary.BigEndian.Uint32 (data[size:])
		size += 8
		this.packetNumber = binary.BigEndian.Uint32 (data[size:])

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
		size += 4
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
