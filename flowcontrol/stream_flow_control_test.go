package flowcontrol

import (
	"time"
	"fmt"
	"testing"
	"../protocol"
	"../utils"
)

func streamFlowControlTestBefore() *streamFlowControl {
	rtt := &utils.RTTStat { }
	ctr := &streamFlowControl {
		streamID: *protocol.StreamIDNew(uint64(10)),
		connectionControl: ConnectionFlowControlNew(1000, 1000, rtt).(*connectionFlowControl),
	}

	ctr.maxRecvWindowSize = uint64(10000)
	ctr.rttStat = rtt
	
	return ctr
}

func TestStreamFlowControlConstrctor(t *testing.T) {
	rtt := &utils.RTTStat { }
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
	if fc.sendSize != sendWindow {
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

	ctr.recvCapacity = 1337
	err := ctr.UpdateRecvCapacity(1338, false)
	if err != nil {
		t.Fail()
	}
	if ctr.recvCapacity != uint64(1338) {
		t.Fail()
	}
}
func TestStreamReceivingDataUpdateInformConnectionController(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvCapacity = 10
	ctr.influenceConnection = true
	ctr.connectionControl.recvCapacity = 100
	err := ctr.UpdateRecvCapacity(20, false)
	if err != nil {
		t.Fail()
	}
	if ctr.connectionControl.recvCapacity != uint64(100 + 10) {
		t.Fail()
	}
}

func TestStreamReceivingDataUpdateDoesntInformConnectionController(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvCapacity = 10
	ctr.connectionControl.recvCapacity = 100
	err := ctr.UpdateRecvCapacity(20, false)
	if err != nil {
		t.Fail()
	}
	if ctr.connectionControl.recvCapacity != uint64(100) {
		t.Fail()
	}
}

func TestStreamDoesNotDecreaseHighestReceived(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	ctr.recvCapacity = 1337
	err := ctr.UpdateRecvCapacity(1337, false)
	if err != nil {
		t.Fail()
	}

	err = ctr.UpdateRecvCapacity(1000, false)
	if err != nil {
		t.Fail()
	}
	if ctr.recvCapacity != uint64(1337) {
		t.Fail()
	}
}

func TestStreamDetectsFlowControlViolation(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvCapacity(receiveSize + 1, false)
	if err == nil {
		t.Fail()
	}
	fmt.Println(err.Error())

	ctr.recvCapacity = 100
	err = ctr.UpdateRecvCapacity(101, true)
	if err != nil {
		t.Fail()
	}
	if ctr.recvCapacity != uint64(101) {
		t.Fail()
	}

	ctr.recvCapacity = 100
	err = ctr.UpdateRecvCapacity(99, true)
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

	err := ctr.UpdateRecvCapacity(300, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvCapacity(250, false)
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

	err := ctr.UpdateRecvCapacity(200, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvCapacity(250, false)
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

	err := ctr.UpdateRecvCapacity(200, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvCapacity(200, true)
	if err != nil {
		t.Fail()
	}
	if ctr.recvCapacity != uint64(200) {
		t.Fail()
	}
}

func TestStreamReceivingInconsistentFinalOffset(t *testing.T) {
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(600)
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = receiveSize
	ctr.recvWindowSize = receiveWindowSize

	err := ctr.UpdateRecvCapacity(200, true)
	if err != nil {
		t.Fail()
	}
	err = ctr.UpdateRecvCapacity(201, true)
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
	if ctr.recvedBytesCount != uint64(100) {
		t.Fail()
	}
	if ctr.connectionControl.recvedBytesCount != 0 {
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
	if ctr.recvedBytesCount != uint64(200) {
		t.Fail()
	}
	if ctr.connectionControl.recvedBytesCount != uint64(200) {
		t.Fail()
	}
}

func TestStreamTellWindowUpdate(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
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
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	oldWindowSize := ctr.recvWindowSize
	
	oldoffset := ctr.recvedBytesCount
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
func TestStreamTellConnectionDoesntContribute(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	oldWindowSize := ctr.recvWindowSize
	
	oldoffset := ctr.recvedBytesCount
	ctr.influenceConnection = false
	ctr.rttStat.Update(20 * time.Millisecond, 0, time.Time { })
	ctr.startAutoTuringOffset = oldoffset
	ctr.startAutoTuringTime = time.Now().Add(-time.Millisecond)
	ctr.AddRecvedBytesCount(55)
	offset := ctr.RecvWindowUpdate()
	if offset == 0 {
		t.Fail()
	}
	if ctr.recvWindowSize != 2 * oldWindowSize {
		t.Fail()
	}
	if ctr.connectionControl.recvWindowSize != 2 * oldWindowSize {
		t.Fail()
	}
}
func TestStreamDoesntIncreaseWindowAfterFinalOffset(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	
	ctr.AddRecvedBytesCount(30)
	err := ctr.UpdateRecvCapacity(90, true)
	if err != nil {
		t.Fail()
	}
	if ctr.RecvWindowHasUpdate() == true {
		t.Fail()
	}
	offset := ctr.RecvWindowUpdate()
	if offset != 0 {
		t.Fail()
	}
}
func TestStreamSendingData1(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	
	ctr.SetSendSize(15)
	ctr.AddSendedBytesCount(5)
	if ctr.GetSendWindowSize() != uint64(10) {
		t.Fail()
	}
}
func TestStreamSendingData2(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	
	ctr.SetSendSize(15)
	ctr.connectionControl.SetSendSize(1)
	ctr.AddSendedBytesCount(5)
	if ctr.GetSendWindowSize() != uint64(10) {
		t.Fail()
	}
}
func TestStreamSendingData3(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	
	ctr.influenceConnection = true
	ctr.connectionControl.SetSendSize(12)
	ctr.SetSendSize(20)
	ctr.AddSendedBytesCount(10)
	if ctr.GetSendWindowSize() != uint64(2) {
		fmt.Printf("%d\n", ctr.GetSendWindowSize())
		t.Fail()
	}
}
func TestStreamSendingData4(t *testing.T) {
	ctr := streamFlowControlTestBefore()
	ctr.recvSize = 100
	ctr.recvWindowSize = 60
	ctr.recvedBytesCount = 100 - 60
	ctr.connectionControl.recvWindowSize = 120
	
	ctr.influenceConnection = true
	ctr.connectionControl.SetSendSize(50)
	ctr.SetSendSize(100)
	ctr.AddSendedBytesCount(50)
	blocked, _ := ctr.connectionControl.IsNewlyBlocked()
	if blocked == false {
		t.Fail()
	}
	blocked, _ = ctr.IsBlocked()
	if blocked == false {
		t.Fail()
	}
}


