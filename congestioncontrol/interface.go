package congestioncontrol

import (
	"time"
)

type SendAlgo interface {
	TimeUntilSend(inSending uint64) time.Duration
	PacketSent(sentTime time.Time, inSending uint64, packetNumber uint64, bytesCount uint64) bool
	GetCongestionWindow() uint64
	TryExitSlowStart()
	Acked(n uint64, ackedBytes uint64, inSending uint64)
	Lost(n uint64, lostBytes uint64, inSending uint64)
	SetConnectionNumber(n int)
	RetransmissionTimeout(packetsRetransmitted bool)
	ConnectionMigration()
	RetransmissionDelay() time.Duration

	SetSlowStartLargeReduction(enable bool)
}

type SendAlgoDebug interface {
	SendAlgo

	BandwidthEstimate() uint64

	GetSlowStart() *SlowStart
	SlowStartThreshold() uint64
	RenoBeta() float32
	InRecovery() bool
}