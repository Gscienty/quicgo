package flowcontrol

import(
	"sync"
	"time"
	"../protocol"
)


type IFlowControl interface {
	GetSendWindowSize() uint64
	SetSendWindowOffset(uint64)

	AddSendedBytesCount(uint64)
	AddRecvedBytesCount(uint64)

	RecvWindowUpdate() uint64
}

type FlowControl struct {
	sendedBytesCount		uint64
	sendOffset				uint64

	recvRWLock				sync.RWMutex

	recvOffset				uint64
	recvSize				uint64
	recvWindowSize			uint64
	maxRecvWindowSize		uint64
	recvBytesCount			uint64
	recvHighestOffset		uint64

	startAutoTuringTime		time.Time
	startAutoTuringOffset	uint64
	rttStat					*RTTStat
}

func (this *FlowControl) GetSendWindowSize() uint64 {
	if this.sendedBytesCount > this.sendOffset {
		return 0
	}
	return this.sendOffset - this.sendedBytesCount
}

func (this *FlowControl) SetSendWindowOffset(offset uint64) {
	if offset > this.sendOffset {
		this.sendOffset = offset
	}
}

func (this *FlowControl) AddSendedBytesCount(n uint64) {
	this.sendedBytesCount += n
}

func (this *FlowControl) AddRecvedBytesCount(n uint64) {
	this.recvRWLock.Lock()
	defer this.recvRWLock.Unlock()

	if this.recvBytesCount == 0 {
		this.startAutoTuring()
	}
	this.recvBytesCount += n
}

func (this *FlowControl) adjustWindowSize() {
	bytesReadInDuringCount := this.recvBytesCount - this.startAutoTuringOffset

	if bytesReadInDuringCount < this.recvWindowSize / 2 {
		return
	}

	rtt := this.rttStat.smoothedRTT
	if rtt == 0 {
		return
	}

	fraction := float64(bytesReadInDuringCount) / float64(this.recvBytesCount)
	if time.Since(this.startAutoTuringTime) < time.Duration(4 * fraction * float64(rtt)) {
		if(this.recvWindowSize << 1) < this.maxRecvWindowSize {
			this.recvWindowSize <<= 1
		} else {
			this.recvWindowSize = this.maxRecvWindowSize
		}
	}

	this.startAutoTuring()
}

func (this *FlowControl) startAutoTuring() {
	this.startAutoTuringTime = time.Now()
	this.startAutoTuringOffset = this.recvBytesCount
}

func (this *FlowControl) recvWindowHasUpdate() bool {
	remain := this.recvSize - this.recvBytesCount
	return remain <= uint64(float64(this.recvWindowSize) * float64(1 - protocol.RECV_WINDOW_UPDATE_THREHOLD))
}

func (this *FlowControl) recvWindowUpdate() uint64 {
	if this.recvWindowHasUpdate() {
		this.adjustWindowSize()
		this.recvSize = this.recvBytesCount + this.recvWindowSize
	}

	return this.recvSize
}
