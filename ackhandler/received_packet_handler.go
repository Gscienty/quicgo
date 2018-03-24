package ackhandler

import (
	"time"
	"../protocol"
)

type receivedPacketHandler struct {
	largestObserved				protocol.PacketNumber
	ignoreBelow					protocol.PacketNumber
	largestObservedReceivedTime	time.Time

	
}