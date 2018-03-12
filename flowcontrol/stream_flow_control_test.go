package flowcontrol

import (
	"time"
	"fmt"
	"testing"
	"../protocol"
)

func streamFlowControlTestBefore() *streamFlowControl {
	rtt := &RTTStat { }
	ctr := &streamFlowControl {
		streamID: *protocol.StreamIDNew(uint64(10)),
		connectionControl: ConnectionFlowControlNew(1000, 1000, rtt).(*connectionFlowControl),
	}

	ctr.maxRecvWindowSize = uint64(10000)
	ctr.rttStat = rtt
	
	return ctr
}

func TestStreamFlowControlConstrctor(t *testing.T) {
	rtt := &RTTStat { }
	receiveSize := uint64(2000)
	maxReceiveWindow := uint64(3000)
	sendWindow := uint64(4000)

	cc := ConnectionFlowControlNew(0, 0, nil).(*connectionFlowControl)
	fc := StreamFlowControlNew(*protocol.StreamIDNew(5), true, cc, receiveSize, maxReceiveWindow, sendWindow, rtt).(*streamFlowControl)

	if fc.streamID.Equal(protocol.StreamIDNew(5)) == false {
		t.Fail()
	}
	if fc.recvSize != receiveSize {
		t.Fail()
	}
	if fc.maxRecvWindowSize != maxReceiveWindow {
		t.Fail()
	}
	if fc.sendOffset != sendWindow {
		t.Fail()
	}
	if fc.influenceConnection == false {
		t.Fail()
	}
}
func TestStreamReceivingDataUpdateHighestReceived(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)

	ctr := streamFlowControlTestBefore()

	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvHighestOffset = 1337
	err := ctr.UpdateRecvHighestOffset(1338, false)
	if err != nil {
		t.Fail()
	}
	if ctr.recvHighestOffset != uint64(1338) {
		t.Fail()
	}
}
func TestStreamReceivingDataUpdateInformConnectionController(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvHighestOffset = 10
	ctr.influenceConnection = true
	ctr.connectionControl.recvHighestOffset = 100
	err := ctr.UpdateRecvHighestOffset(20, false)
	if err != nil {
		t.Fail()
	}
	if ctr.connectionControl.recvHighestOffset != uint64(100 + 10) {
		t.Fail()
	}
}

func TestStreamReceivingDataUpdateDoesntInformConnectionController(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvHighestOffset = 10
	ctr.connectionControl.recvHighestOffset = 100
	err := ctr.UpdateRecvHighestOffset(20, false)
	if err != nil {
		t.Fail()
	}
	if ctr.connectionControl.recvHighestOffset != uint64(100) {
		t.Fail()
	}
}

func TestStreamDoesNotDecreaseHighestReceived(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvHighestOffset = 1337
	err := ctr.UpdateRecvHighestOffset(1337, false)
	if err != nil {
		t.Fail()
	}

	err = ctr.UpdateRecvHighestOffset(1000, false)
	if err != nil {
		t.Fail()
	}
	if ctr.recvHighestOffset != uint64(1337) {
		t.Fail()
	}
}

func TestStreamDetectsFlowControlViolation(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvHighestOffset(receiveSize + 1, false)
	if err == nil {
		t.Fail()
	}
	fmt.Println(err.Error())

	ctr.recvHighestOffset = 100
	err = ctr.UpdateRecvHighestOffset(101, true)
	if err != nil {
		t.Fail()
	}
	if ctr.recvHighestOffset != uint64(101) {
		t.Fail()
	}

	ctr.recvHighestOffset = 100
	err = ctr.UpdateRecvHighestOffset(99, true)
	if err == nil {
		t.Fail()
	}
	fmt.Println(err.Error())
}

func TestStreamReceivingAFinalOffset(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvHighestOffset(300, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvHighestOffset(250, false)
	if err != nil {
		t.Fail()
	}
}

func TestStreamReceivingAHigherOffsetAfterFinalOffset(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvHighestOffset(200, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvHighestOffset(250, false)
	if err == nil {
		t.Fail()
	}
	fmt.Println(err.Error())
}

func TestStreamAcceptsDuplicateFinalOffset(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvHighestOffset(200, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvHighestOffset(200, true)
	if err != nil {
		t.Fail()
	}
	if ctr.recvHighestOffset != uint64(200) {
		t.Fail()
	}
}

func TestStreamReceivingInconsistentFinalOffset(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvHighestOffset(200, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvHighestOffset(201, true)
	if err == nil {
		t.Fail()
	}
	fmt.Println(err.Error())
}

func TestStreamSaveReadOnStreamNotContributingConnection(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.AddRecvedBytesCount(100)
	if ctr.recvBytesCount != uint64(100) {
		t.Fail()
	}
	if ctr.connectionControl.recvBytesCount != 0 {
		t.Fail()
	}
}

func TestStreamSaveReadOnStreamNotContributingConnection2(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.influenceConnection = true
	ctr.AddRecvedBytesCount(200)
	if ctr.recvBytesCount != uint64(200) {
		t.Fail()
	}
	if ctr.connectionControl.recvBytesCount != uint64(200) {
		t.Fail()
	}
}

func TestStreamTellWindowUpdate(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	
	if ctr.RecvWindowHasUpdate() {
		t.Fail()
	}
	ctr.AddRecvedBytesCount(30)
	if ctr.RecvWindowHasUpdate() == false {
		t.Fail()
	}
	if ctr.RecvWindowUpdate() == 0 {
		t.Fail()
	}
	if ctr.RecvWindowHasUpdate() {
		t.Fail()
	}
}
func TestStreamTellConnectionFlowControlAutotuned(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	oldWindowSize := ctr.recvWindowSize
	
	oldoffset := ctr.recvBytesCount
	ctr.influenceConnection = true
	ctr.rttStat.Update(20 * time.Millisecond, 0, time.Time { })
	ctr.startAutoTuringOffset = oldoffset
	ctr.startAutoTuringTime = time.Now().Add(-time.Millisecond)
	ctr.AddRecvedBytesCount(55)
	offset := ctr.RecvWindowUpdate()
	if offset != oldoffset + 55 + 2 * oldWindowSize {
		t.Fail()
	}
	if ctr.recvWindowSize != 2 * oldWindowSize {
		t.Fail()
	}
	if ctr.connectionControl.recvWindowSize != uint64(float64(ctr.recvWindowSize) * protocol.CONNECTION_FLOW_CONTROL_MULTIPLIER) {
		fmt.Printf("%d %d\n", ctr.connectionControl.recvWindowSize, uint64(float64(ctr.recvWindowSize) * protocol.CONNECTION_FLOW_CONTROL_MULTIPLIER))
		t.Fail()
	}
}


