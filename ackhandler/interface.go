package ackhandler

import (
	"time"
	"../frame"
	"../protocol"
)

type ISentPacketHandler interface {
	SentPacket(packet *Packet)
	ReceivedAck(ackFrame *frame.AckFrame, withPacketNumber protocol.PacketNumber, recvTime time.Time) error
	SetHandshakeComplete()

	SendingAllowed()
	TimeUntilSend() time.Time
	ShouldSendNumPackets() int

	GetLowestPacketNotConfirmedAcked() protocol.PacketNumber
	DequeuePacketForRetransmission() *Packet
	GetPacketNumberLen(protocol.PacketNumber) uint8

	GetAlarmTimeout() time.Time
	OnAlarm()
}

type IReceivedPacketHandler interface {
	ReceivedPacket(packetNumber protocol.PacketNumber, recvTime time.Time, shouldInstigateAck bool) error
	IgnoreBelow(protocol.PacketNumber)

	GetAlarmTimeout() time.Time
	GetAckFrame() *frame.AckFrame
}