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