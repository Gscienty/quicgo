package flowcontrol

import (
	"math"
	"time"
)

const (
	MAX_DURATION time.Duration	= time.Duration (math.MaxInt64)

	INITIAL_RTT_US				= 100 * 1000
	RTT_ALPHA					= 0.125
	RTT_BETA					= 0.25
	HARF_WINDOW					= 0.5
	QUARTER_WINDOW				= 0.25
)

type RTT struct {
	rtt		time.Duration
	time	time.Time
}

type RTTStat struct {
	recentMinRTTWindow	time.Duration
	minRTT				time.Duration
	lastestRTT			time.Duration
	smoothedRTT			time.Duration
	meanDeviation		time.Duration

	samplesCount		uint64

	currentMinRTT		RTT
	recentMinRTT		RTT
	halfWindowRTT		RTT
	quarterWindowRTT	RTT
}

func absDuration (t time.Duration) time.Duration {
	if t < 0 {
		return -t
	}
	return t
}
func maxDuration (a time.Duration, b time.Duration) time.Duration {
	if a > b {
		return a
	} else {
		return b
	}
}

func RTTStatNew () *RTTStat {
	return &RTTStat { recentMinRTTWindow: MAX_DURATION }
}

func (this *RTTStat) Update (delta time.Duration, delay time.Duration, current time.Time) {
	if delta == MAX_DURATION || delta <= 0 {
		return
	}

	if this.minRTT == 0 || this.minRTT > delta {
		this.minRTT = delta
	}

	if this.samplesCount > 0 {
		this.samplesCount--
		if this.currentMinRTT.rtt == 0 || delta <= this.currentMinRTT.rtt {
			this.currentMinRTT = RTT { delta, current }
		}
		if this.samplesCount == 0 {
			this.recentMinRTT = this.currentMinRTT
			this.halfWindowRTT = this.currentMinRTT
			this.quarterWindowRTT = this.currentMinRTT
		}
	}

	if this.recentMinRTT.rtt == 0 || delta <= this.recentMinRTT.rtt {
		this.recentMinRTT = RTT { delta, current }
		this.halfWindowRTT = this.recentMinRTT
		this.quarterWindowRTT = this.recentMinRTT
	} else if delta <= this.halfWindowRTT.rtt {
		this.halfWindowRTT = RTT { delta, current }
		this.quarterWindowRTT = this.halfWindowRTT
	} else if delta <= this.quarterWindowRTT.rtt {
		this.quarterWindowRTT = RTT { delta, current }
	}

	if this.recentMinRTT.time.Before (current.Add (-this.recentMinRTTWindow)) {
		this.recentMinRTT = this.halfWindowRTT
		this.halfWindowRTT = this.quarterWindowRTT
		this.quarterWindowRTT = RTT { delta, current }
	} else if this.halfWindowRTT.time.Before (current.Add (
		-time.Duration (float32 (this.recentMinRTTWindow / time.Millisecond) * HARF_WINDOW) * time.Microsecond)) {
		this.halfWindowRTT = this.quarterWindowRTT
		this.quarterWindowRTT = RTT { delta, current }
	} else if this.quarterWindowRTT.time.Before (current.Add(
		-time.Duration (float32 (this.recentMinRTTWindow / time.Millisecond) * QUARTER_WINDOW) * time.Microsecond)) {
		this.quarterWindowRTT = RTT { delta, current }
	}

	if delta - this.minRTT >= delay {
		delta -= delay
	}
	this.lastestRTT = delta
	if this.smoothedRTT == 0 {
		this.smoothedRTT = delta
		this.meanDeviation = delta / 2
	} else {
		this.meanDeviation = time.Duration (
			(1 - RTT_BETA) * float32 (this.meanDeviation / time.Millisecond) +
			RTT_BETA * float32 (absDuration((this.smoothedRTT - delta) / time.Millisecond))) * time.Microsecond
		this.smoothedRTT = time.Duration (
			(1 - RTT_ALPHA) * float32 (this.smoothedRTT / time.Millisecond) +
			RTT_ALPHA * float32 (delta / time.Millisecond)) * time.Millisecond
	}
}

func (this *RTTStat) SimpleCurrentRecentMinRTT (n uint64) {
	this.samplesCount = n
	this.currentMinRTT = RTT { }
}

func (this *RTTStat) OnConnectionMigration () {
	this.lastestRTT = 0
	this.minRTT = 0
	this.smoothedRTT = 0
	this.meanDeviation = 0

	this.samplesCount = 0

	this.recentMinRTTWindow = MAX_DURATION
	this.recentMinRTT = RTT { }
	this.halfWindowRTT = RTT { }
	this.quarterWindowRTT = RTT { }
}

func (this *RTTStat) ExpireSmoothed () {
	this.meanDeviation = maxDuration (this.meanDeviation, absDuration (this.smoothedRTT - this.lastestRTT))
	this.smoothedRTT = maxDuration (this.smoothedRTT, this.lastestRTT)
}
