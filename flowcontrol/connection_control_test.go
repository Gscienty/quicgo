package flowcontrol

import (
	"fmt"
	"testing"
	"time"
)

func ConnectionControlTestBefore(t time.Duration) *connectionFlowControl {
	conn := &connectionFlowControl { }
	conn.rttStat = &RTTStat { }
	return conn
}

func ConnectionControlTestBeforeSetRTT(conn *connectionFlowControl, t time.Duration) {
	conn.rttStat.Update(t, 0, time.Time { })
}

func TestConstructor(t *testing.T) {
	recWin := uint64(2000)
	maxRecWin := uint64(3000)
	fc := ConnectionFlowControlNew(recWin, maxRecWin, &RTTStat { }).(*connectionFlowControl)

	if fc.recvSize != recWin {
		t.Fail()
	}
	if fc.maxRecvWindowSize != maxRecWin {
		t.Fail()
	}
}

func TestReceiveFlowControl(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))

	conn.recvHighestOffset = 1337
	conn.AddHighestOffset(123)
	if conn.recvHighestOffset != uint64(1337 + 123) {
		fmt.Printf("%d %d\n", conn.recvHighestOffset, uint64(1337 + 123))
		t.Fail()
	}
}

func TestConnectionWindowUpdate(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))

	conn.recvHighestOffset = 1337
	conn.AddHighestOffset(123)
	if conn.recvHighestOffset != uint64(1337 + 123) {
		fmt.Printf("%d %d\n", conn.recvHighestOffset, uint64(1337 + 123))
		t.Fail()
	}

	conn.recvSize = 100
	conn.recvWindowSize = 60
	conn.maxRecvWindowSize = 1000
	conn.recvBytesCount = 100 - 60

	ws := conn.recvWindowSize
	of := conn.recvBytesCount
	dr := ws / 2 - 1
	conn.AddRecvedBytesCount(dr)
	f := conn.RecvWindowUpdate()

	if f != of + dr + 60 {
		t.Fail()
	}
}

func TestConnectionAutoTurnWindow(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))

	conn.recvHighestOffset = 1337
	conn.AddHighestOffset(123)
	if conn.recvHighestOffset != uint64(1337 + 123) {
		fmt.Printf("%d %d\n", conn.recvHighestOffset, uint64(1337 + 123))
		t.Fail()
	}

	conn.recvSize = 100
	conn.recvWindowSize = 60
	conn.maxRecvWindowSize = 1000
	conn.recvBytesCount = 100 - 60

	of := conn.recvBytesCount
	ows := conn.recvWindowSize
	rtt := scaleDuration(20 * time.Millisecond)
	ConnectionControlTestBeforeSetRTT(conn, rtt)
	conn.startAutoTuringTime = time.Now().Add(-time.Millisecond)
	conn.startAutoTuringOffset = of
	dr := ows / 2 + 1
	conn.AddRecvedBytesCount(dr)
	f := conn.RecvWindowUpdate()
	nws := conn.recvWindowSize
	if nws != 2 * ows {
		t.Fail()
	}
	if f != uint64(of + dr + nws) {
		t.Fail()
	}
}

func TestSendFlowControlBlocked(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))

	conn.SetSendOffset(100)
	if b, _ := conn.IsNewlyBlocked(); b {
		t.Fail()
	}
	conn.AddSendedBytesCount(100)
	blocked, off := conn.IsNewlyBlocked()
	if blocked == false {
		t.Fail()
	}
	if off != uint64(100) {
		t.Fail()
	}
}

func TestSendFlowControlNewlyBlocked(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))

	conn.SetSendOffset(100)
	conn.AddSendedBytesCount(100)
	newlyBlocked, offset := conn.IsNewlyBlocked()
	if newlyBlocked != true {
		t.Fail()
	}
	if offset != uint64(100) {
		t.Fail()
	}
	newlyBlocked, _ = conn.IsNewlyBlocked()
	if newlyBlocked != false {
		t.Fail()
	}
	conn.SetSendOffset(150)
	conn.AddSendedBytesCount(150)
	newlyBlocked, _ = conn.IsNewlyBlocked()
	if newlyBlocked != true {
		t.Fail()
	}
}

func TestSettingMinimumWindowSize(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(1000)

	conn.recvSize = receiveSize
	conn.recvWindowSize = receiveWindowSize
	oldWindowSize := conn.recvWindowSize
	conn.maxRecvWindowSize = 3000

	conn.EnsureRecvMinimumWindowSize(1800)
	if conn.recvWindowSize != 1800 {
		fmt.Printf("%d %d\n", conn.recvWindowSize, oldWindowSize)
		t.Fail()
	}
}

func TestSettingReduceWindowSize(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(1000)

	conn.recvSize = receiveSize
	conn.recvWindowSize = receiveWindowSize
	oldWindowSize := conn.recvWindowSize
	conn.maxRecvWindowSize = 3000

	conn.EnsureRecvMinimumWindowSize(1)
	if conn.recvWindowSize != oldWindowSize {
		fmt.Printf("%d %d\n", conn.recvWindowSize, oldWindowSize)
		t.Fail()
	}
}

func TestSettingBeyoundWindowSize(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(1000)

	conn.recvSize = receiveSize
	conn.recvWindowSize = receiveWindowSize
	conn.maxRecvWindowSize = 3000

	max := conn.maxRecvWindowSize
	conn.EnsureRecvMinimumWindowSize(2 * max)
	if conn.recvWindowSize != max {
		fmt.Printf("%d %d\n", conn.recvWindowSize, max)
		t.Fail()
	}
}

func TestSettingAfterWindowSize(t *testing.T) {
	conn := ConnectionControlTestBefore(time.Duration(0))
	receiveSize := uint64(10000)
	receiveWindowSize := uint64(1000)

	conn.recvSize = receiveSize
	conn.recvWindowSize = receiveWindowSize
	conn.maxRecvWindowSize = 3000

	conn.EnsureRecvMinimumWindowSize(1912)
	if conn.startAutoTuringTime.Before(time.Now().Add(-100 * time.Millisecond)) && conn.startAutoTuringTime.After(time.Now().Add(100 * time.Millisecond)) {
		t.Fail()
	}
}
