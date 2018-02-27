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