package utils

import (
	"bytes"
	"io"
)

type ByteOrder interface {
	ReadUInt (io.Reader, uint8) (uint64, error)
	WriteUInt (*bytes.Buffer, uint64, uint8)
}