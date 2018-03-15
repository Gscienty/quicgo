package handshake

import (
	"fmt"
	"bytes"
	"testing"
	"../protocol"
)

func TestEncryptExtension1(t *testing.T) {
	tester := &encryptedExtension {
		NegotiatedVersion:	protocol.Version(0x01020304),
		SupportedVersions: 	[]protocol.Version {
			protocol.Version(0x01020304), protocol.Version(0x02030405), protocol.Version(0x03040506),
		},
		Parameters:			[]TransportParameter {
			{ 0x0708, []byte { 0x0A, 0x0B, 0x0C, 0x0D } }, { 0x0809, []byte { 0x0E, 0x0F, 0x10, 0x11 } },
		},
	}

	b := bytes.NewBuffer([]byte { })

	tester.Serialize(b)

	fmt.Println(b.Bytes())
}

func TestEncryptExtension2(t *testing.T) {
	tester := []byte { 1, 2, 3, 4, 12, 1, 2, 3, 4, 2, 3, 4, 5, 3, 4, 5, 6, 0, 16, 7, 8, 0, 4, 10, 11, 12, 13, 8, 9, 0, 4, 14, 15, 16, 17 }
	a := &encryptedExtension { }
	a.Parse(bytes.NewReader(tester))
	fmt.Println(a)
}
