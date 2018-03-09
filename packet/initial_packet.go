package packet

import(
	"../protocol"
	"../frame"
)

type InitialPacket struct {
	header	protocol.Header
	payload	[]frame.Frame
}

