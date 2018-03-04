package utils

import (
	"bytes"
	"testing"
)

func TestLittleEndianWriteUInt (t *testing.T) {
	var b bytes.Buffer
	
	LittleEndian.WriteUInt(&b, 0x010203, 3)

	if (b.Bytes()[0] != 0x03) || (b.Bytes()[1] != 0x02) || (b.Bytes()[2] != 0x01) {
		t.Fail ()
	}
}

func TestLittleEndianReadUInt (t *testing.T) {
	b := bytes.NewReader([] byte { 0x01, 0x02, 0x03, 0x04 })
	v, _ := LittleEndian.ReadUInt(b, 4)
	if v != 0x04030201 {
		t.Log (v)
		t.Fail ()
	}
}