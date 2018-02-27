package packet

import "errors"
import "encoding/binary"

type VariableLengthIntegerEncoder struct {
	val uint64
	usableBits int
}

func (this *VariableLengthIntegerEncoder) Parse (data []byte) (size int, err error) {
	this.usableBits = int((data[0] & 0xC0) >> 6)
	dataLen := len (data)
	switch this.usableBits {
	case 0:
		if dataLen < 1 {
			err = errors.New ("VariableLengthIntegerEncoder.Parse: data too small")
			return
		}
		this.val = uint64 (data[0] & 0x3F)
		size = 1
		break
	case 1:
		if dataLen < 2 {
			err = errors.New ("VariableLengthIntegerEncoder.Parse: data too small")
			return
		}
		this.val = uint64 (binary.BigEndian.Uint16(data) & 0x3FFF)
		size = 2
		break
	case 2:
		if dataLen < 4 {
			err = errors.New ("VariableLengthIntegerEncoder.Parse: data too small")
			return
		}
		this.val = uint64 (binary.BigEndian.Uint32(data) & 0x3FFFFFFF)
		size = 4
		break
	case 3:
		if dataLen < 8 {
			err = errors.New ("VariableLengthIntegerEncoder.Parse: data too small")
			return
		}
		this.val = (binary.BigEndian.Uint64(data) & 0x3FFFFFFFFFFFFFFF)
		size = 8
		break
	}
	return
}

func (this *VariableLengthIntegerEncoder) GetSerializedSize () int {
	switch this.usableBits {
	case 0:
		return 1
	case 1:
		return 2
	case 2:
		return 4
	case 4:
		return 8
	default:
		return 0
	}
}

func (this *VariableLengthIntegerEncoder) Serialize (data []byte) (size int, err error) {
	size = 0
	if len (data) < this.GetSerializedSize () {
		err = errors.New ("VariableLengthIntegerEncoder.Serialize error: data too small")
		return
	}
	switch this.usableBits {
	case 0:
		data[size] = uint8 (this.val) | 0x00
		break
	case 1:
		binary.BigEndian.PutUint16 (data, uint16 (this.val))
		data[0] |= 0x40
		size += 1
		break
	case 2:
		binary.BigEndian.PutUint32 (data, uint32 (this.val))
		data[0] |= 0x80
		size += 2
		break
	case 3:
		binary.BigEndian.PutUint64 (data, this.val)
		data[0] |= 0xC0
		size += 4
		break
	default:
		err = errors.New ("VariableLengthIntegerEncoder.Serialize error: internal error, have not legal usableBits")
	}
	return
}

func (this *VariableLengthIntegerEncoder) SetVal (val uint64) (err error) {
	if val < 0x40 {
		this.val = val
		this.usableBits = 0
	} else if val < 0x4000 {
		this.val = val
		this.usableBits = 1
	} else if val < 0x40000000 {
		this.val = val
		this.usableBits = 2
	} else if val <= 0x4000000000000000 {
		this.val = val
		this.usableBits = 3
	} else {
		err = errors.New ("VariableLengthIntegerEncoder.SetVal error: val too large")
	}
	return
}