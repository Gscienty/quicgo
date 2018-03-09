package utils

import(
	"bytes"
	"testing"
)

func TestBigEndianWriteUInt(t *testing.T) {
	var b bytes.Buffer
	
	BigEndian.WriteUInt(&b, 0x010203, 3)

	if(b.Bytes()[0] != 0x01) ||(b.Bytes()[1] != 0x02) ||(b.Bytes()[2] != 0x03) {
		t.Fail()
	}
}

func TestBigEndianReadUInt(t *testing.T) {
	b := bytes.NewReader([] byte { 0x01, 0x02, 0x03, 0x04 })
	v, _ := BigEndian.ReadUInt(b, 4)
	if v != 0x01020304 {
		t.Fail()
	}
}