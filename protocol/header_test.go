package protocol

import (
	"testing"
	"bytes"
	"fmt"
)

func TestHeaderParse (t *testing.T) {
	b := bytes.NewReader ([]byte {
		0x80 | 0x7F,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C,
		0x10, 0x20, 0x30, 0x40,
	})

	header, err := HeaderParse (b)
	if err != nil {
		t.FailNow ()
	}
	if header.isLongHeader == false {
		t.Fail ()
	}
	if header.connectionID != ConnectionID (0x0102030405060708) {
		t.Fail ()
		fmt.Printf ("connectionID %x", header.connectionID)
	}
	if header.version != Version (0x090A0B0C) {
		t.Fail ()
	}
	if header.packetNumber != PacketNumber (0x10203040) {
		t.Fail ()
	}
}

func TestHeaderSerializedLength (t *testing.T) {
	ori := []byte {
		0x80 | 0x7F,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C,
		0x10, 0x20, 0x30, 0x40,
	}

	b := bytes.NewReader (ori)
	header, _ := HeaderParse (b)

	if header.SerializedLength () != uint8 (len (ori)) {
		t.Fail ()
	}
}

func TestHeaderSerialize (t *testing.T) {
	ori := []byte {
		0x80 | 0x7F,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C,
		0x10, 0x20, 0x30, 0x40,
	}
	b := bytes.NewReader (ori)

	header, _ := HeaderParse (b)
	var ret bytes.Buffer
	err := header.Serialize (&ret)
	if err != nil {
		t.FailNow ()
	}

	if len (ori) != len (ret.Bytes ()) {
		t.FailNow ()
	}

	for i, v := range ret.Bytes () {
		if ori[i] != v {
			t.FailNow ()
		}
	}

}