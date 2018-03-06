package protocol

import (
	"testing"
	"bytes"
	"fmt"
)

func TestStreamIDParse1 (t *testing.T) {
	// 1110
	b := bytes.NewReader ([]byte { 0x3A })

	streamID, err := StreamIDParse (b)
	if err != nil {
		fmt.Println (err.Error ())
		t.FailNow ()
	}

	if streamID.GetID () != 0xE {
		fmt.Printf ("%x\n", streamID.GetID ())
		t.FailNow ()
	}
}

func TestStreamIDParse2 (t *testing.T) {
	// 1110 1010 1011
	b := bytes.NewReader ([]byte { 0x7A, 0xAC })

	streamID, err := StreamIDParse (b)
	if err != nil {
		fmt.Println (err.Error ())
		t.FailNow ()
	}

	if streamID.GetID () != 0xEAB {
		fmt.Printf ("%x\n", streamID.GetID ())
		t.FailNow ()
	}
}

func TestStreamIDParse3 (t *testing.T) {
	// 1110 1010 1011 0001 0000 1000 0101
	b := bytes.NewReader ([]byte { 0xBA, 0xAC, 0x42, 0x15 })

	streamID, err := StreamIDParse (b)
	if err != nil {
		fmt.Println (err.Error ())
		t.FailNow ()
	}

	if streamID.GetID () != 0xEAB1085 {
		fmt.Printf ("%x\n", streamID.GetID ())
		t.FailNow ()
	}
}

func TestStreamIDParse4 (t *testing.T) {
	// 1110 1010 1011 0001 0000 1000 0101 0101 1010 0001 1111 0101 0111 0010 0100
	b := bytes.NewReader ([]byte { 0xFA, 0xAC, 0x42, 0x15, 0x68, 0x7D, 0x5C, 0x93 })

	streamID, err := StreamIDParse (b)
	if err != nil {
		fmt.Println (err.Error ())
		t.FailNow ()
	}

	if streamID.GetID () != 0xEAB10855A1F5724 {
		fmt.Printf ("%x\n", streamID.GetID ())
		t.FailNow ()
	}
}
