package ackhandler

import (
	"time"
	"../protocol"
	"../frame"
)

type Packet struct {
	PacketNumber	protocol.PacketNumber
	PacketType		protocol.PacketType
	Frames			[]frame.IFrame
	Length			uint64

	largestAcked	protocol.PacketNumber
	sendTime		time.Time
}

func (this *Packet) GetFramesForRetransmission() []frame.IFrame {
	var fs []frame.IFrame
	for _, f := range this.Frames {
		if f.GetType() == frame.FRAME_TYPE_ACK {
			continue
		}
		fs = append(fs, f)
	}

	return fs
}