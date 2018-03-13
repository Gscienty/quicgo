package congestioncontrol

import (
	"time"
)

const (
	SLOWSTART_LOW_WINDOW		= uint64(16)
	SLOWSTART_MIN_SAMPLE		= 8
	SLOWSTART_DELAY_FACTOR_EXP	= 3
	SLOWSTART_MIN_THREHOLD		= uint64(4000)
	SLOWSTART_MAX_THREHOLD		= uint64(16000)

	SLOWSTART_START_LOW_WINDOW = uint64(16)
)

type SlowStart struct {
	end				uint64
	lastSend		uint64
	started			bool
	minRTT			time.Duration
	sampleCount		int
	hystartFound	bool
}

func (this *SlowStart) StartReceiveRound(lastNumber uint64) {
	this.end = lastNumber
	this.minRTT = 0
	this.sampleCount = 0
	this.started = true
}

func (this *SlowStart) IsEndOfRound(n uint64) bool {
	return this.end < n
}

func (this *SlowStart) ShouldExit(lastRTT time.Duration, minRTT time.Duration, congestionWindow uint64) bool {
	if this.started == false {
		this.StartReceiveRound(this.lastSend)
	}
	if this.hystartFound {
		return true
	}
	this.sampleCount++
	if (this.sampleCount <= SLOWSTART_MIN_SAMPLE) && (this.minRTT == 0 || this.minRTT > lastRTT) {
		this.minRTT = lastRTT
	}

	if this.sampleCount == SLOWSTART_MIN_SAMPLE {
		minRTTIncreaseThreholdus := uint64(minRTT / time.Microsecond >> SLOWSTART_DELAY_FACTOR_EXP)
		if minRTTIncreaseThreholdus > SLOWSTART_MAX_THREHOLD {
			minRTTIncreaseThreholdus = SLOWSTART_MAX_THREHOLD
		}
		if minRTTIncreaseThreholdus < SLOWSTART_MIN_THREHOLD {
			minRTTIncreaseThreholdus = SLOWSTART_MIN_THREHOLD
		}
		minRTTIncreaseThrehold := time.Duration(minRTTIncreaseThreholdus)
		if this.minRTT > (minRTT + minRTTIncreaseThrehold) {
			this.hystartFound = true
		}
	}

	return congestionWindow >= SLOWSTART_START_LOW_WINDOW && this.hystartFound
}

func (this *SlowStart) SetSendCount(n uint64) {
	this.lastSend = n
}

func (this *SlowStart) Acked(ackedNumber uint64) {
	if this.IsEndOfRound(ackedNumber) {
		this.started = false
	}
}

func (this *SlowStart) Started() bool {
	return this.started
}

func (this *SlowStart) Restart() {
	this.started = false
	this.hystartFound = false
}