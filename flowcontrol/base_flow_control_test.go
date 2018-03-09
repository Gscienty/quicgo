package flowcontrol

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"testing"
	"../protocol"
)

func TestSendAddBytesSent(t *testing.T) {
	ctr := &FlowControl { }

	ctr.sendedBytesCount = 5
	ctr.AddSendedBytesCount(6)
	if ctr.sendedBytesCount != 5 + 6 {
		t.Fail()
	}
}

func TestSendRemainWindow(t *testing.T) {
	ctr := &FlowControl { }
	ctr.sendedBytesCount = 5
	ctr.sendOffset = 12
	if ctr.GetSendWindowSize() != 12 - 5 {
		t.Fail()
	}
}

func TestSendUpdateWindow(t *testing.T) {
	ctr := &FlowControl { }
	ctr.AddSendedBytesCount(5)
	ctr.SetSendOffset(15)
	if ctr.sendOffset != 15 {
		t.Fail()
	}
	if ctr.GetSendWindowSize() != 15 - 5 {
		t.Fail()
	}
}

func TestSendZeroWindow(t *testing.T) {
	ctr := &FlowControl { }
	ctr.AddSendedBytesCount(15)
	ctr.SetSendOffset(10)
	if ctr.GetSendWindowSize() != 0 {
		t.Fail()
	}
}

func TestSendDecreaseWindow(t *testing.T) {
	ctr := &FlowControl { }
	ctr.SetSendOffset(20)
	if ctr.GetSendWindowSize() != 20 {
		t.Fail()
	}
	ctr.SetSendOffset(10)
	if ctr.GetSendWindowSize() != 20 {
		t.Fail()
	}
}

func TestRecvAddRecv(t *testing.T) {
	ctr := &FlowControl { }
	ctr.recvBytesCount = 10000 - 1000
	ctr.recvWindowSize = 1000
	ctr.recvSize = 10000

	ctr.recvBytesCount = 5
	ctr.AddRecvedBytesCount(6)
	if ctr.recvBytesCount != 5 + 6 {
		t.Fail()
	}
}

func TestRecvTrigger(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	byteConsumed := float64(receiveWindowSize) * protocol.RECV_WINDOW_UPDATE_THREHOLD + 1
	byteRemaining := receiveWindowSize - uint64(byteConsumed)
	readPosition := receiveLength - byteRemaining
	ctr.recvBytesCount = readPosition
	offset := ctr.recvWindowUpdate()
	if offset != readPosition + 1000 {
		t.Fail()
	}
	if ctr.recvSize != readPosition + 1000 {
		t.Fail()
	}
}

func TestRectNonTrigger(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	byteConsumed := float64(receiveWindowSize) * protocol.RECV_WINDOW_UPDATE_THREHOLD - 1
	byteRemaining := receiveWindowSize - uint64(byteConsumed)
	readPosition := receiveLength - byteRemaining
	ctr.recvBytesCount = readPosition
	offset := ctr.recvWindowUpdate()
	if offset != 0 {
		t.Fail()
	}
}

func TestRectAutoTuring(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	oldWindowSize := ctr.recvWindowSize
	ctr.maxRecvWindowSize = 5000

	ctr.adjustWindowSize()
	if ctr.recvWindowSize != oldWindowSize {
		t.Fail()
	}
}


func setRTT(ctr *FlowControl, d time.Duration, t *testing.T) {
	ctr.rttStat.Update(d, 0, time.Now())
	if (ctr.rttStat.smoothedRTT != d) {
		t.Fail()
	}
}

func TestRectAutoTuring2(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	oldWindowSize := ctr.recvWindowSize
	ctr.maxRecvWindowSize = 5000

	setRTT(ctr, 0, t)
	ctr.startAutoTuring()
	ctr.AddRecvedBytesCount(400)
	offset := ctr.recvWindowUpdate()
	if offset == 0 {
		t.Fail()
	}
	if ctr.recvWindowSize != oldWindowSize {
		t.Fail()
	}
}

func scaleDuration(t time.Duration) time.Duration {
	scaleFactor := 1
	if f, err := strconv.Atoi(os.Getenv("TIMESCALE_FACTOR")); err == nil {
		scaleFactor = f
	}
	return time.Duration(scaleFactor) * t
}

func TestRectAutoTuring3(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	oldWindowSize := ctr.recvWindowSize
	ctr.maxRecvWindowSize = 5000

	bytesread := ctr.recvBytesCount
	rtt := scaleDuration(20 * time.Millisecond)
	setRTT(ctr, rtt, t)
	dataRead := receiveWindowSize * 2 / 3 + 1
	ctr.startAutoTuringOffset = ctr.recvBytesCount
	ctr.startAutoTuringTime = time.Now().Add(-rtt * 4 * 2 / 3)
	ctr.AddRecvedBytesCount(dataRead)
	offset := ctr.recvWindowUpdate()
	if offset == 0 {
		fmt.Println("offset == 0 error")
		t.Fail()
	}
	newWindowSize := ctr.recvWindowSize
	if newWindowSize != 2 * oldWindowSize {
		fmt.Printf("newWindowSize != 2 * oldWIndowSize [%d] - [%d] \n", newWindowSize, oldWindowSize)
		t.Fail()
	}
	if offset != bytesread + dataRead + newWindowSize {
		fmt.Println("offset != bytesread + dataRead + newWindowSize")
		t.Fail()
	}
}

func TestRectAutoTuring4(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	oldWindowSize := ctr.recvWindowSize
	ctr.maxRecvWindowSize = 5000

	bytesread := ctr.recvBytesCount
	rtt := scaleDuration(20 * time.Millisecond)
	setRTT(ctr, rtt, t)
	dataRead := receiveWindowSize * 1 / 3 + 1
	ctr.startAutoTuringOffset = ctr.recvBytesCount
	ctr.startAutoTuringTime = time.Now().Add(-rtt * 4 * 1 / 3)
	ctr.AddRecvedBytesCount(dataRead)
	offset := ctr.recvWindowUpdate()
	if offset == 0 {
		fmt.Println("offset == 0 error")
		t.Fail()
	}
	newWindowSize := ctr.recvWindowSize
	if newWindowSize != oldWindowSize {
		fmt.Printf("newWindowSize != 2 * oldWIndowSize [%d] - [%d] \n", newWindowSize, oldWindowSize)
		t.Fail()
	}
	if offset != bytesread + dataRead + newWindowSize {
		fmt.Println("offset != bytesread + dataRead + newWindowSize")
		t.Fail()
	}
}

func TestRectAutoTuring5(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	oldWindowSize := ctr.recvWindowSize
	ctr.maxRecvWindowSize = 5000

	bytesread := ctr.recvBytesCount
	rtt := scaleDuration(20 * time.Millisecond)
	setRTT(ctr, rtt, t)
	dataRead := receiveWindowSize * 2 / 3 - 1
	ctr.startAutoTuringOffset = ctr.recvBytesCount
	ctr.startAutoTuringTime = time.Now().Add(-rtt * 4 * 2 / 3)
	ctr.AddRecvedBytesCount(dataRead)
	offset := ctr.recvWindowUpdate()
	if offset == 0 {
		fmt.Println("offset == 0 error")
		t.Fail()
	}
	newWindowSize := ctr.recvWindowSize
	if newWindowSize != oldWindowSize {
		fmt.Printf("newWindowSize != 2 * oldWIndowSize [%d] - [%d] \n", newWindowSize, oldWindowSize)
		t.Fail()
	}
	if offset != bytesread + dataRead + newWindowSize {
		fmt.Println("offset != bytesread + dataRead + newWindowSize")
		t.Fail()
	}
}

func TestRectAutoTuring6(t *testing.T) {
	receiveLength := uint64(10000)
	receiveWindowSize := uint64(1000)
	ctr := &FlowControl { }
	ctr.rttStat = &RTTStat { }
	ctr.recvBytesCount = receiveLength - receiveWindowSize
	ctr.recvWindowSize = receiveWindowSize
	ctr.recvSize = receiveLength

	oldWindowSize := ctr.recvWindowSize
	ctr.maxRecvWindowSize = 5000

	reset := func() {
		ctr.startAutoTuringTime = time.Now().Add(-time.Millisecond)
		ctr.startAutoTuringOffset = ctr.recvBytesCount
		ctr.AddRecvedBytesCount(ctr.recvWindowSize / 2 + 1)
	}

	setRTT(ctr, scaleDuration(20 * time.Millisecond), t)
	reset()
	ctr.adjustWindowSize()
	if ctr.recvWindowSize != 2 * oldWindowSize {
		fmt.Printf("ctr.recvWindowSize != 2 * oldWindowSize [%d]\n", ctr.recvWindowSize)
		t.Fail()
	}
	reset()
	ctr.adjustWindowSize()
	if ctr.recvWindowSize != 2 * 2 * oldWindowSize {
		fmt.Printf("ctr.recvWindowSize != 2 * 2 * oldWindowSize [%d]\n", ctr.recvWindowSize)
		t.Fail()
	}
	reset()
	ctr.adjustWindowSize()
	if ctr.recvWindowSize != ctr.maxRecvWindowSize {
		fmt.Printf("ctr.recvWindowSize != ctr.maxRecvWindowSize (once) [%d]\n", ctr.recvWindowSize)
		t.Fail()
	}
	ctr.adjustWindowSize()
	if ctr.recvWindowSize != ctr.maxRecvWindowSize {
		fmt.Printf("ctr.recvWindowSize != ctr.maxRecvWindowSize (twice) [%d]\n", ctr.recvWindowSize)
		t.Fail()
	}
}