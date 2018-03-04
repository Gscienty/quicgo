package utils

import (
	"bytes"
	"io"
	"errors"
	"encoding/binary"
)

var BigEndian ByteOrder = bigEndianStruct {}

type bigEndianStruct struct { }

var _ ByteOrder = &bigEndianStruct {}

func (bigEndianStruct) ReadUInt (b io.Reader, len uint8) (uint64, error) {
	buf := make([]byte, len)
	readedLen, readErr := b.Read (buf)
	if readErr != nil {
		return 0, readErr
	}
	if readedLen != int (len) {
		return 0, errors.New ("bigEndianStruct.ReadUInt error: readed length not equal len")
	}
	var retval uint64 = 0
	for i := uint8 (0); i < len; i++ {
		retval |= uint64 (buf[i]) << ((len - 1 - i) * 8)
	}
	return retval, nil
}

func (bigEndianStruct) WriteUInt (b *bytes.Buffer, value uint64, len uint8) {
	buf := make ([]byte, 8)
	binary.BigEndian.PutUint64 (buf, value)
	b.Write (buf[8 - len:])
}