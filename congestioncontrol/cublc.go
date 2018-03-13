package congestioncontrol

import (
	"math"
	"time"
)

const (
	CUBIC_DEFAULT_CONNECTION_NUMBER 	= 2
	CUBIC_BETA							= 0.7
	CUBIC_BETA_LAST_MAX					= 0.85
	CUBIC_MAX_TIME_INTERVAL				= 30 * time.Millisecond
	CUBIC_CUBE_SCALE					= 40
	CUBIC_CUBE_CONGESTION_WINDOW_SCALE	= 410
	CUBIC_CUBE_FACTOR					= uint64(1) << CUBIC_CUBE_SCALE / CUBIC_CUBE_CONGESTION_WINDOW_SCALE
)

type Cubic struct {
	clock							time.Time
	connectionNumber				int
	epoch							time.Time
	appLimitedStartTime				time.Time
	lastUpdateTime					time.Time
	lastCongestionWindow			uint64
	lastMaxCongestionWindow			uint64
	ackedCount						uint64
	estimatedTCPCongestionWindow	uint64
	originPointCongestionWindow		uint64
	timeToOriginPoint				uint32
	lastTargetCongestionWindow		uint64
}

func CubicNew(clock time.Time) *Cubic {
	ret := &Cubic {
		clock:				clock,
		connectionNumber:	2,
	}
	ret.Reset()
	return ret
}

func (this *Cubic) Reset() {
	this.epoch = time.Time { }
	this.appLimitedStartTime = time.Time { }
	this.lastUpdateTime = time.Time { }
	this.lastCongestionWindow = 0
	this.lastMaxCongestionWindow = 0
	this.ackedCount = 0
	this.estimatedTCPCongestionWindow = 0
	this.originPointCongestionWindow = 0
	this.timeToOriginPoint = 0
	this.lastTargetCongestionWindow = 0
}

func (this *Cubic) alpha() float32 {
	return 3 * float32(this.connectionNumber) * float32(this.connectionNumber) * (1 - this.beta()) / (1 + this.beta())
}

func (this *Cubic) beta() float32 {
	return (float32(this.connectionNumber) - 1 + CUBIC_BETA) / float32(this.connectionNumber)
}

func (this *Cubic) ApplicationLimited() {
	if this.appLimitedStartTime.IsZero() {
		this.appLimitedStartTime = time.Now()
	}
}

func (this *Cubic) CongestionWindowAfterPacketLoss(congestionWindow uint64) uint64 {
	if congestionWindow < this.lastMaxCongestionWindow {
		this.lastMaxCongestionWindow = uint64(CUBIC_BETA_LAST_MAX * float32(congestionWindow))
	} else {
		this.lastMaxCongestionWindow = congestionWindow
	}
	this.epoch = time.Time { }
	return uint64(float32(congestionWindow) * this.beta())
}

func (this *Cubic) CongestionWindowAfterAck(congestionWindow uint64, delayMin time.Duration) uint64 {
	this.ackedCount++
	currentTime := time.Now()

	if this.lastCongestionWindow == congestionWindow && (currentTime.Sub(this.lastUpdateTime) <= CUBIC_MAX_TIME_INTERVAL) {
		if this.lastTargetCongestionWindow < this.estimatedTCPCongestionWindow {
			return this.estimatedTCPCongestionWindow
		} else {
			return this.lastTargetCongestionWindow
		}
	}
	this.lastCongestionWindow = congestionWindow
	this.lastUpdateTime = currentTime

	if this.epoch.IsZero() {
		this.epoch = currentTime
		this.ackedCount = 1
		this.estimatedTCPCongestionWindow = congestionWindow
		if this.lastMaxCongestionWindow <= congestionWindow {
			this.timeToOriginPoint = 0
			this.originPointCongestionWindow = congestionWindow
		} else {
			this.timeToOriginPoint = uint32(math.Cbrt(float64(CUBIC_CUBE_FACTOR * (this.lastMaxCongestionWindow - congestionWindow))))
			this.originPointCongestionWindow = this.lastMaxCongestionWindow
		}
	} else {
		if this.appLimitedStartTime.IsZero() == false {
			shift := currentTime.Sub(this.appLimitedStartTime)
			this.epoch = this.epoch.Add(shift)
			this.appLimitedStartTime = time.Time { }
		}
	}

	elapsedTime := int64((currentTime.Add(delayMin).Sub(this.epoch) / time.Microsecond) << 10) / 1000000
	offset := int64(this.timeToOriginPoint) - elapsedTime
	if offset < 0 {
		offset = -offset
	}
	deltaCongestionWindow := uint64 ((CUBIC_CUBE_CONGESTION_WINDOW_SCALE * offset * offset * offset) >> CUBIC_CUBE_SCALE)
	var targetCongestionWindow uint64
	if elapsedTime > int64(this.timeToOriginPoint) {
		targetCongestionWindow = this.originPointCongestionWindow + deltaCongestionWindow
	} else {
		targetCongestionWindow = this.originPointCongestionWindow - deltaCongestionWindow
	}

	for {
		requiredAckCount := uint64(float32(this.estimatedTCPCongestionWindow) / this.alpha())
		if this.ackedCount < requiredAckCount {
			break
		}
		this.ackedCount -= requiredAckCount
		this.estimatedTCPCongestionWindow++
	}

	this.lastTargetCongestionWindow = targetCongestionWindow
	if targetCongestionWindow < this.estimatedTCPCongestionWindow {
		targetCongestionWindow = this.estimatedTCPCongestionWindow
	}
	return targetCongestionWindow
}

func (this *Cubic) SetConnectionNumber(n int) {
	this.connectionNumber = n
}