package flowcontrol

import (
	"errors"
	"../protocol"
	"../utils"
)

type IStreamFlowControl interface {
	IFlowControl

	IsBlocked() (bool, uint64)
	UpdateRecvCapacity(uint64, bool) error
	RecvWindowHasUpdate() bool
}

type streamFlowControl struct {
	FlowControl
	streamID			protocol.StreamID
	connectionControl	*connectionFlowControl
	influenceConnection	bool
	recvFinalOffset		bool
}

var _ IStreamFlowControl = &streamFlowControl { }

func StreamFlowControlNew(
	streamID			protocol.StreamID,
	influenceConnection	bool,
	connectionControl	*connectionFlowControl,
	recvSize			uint64,
	maxRecvWindowSize	uint64,
	initialSendSize		uint64,
	rttStat				*utils.RTTStat,
) IStreamFlowControl {
	return &streamFlowControl {
		streamID:				streamID,
		connectionControl:		connectionControl,
		influenceConnection:	influenceConnection,
		FlowControl:			FlowControl {
			rttStat:			rttStat,
			recvSize:			recvSize,
			recvWindowSize:		recvSize,
			maxRecvWindowSize:	maxRecvWindowSize,
			sendSize:			initialSendSize,
		},
	}
}

func (this *streamFlowControl) IsBlocked() (bool, uint64) {
	if this.GetSendWindowSize() != 0 {
		return false, 0
	}

	return true, this.sendSize
}

func (this *streamFlowControl) UpdateRecvCapacity(offset uint64, final bool) error {
	this.recvRWLock.Lock()
	defer this.recvRWLock.Unlock()

	if final && this.recvFinalOffset && offset != this.recvCapacity {
		return errors.New("final error")
	}
	if this.recvFinalOffset && offset > this.recvCapacity {
		return errors.New("overflow")
	}
	this.recvFinalOffset = final
	if offset == this.recvCapacity {
		return nil
	}
	if offset < this.recvCapacity {
		if final {
			return errors.New("termination early")
		}
		return nil
	}
	increment := offset - this.recvCapacity
	this.recvCapacity = offset
	if this.recvCapacity > this.recvSize {
		return errors.New("recevied too much data")
	}
	if this.influenceConnection {
		return this.connectionControl.AddRecvCapacity(increment)
	}
	return nil
}

func (this *streamFlowControl) RecvWindowHasUpdate() bool {
	this.recvRWLock.Lock()
	hasWindowUpdate := !this.recvFinalOffset && this.recvWindowHasUpdate()
	this.recvRWLock.Unlock()
	return hasWindowUpdate
}

func (this *streamFlowControl) RecvWindowUpdate() uint64 {
	this.recvRWLock.Lock()

	if this.recvFinalOffset {
		this.recvRWLock.Unlock()
		return 0
	}

	originWindowSize := this.recvWindowSize
	offset := this.FlowControl.recvWindowUpdate()
	if this.recvWindowSize > originWindowSize {
		if this.influenceConnection {
			this.connectionControl.EnsureRecvMinimumWindowSize(uint64(float64(this.recvWindowSize) * protocol.CONNECTION_FLOW_CONTROL_MULTIPLIER))
		}
	}

	this.recvRWLock.Unlock()
	return offset
}

func (this *streamFlowControl) AddSendedBytesCount(n uint64) {
	this.FlowControl.AddSendedBytesCount(n)
	if this.influenceConnection {
		this.connectionControl.AddSendedBytesCount(n)
	}
}

func (this *streamFlowControl) AddRecvedBytesCount(n uint64) {
	this.FlowControl.AddRecvedBytesCount(n)
	if this.influenceConnection {
		this.connectionControl.AddRecvedBytesCount(n)
	}
}

func (this *streamFlowControl) GetSendWindowSize() uint64 {
	windowSize := this.FlowControl.GetSendWindowSize()
	if this.influenceConnection && this.connectionControl.GetSendWindowSize() < windowSize {
		windowSize = this.connectionControl.GetSendWindowSize()
	}
	return windowSize
}