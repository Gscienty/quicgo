package flowcontrol

import (
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
}