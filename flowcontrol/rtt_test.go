package flowcontrol

import (
	"fmt"
	"time"
	"testing"
)

func RTTStatTestBefore() *RTTStat {
	return RTTStatNew()
}

func TestRTTStatDefault(t *testing.T) {
	rtt := RTTStatTestBefore()
	if rtt.minRTT != time.Duration(0) {
		t.Fail()
	}
	if rtt.smoothedRTT != time.Duration(0) {
		t.Fail()
	}
}

func TestRTTStatSmoothed(t *testing.T) {
	rtt := RTTStatTestBefore()
	rtt.Update(300 * time.Millisecond, 100 * time.Millisecond, time.Time { })
	if rtt.lastestRTT != 300 * time.Millisecond {
		t.Fail()
	}
	if rtt.smoothedRTT != 300 * time.Millisecond {
		t.Fail()
	}
	rtt.Update(300 * time.Millisecond, 50 * time.Millisecond, time.Time { })
	if rtt.lastestRTT != 300 * time.Millisecond {
		t.Fail()
	}
	if rtt.smoothedRTT != 300 * time.Millisecond {
		t.Fail()
	}
	rtt.Update(200 * time.Millisecond, 300 * time.Millisecond, time.Time { })
	if rtt.lastestRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.smoothedRTT != 287500 * time.Microsecond {
		t.Fail()
	}
}

func TestMinRTT(t *testing.T) {
	rtt := RTTStatTestBefore()
	rtt.Update(200 * time.Millisecond, 0, time.Time { })
	if rtt.minRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 200 * time.Millisecond {
		t.Fail()
	}

	rtt.Update(10 * time.Millisecond, 0, time.Time { }.Add(10 * time.Millisecond))
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
		t.Fail()
	}

	
	rtt.Update(50 * time.Millisecond, 0, time.Time { }.Add(20 * time.Millisecond))
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
		t.Fail()
	}

	
	rtt.Update(50 * time.Millisecond, 0, time.Time { }.Add(30 * time.Millisecond))
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
		t.Fail()
	}

	
	rtt.Update(50 * time.Millisecond, 0, time.Time { }.Add(40 * time.Millisecond))
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
		t.Fail()
	}
	
	rtt.Update(7 * time.Millisecond, 2 * time.Millisecond, time.Time { }.Add(50 * time.Millisecond))
	if rtt.minRTT != 7 * time.Millisecond {
		fmt.Printf("%d %d", rtt.minRTT, 7 * time.Millisecond)
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 7 * time.Millisecond {
		fmt.Printf("%d %d", rtt.recentMinRTT.rtt, 7 * time.Millisecond)
		t.Fail()
	}
}

func TestRecentMinRTT(t *testing.T) {
	rtt := RTTStatTestBefore()
	rtt.Update(10 * time.Millisecond, 0, time.Time { })
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
		t.Fail()
	}

	rtt.SimpleCurrentRecentMinRTT(4)
	for i := 0; i < 3; i++ {
		rtt.Update(50 * time.Millisecond, 0, time.Time { })
		if rtt.minRTT != 10 * time.Millisecond {
			t.Fail()
		}
		if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
			t.Fail()
		}
	}

	rtt.Update(50 * time.Millisecond, 0, time.Time { })
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 50 * time.Millisecond {
		t.Fail()
	}
}

func TestWindowRecentMinRTT(t *testing.T) {
	
}