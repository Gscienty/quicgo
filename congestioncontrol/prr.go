package congestioncontrol

import (
	"math"
	"time"
	"../protocol"
)

type PRR struct {
	sentBytesCountSinceLoss			uint64
	deliveredBytesCountSinceLoss	uint64
	ackCountSinceLoss				uint64
	inSendingBytesCountBeforeLoss	uint64
}

func (this *PRR) SendCount(n uint64) {
	this.sentBytesCountSinceLoss += n
}

func (this *PRR) Lost(inSendingCount uint64) {
	this.sentBytesCountSinceLoss = 0
	this.deliveredBytesCountSinceLoss = 0
	this.ackCountSinceLoss = 0
	this.inSendingBytesCountBeforeLoss = inSendingCount
}

func (this *PRR) Acked(ackedBytes uint64) {
	this.ackCountSinceLoss += ackedBytes
	this.ackCountSinceLoss++
}

func (this *PRR) TimeUntilSend (congestionWindow uint64, inSendingBytesCount uint64, threhold uint64) time.Duration {
	if this.sentBytesCountSinceLoss == 0 || inSendingBytesCount < protocol.DEFAULT_TCP_MSS {
		return 0
	}
	if congestionWindow > inSendingBytesCount {
		if this.deliveredBytesCountSinceLoss + this.ackCountSinceLoss * protocol.DEFAULT_TCP_MSS <= this.sentBytesCountSinceLoss {
			return time.Duration(math.MaxInt64)
		}
		return 0
	}
	if this.deliveredBytesCountSinceLoss * threhold > this.sentBytesCountSinceLoss * this.inSendingBytesCountBeforeLoss {
		return 0
	}
	return time.Duration(math.MaxInt64)
}