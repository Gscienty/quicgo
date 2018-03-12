package flowcontrol

import(
	"../utils"
)

type IConnectionFlowControl interface {
	IFlowControl

	IsNewlyBlocked() (bool, uint64)
}

type connectionFlowControl struct {
	FlowControl
	lastBlockedAt	uint64
}

var _ IConnectionFlowControl = &connectionFlowControl { }

func ConnectionFlowControlNew(recvWindow uint64, maxRecvWindow uint64, rttStat *utils.RTTStat) IConnectionFlowControl {
	return &connectionFlowControl {
		FlowControl: FlowControl {
			rttStat:			rttStat,
			recvSize:			recvWindow,
			recvWindowSize:		recvWindow,
			maxRecvWindowSize:	maxRecvWindow,
		},
	}
}

func (this *connectionFlowControl) IsNewlyBlocked() (bool, uint64) {
	if this.FlowControl.GetSendWindowSize() != 0 || this.sendSize == this.lastBlockedAt {
		return false, 0
	}

	this.lastBlockedAt = this.sendSize
	return true, this.sendSize
}

func (this *connectionFlowControl) AddHighestOffset(increment uint64) error {
	this.recvRWLock.Lock()
	defer this.recvRWLock.Unlock()

	this.recvHighestOffset += increment

	return nil
}

func (this *connectionFlowControl) EnsureRecvMinimumWindowSize(size uint64) {
	this.recvRWLock.Lock()
	if size > this.recvWindowSize {
		if size > this.maxRecvWindowSize {
			this.recvWindowSize = this.maxRecvWindowSize
		} else {
			this.recvWindowSize = size
		}

		this.startAutoTuring()
	}
	this.recvRWLock.Unlock()
}

func (this *connectionFlowControl) RecvWindowUpdate() uint64 {
	this.recvRWLock.Lock()
	offset := this.FlowControl.recvWindowUpdate()
	this.recvRWLock.Unlock()
	return offset
}