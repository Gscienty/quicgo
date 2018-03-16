package handshake

import (
	"fmt"
	"bytes"
	"testing"
	"../protocol"
)

func TestTransportParametersClientHandler1(t *testing.T) {
	b := bytes.NewBuffer([]byte { })
	s := &clientHelloTransportParameters {
		InitialVersion: protocol.Version(0x01020304),
		Parameters:		[]TransportParameter {
			{ transportParameterID(0x0102), []byte { 1,2,3,4 } },
			{ transportParameterID(0x0103), []byte { 4,5,6,7,8 } },
			{ transportParameterID(0x0104), []byte { 5,6,7,8,9,10 } },
		},
	}

	s.Serialize(b)

	fmt.Println(b.Bytes())
}

func TestTransportParametersClientHandler2(t *testing.T) {
	b := bytes.NewReader([]byte {1, 2, 3, 4, 0, 27, 1, 2, 0, 4, 1, 2, 3, 4, 1, 3, 0, 5, 4, 5, 6, 7, 8, 1, 4, 0, 6, 5, 6, 7, 8, 9, 10})
	
	s := &clientHelloTransportParameters { }

	s.Parse(b)

	fmt.Println(s)
}