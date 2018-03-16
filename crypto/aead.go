package crypto

import (
	"../protocol"
)

type AEAD interface {
	Open(dst []byte, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) ([]byte, error)
	Seal(dst []byte, src []byte, packetNumber protocol.PacketNumber, associatedData []byte) []byte
	Overhead() int
}