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
	rtt := RTTStatTestBefore()
	rtt.recentMinRTTWindow = 99 * time.Millisecond

	now := time.Time { }
	rttSimple := 10 * time.Millisecond
	rtt.Update(rttSimple, 0, now)
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
		t.Fail()
	}

	for i := 0; i < 8; i++ {
		now = now.Add(25 * time.Millisecond)
		rttSimple += 10 * time.Millisecond
		rtt.Update(rttSimple, 0, now)
		if rtt.minRTT != 10 * time.Millisecond {
			t.Fail()
		}
		if rtt.quarterWindowRTT.rtt != rttSimple {
			t.Fail()
		}
		if rtt.halfWindowRTT.rtt != rttSimple - (10 * time.Millisecond) {
			t.Fail()
		}
		if i < 3 {
			if rtt.recentMinRTT.rtt != 10 * time.Millisecond {
				t.Fail()
			}
		} else if i < 5 {
			if rtt.recentMinRTT.rtt != 30 * time.Millisecond {
				t.Fail()
			}
		} else if i < 7 {
			if rtt.recentMinRTT.rtt != 50 * time.Millisecond {
				t.Fail()
			}
		} else {
			if rtt.recentMinRTT.rtt != 70 * time.Millisecond {
				t.Fail()
			}
		}
	}

	rttSimple -= 5 * time.Millisecond
	rtt.Update(rttSimple, 0, now)
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.quarterWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.halfWindowRTT.rtt != (rttSimple - (5 * time.Millisecond)) {
		fmt.Printf("%d %d\n", rtt.halfWindowRTT.rtt, (rttSimple - (5 * time.Millisecond)))
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 70 * time.Millisecond {
		t.Fail()
	}
	
	rttSimple -= 15 * time.Millisecond
	rtt.Update(rttSimple, 0, now)
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.quarterWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.halfWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 70 * time.Millisecond {
		t.Fail()
	}
	
	rttSimple = 65 * time.Millisecond
	rtt.Update(rttSimple, 0, now)
	if rtt.minRTT != 10 * time.Millisecond {
		t.Fail()
	}
	if rtt.quarterWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.halfWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != rttSimple {
		t.Fail()
	}
	
	rttSimple = 5 * time.Millisecond
	rtt.Update(rttSimple, 0, now)
	if rtt.minRTT != rttSimple {
		t.Fail()
	}
	if rtt.quarterWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.halfWindowRTT.rtt != rttSimple {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != rttSimple {
		t.Fail()
	}
}

func TestSmoothedMetrics(t *testing.T) {
	rtt := RTTStatTestBefore()
	init := 10 * time.Millisecond
	rtt.Update(init, 0, time.Time { })
	if rtt.minRTT != init {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != init {
		t.Fail()
	}
	if rtt.smoothedRTT != init {
		t.Fail()
	}
	if rtt.meanDeviation != init / 2 {
		t.Fail()
	}

	double := init * 2
	rtt.Update(double, 0, time.Time { })
	if rtt.smoothedRTT != time.Duration(float32(init) * 1.125) {
		t.Fail()
	}
	rtt.ExpireSmoothed()
	if rtt.smoothedRTT != double {
		t.Fail()
	}
	if rtt.meanDeviation != time.Duration(float32(init) * 0.875) {
		t.Fail()
	}
}

func TestBadSendDelta(t *testing.T) {
	rtt := RTTStatTestBefore()
	init := 10 * time.Millisecond
	rtt.Update(init, 0, time.Time { })
	if rtt.minRTT != init {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != init {
		t.Fail()
	}
	if rtt.smoothedRTT != init {
		t.Fail()
	}

	baddelta := []time.Duration {
		0,
		MAX_DURATION,
		-1000 * time.Microsecond,
	}

	for _, delta := range baddelta {
		rtt.Update(delta, 0, time.Time { })
		if rtt.minRTT != init {
			t.Fail()
		}
		if rtt.recentMinRTT.rtt != init {
			t.Fail()
		}
		if rtt.smoothedRTT != init {
			t.Fail()
		}
	}
}

func TestConnection(t *testing.T) {
	rtt := RTTStatTestBefore()
	rtt.Update(200 * time.Millisecond, 0, time.Time { })
	if rtt.lastestRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.smoothedRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.minRTT != 200 * time.Millisecond {
		t.Fail()
	}
	rtt.Update(300 * time.Millisecond, 100 * time.Millisecond, time.Time { })
	if rtt.lastestRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.smoothedRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.minRTT != 200 * time.Millisecond {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != 200 * time.Millisecond {
		t.Fail()
	}

	rtt.OnConnectionMigration()
	if rtt.lastestRTT != time.Duration(0) {
		t.Fail()
	}
	if rtt.smoothedRTT != time.Duration(0) {
		t.Fail()
	}
	if rtt.minRTT != time.Duration(0) {
		t.Fail()
	}
	if rtt.recentMinRTT.rtt != time.Duration(0) {
		t.Fail()
	}
}