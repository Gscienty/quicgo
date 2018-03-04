package utils

import (
	"bytes"
	"io"
	"errors"
	"encoding/binary"
)

var LittleEndian ByteOrder = littleEndianStruct { }

type littleEndianStruct struct { }

var _ ByteOrder = &littleEndianStruct { }

func (littleEndianStruct) ReadUInt (b io.Reader, len uint8) (uint64, error) {
	buf := make ([]byte, len)
	readedLen, readErr := b.Read (buf)
	if readErr != nil {
		return 0, readErr
	}

	if readedLen != int (len) {
		return 0, errors.New ("littleEndianStruct.ReadUInt: readed length not equal len")
	}
	var retval uint64 = 0
	for i := uint8 (0); i < len; i++ {
		retval |= uint64 (buf[i]) << (i << 3)
	}
	return retval, nil
}

func (littleEndianStruct) WriteUInt (b *bytes.Buffer, value uint64, len uint8) {
	buf := make ([]byte, 8)
	binary.LittleEndian.PutUint64 (buf, value)
	b.Write (buf[:len])
}