package flowcontrol

import(
	"sync"
	"time"
	"../protocol"
	"../utils"
)


type IFlowControl interface {
	GetSendWindowSize() uint64
	SetSendSize(uint64)

	AddSendedBytesCount(uint64)
	AddRecvedBytesCount(uint64)

	RecvWindowUpdate() uint64
}

type FlowControl struct {
	sendedBytesCount		uint64
	sendSize				uint64

	recvRWLock				sync.RWMutex

	recvSize				uint64
	recvedBytesCount		uint64
	recvWindowSize			uint64
	maxRecvWindowSize		uint64
	recvCapacity			uint64

	startAutoTuringTime		time.Time
	startAutoTuringOffset	uint64
	rttStat					*utils.RTTStat
}

func (this *FlowControl) GetSendWindowSize() uint64 {
	if this.sendedBytesCount > this.sendSize {
		return 0
	}
	return this.sendSize - this.sendedBytesCount
}

func (this *FlowControl) SetSendSize(size uint64) {
	if size > this.sendSize {
		this.sendSize = size
	}
}

func (this *FlowControl) AddSendedBytesCount(n uint64) {
	this.sendedBytesCount += n
}

func (this *FlowControl) AddRecvedBytesCount(n uint64) {
	this.recvRWLock.Lock()
	defer this.recvRWLock.Unlock()

	if this.recvedBytesCount == 0 {
		this.startAutoTuring()
	}
	this.recvedBytesCount += n
}

func (this *FlowControl) adjustRecvWindowSize() {
	bytesReadInDuringCount := this.recvedBytesCount - this.startAutoTuringOffset
	if bytesReadInDuringCount < this.recvWindowSize / 2 {
		return
	}

	rtt := this.rttStat.GetSmoothedRTT()
	if rtt == 0 {
		return
	}

	fraction := float64(bytesReadInDuringCount) / float64(this.recvWindowSize)
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
	this.startAutoTuringOffset = this.recvedBytesCount
}

func (this *FlowControl) recvWindowHasUpdate() bool {
	remain := this.recvSize - this.recvedBytesCount
	return remain <= uint64(float64(this.recvWindowSize) * float64(1 - protocol.RECV_WINDOW_UPDATE_THREHOLD))
}

func (this *FlowControl) recvWindowUpdate() uint64 {
	if this.recvWindowHasUpdate() == false {
		return 0
	}

	this.adjustRecvWindowSize()
	this.recvSize = this.recvedBytesCount + this.recvWindowSize
	return this.recvSize
}
