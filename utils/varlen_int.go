package utils

import "io"
import "bytes"
import "encoding/binary"
import "errors"

const (
	MAX_1BYTE_VALUE = 0x3F
	MAX_2BYTE_VALUE = 0x3FFF
	MAX_4BYTE_VALUE = 0x3FFFFFFF
	MAX_8BYTE_VALUE = 0x3FFFFFFFFFFFFFFF
)

type VarLenInteger struct {
	len int
	val uint64
}

func Parse (b io.Reader) (*VarLenInteger, error) {
	byteBuf := make ([]byte, 1)
	_, err := b.Read (byteBuf)
	if err != nil {
		return nil, err
	}

	firstByte := byteBuf[0]
	len := 1 << ((firstByte & 0xC0) >> 6)

	var val uint64
	if len == 1 {
		val = uint64 (firstByte & 0x3F)
	} else if len == 2 {
		_, err := b.Read (byteBuf)
		if err != nil {
			return nil, err
		}
		val = uint64 (byteBuf[0]) + (uint64 (firstByte & 0x3F) << 8)
	} else if len == 4 {
		buf := make ([]byte, 3)
		_, err := b.Read (buf)
		if err != nil {
			return nil, err
		}
		val = (uint64 (binary.BigEndian.Uint16 (buf[1:]))) +
			(uint64 (buf[0]) << 16) +
			(uint64 (firstByte & 0x3F) << 24)
	} else if len == 8 {
		buf := make ([]byte, 7)
		_, err := b.Read (buf)
		if err != nil {
			return nil, err
		}
		val = (uint64 (binary.BigEndian.Uint32 (buf[3:]))) +
			(uint64 (binary.BigEndian.Uint16 (buf[1:])) << 32) +
			(uint64 (buf[0]) << 48) +
			(uint64 (firstByte & 0x3F) << 56)
	} else {
		return nil, errors.New ("VarLenInteger.Parse: len error")
	}

	return &VarLenInteger { len, val }, nil
}

func (this *VarLenInteger) SetVal (val uint64) (err error) {
	if val < MAX_1BYTE_VALUE {
		this.val = val
		this.len = 1
	} else if val < MAX_2BYTE_VALUE {
		this.val = val
		this.len = 2
	} else if val < MAX_4BYTE_VALUE {
		this.val = val
		this.len = 4
	} else if val < MAX_8BYTE_VALUE {
		this.val = val
		this.len = 8
	} else {
		err = errors.New ("VarLenInteger.SetVal: val too large")
	}
	return 
}

func (this *VarLenInteger) Serialize (b bytes.Buffer) (size int, err error) {
	if this.len == 1 {
		b.WriteByte (uint8 (this.val))
		size = 1
	} else if this.len == 2 {
		buf := make ([]byte, 2)
		binary.BigEndian.PutUint16 (buf, uint16 (this.val))
		buf[0] |= 0x40
		b.Write (buf)
		size = 2
	} else if this.len == 4 {
		buf := make ([]byte, 4)
		binary.BigEndian.PutUint32 (buf, uint32 (this.val))
		buf[0] |= 0x80
		b.Write (buf)
		size = 4
	} else if this.len == 8 {
		buf := make ([]byte, 8)
		binary.BigEndian.PutUint64 (buf, this.val)
		buf[0] |= 0xc0
		b.Write (buf)
		size = 8
	} else {
		err = errors.New ("VarLenInteger.Serialize: internal error, .len format error")
	}
	return
}
