package congestioncontrol

import (
	"time"
	"../utils"
	"../protocol"
)

const (
	CUBIC_SENDER_MAX_BURST_BYTESCOUNT					= 3 * protocol.DEFAULT_TCP_MSS
	CUBIC_SENDER_DEFAULT_MINIMUM_CONGRESTION_WINDOW		= 2
	CUBIC_SENDER_RENO_BETA								= 0.7
)

type connectionStats struct {
	slowStartPacketsLost	uint64
	slowStartBytesLost		uint64
}

type cubicSender struct {
	slowStart	SlowStart
	prr			PRR
	rtt			*utils.RTTStat
	stats		connectionStats
	cubic		*Cubic

	reno		bool

	largestSentPacketNumber uint64
	largestAckedPacketNumber uint64
	largestSentAtLastCutback uint64
	congestionWindow uint64
	slowStartThrehold uint64
	lastCutbackExitedSlowStart bool
	slowStartLargeReduction bool
	minCongestionWindow uint64
	maxTCPCongestionWindow uint64
	connectionNumber int
	congestionWindowCount uint64
	initialCongestionWindow uint64
	initialMaxCongestionWindow uint64
}

func CubicSenderNew (clock time.Time, rtt *utils.RTTStat, reno bool, initialCongestionWindow uint64, initialMaxCongestionWindow uint64) SendAlgoDebug {
	return &cubicSender {
		rtt:						rtt,
		initialCongestionWindow:	initialCongestionWindow,
		initialMaxCongestionWindow:	initialMaxCongestionWindow,
		congestionWindow:			initialCongestionWindow,
		minCongestionWindow:		CUBIC_SENDER_DEFAULT_MINIMUM_CONGRESTION_WINDOW,
		slowStartThrehold:			initialMaxCongestionWindow,
		maxTCPCongestionWindow:		initialMaxCongestionWindow,
		connectionNumber:			CUBIC_SENDER_DEFAULT_MINIMUM_CONGRESTION_WINDOW,
		cubic:						CubicNew(clock),
		reno:						reno,
	}
}

func (this *cubicSender) TimeUntilSend(inSending uint64) time.Duration {
	if this.InRecovery() && this.prr.TimeUntilSend(this.GetCongestionWindow(), inSending, this.GetSlowStartThrehold()) == 0 {
		return 0
	}

	delay := this.rtt.GetSmoothedRTT() / time.Duration(2 * this.GetCongestionWindow() / protocol.DEFAULT_TCP_MSS)
	if this.InSlowStart() == false {
		delay = delay * 8 / 5
	}
	return delay
}

func (this *cubicSender) PacketSent(sentTime time.Time, inSending uint64, packetNumber uint64, bytesCount uint64) bool {
	if this.InRecovery() {
		this.prr.SendCount(bytesCount)
	}
	this.largestSentPacketNumber = packetNumber
	this.slowStart.SetSendCount(packetNumber)
	return true
}

func (this *cubicSender) InRecovery() bool {
	return this.largestAckedPacketNumber <= this.largestSentAtLastCutback && this.largestAckedPacketNumber != 0
}

func (this *cubicSender) InSlowStart() bool {
	return this.GetCongestionWindow() < this.GetSlowStartThrehold()
}

func (this *cubicSender) GetCongestionWindow() uint64 {
	return this.congestionWindow * protocol.DEFAULT_TCP_MSS
}

func (this *cubicSender) GetSlowStartThrehold() uint64 {
	return this.slowStartThrehold * protocol.DEFAULT_TCP_MSS
}

func (this *cubicSender) ExitSlowStart() {
	this.slowStartThrehold = this.congestionWindow
}

func (this *cubicSender) SlowStartThreshold() uint64 {
	return this.slowStartThrehold
}

func (this *cubicSender) TryExitSlowStart() {
	if this.InSlowStart() && this.slowStart.ShouldExit(this.rtt.GetLastestRTT(), this.rtt.GetMinRTT(), this.GetCongestionWindow() / protocol.DEFAULT_TCP_MSS) {
		this.ExitSlowStart()
	}
}

func (this *cubicSender) isCWNDLimited(inSending uint64) bool {
	congestionWindow := this.GetCongestionWindow()
	if inSending >= congestionWindow {
		return true
	}
	availableBytes := congestionWindow - inSending
	slowStartLimited := this.InSlowStart() && inSending > congestionWindow / 2
	return slowStartLimited || availableBytes <= CUBIC_SENDER_MAX_BURST_BYTESCOUNT
}

func (this *cubicSender) tryIncreaseCWND(n uint64, ackedBytes uint64, inSending uint64) {
	if this.isCWNDLimited(inSending) {
		this.cubic.ApplicationLimited()
		return
	}
	if this.congestionWindow > this.maxTCPCongestionWindow {
		return
	}
	if this.InSlowStart() {
		this.congestionWindow++
		return
	}
	if this.reno {
		this.congestionWindowCount++
		if this.congestionWindowCount * uint64(this.connectionNumber) >= this.congestionWindow {
			this.congestionWindow++
			this.congestionWindowCount = 0
		}
	} else {
		minCongestionWindow := this.cubic.CongestionWindowAfterAck(this.congestionWindow, this.rtt.GetMinRTT())
		if minCongestionWindow > this.maxTCPCongestionWindow {
			minCongestionWindow = this.maxTCPCongestionWindow
		}
	}
}

func (this *cubicSender) Acked(n uint64, ackedBytes uint64, inSending uint64) {
	if n > this.largestAckedPacketNumber {
		this.largestAckedPacketNumber = n
	}
	if this.InRecovery() {
		this.prr.Acked(ackedBytes)
		return
	}
	this.tryIncreaseCWND(n, ackedBytes, inSending)
	if this.InSlowStart() {
		this.slowStart.Acked(n)
	}
}

func (this *cubicSender) Lost(n uint64, lostBytes uint64, inSending uint64) {
	if n <= this.largestSentAtLastCutback {
		if this.lastCutbackExitedSlowStart {
			this.stats.slowStartPacketsLost++
			this.stats.slowStartBytesLost += lostBytes
			if this.slowStartLargeReduction {
				if this.stats.slowStartPacketsLost == 1 || (this.stats.slowStartBytesLost / protocol.DEFAULT_TCP_MSS) > (this.stats.slowStartBytesLost - lostBytes) / protocol.DEFAULT_TCP_MSS {
					this.congestionWindow = this.congestionWindow - 1
					if this.congestionWindow < this.minCongestionWindow {
						this.congestionWindow = this.minCongestionWindow
					}
				}
				this.slowStartThrehold = this.congestionWindow
			}
		}
		return
	}
	this.lastCutbackExitedSlowStart = this.InSlowStart()
	if this.InSlowStart() {
		this.stats.slowStartPacketsLost++
	}
	this.prr.Lost(inSending)

	if this.slowStartLargeReduction && this.InSlowStart() {
		this.congestionWindow = this.congestionWindow - 1
	} else if this.reno {
		this.congestionWindow = uint64(float32(this.congestionWindow) * this.RenoBeta())
	} else {
		this.congestionWindow = this.cubic.CongestionWindowAfterPacketLoss(this.congestionWindow)
	}

	if this.congestionWindow < this.minCongestionWindow {
		this.congestionWindow = this.minCongestionWindow
	}
	this.slowStartThrehold = this.congestionWindow
	this.largestSentAtLastCutback = this.largestSentPacketNumber
	this.congestionWindowCount = 0
}

func (this *cubicSender) RenoBeta() float32 {
	return (float32(this.connectionNumber) - 1. + CUBIC_SENDER_RENO_BETA) / float32(this.connectionNumber)
}

func (this *cubicSender) SetConnectionNumber(n int) {
	this.connectionNumber = n
	if this.connectionNumber < 1 {
		this.connectionNumber = 1
	}
	this.cubic.SetConnectionNumber(this.connectionNumber)
}

func (this *cubicSender) RetransmissionTimeout(packetsRetransmitted bool) {
	this.largestSentAtLastCutback = 0
	this.slowStart.Restart()
	this.cubic.Reset()
	this.slowStartThrehold = this.congestionWindow / 2
	this.congestionWindow = this.minCongestionWindow
}

func (this *cubicSender) ConnectionMigration() {
	this.slowStart.Restart()
	this.prr = PRR { }
	this.largestSentPacketNumber = 0
	this.largestAckedPacketNumber = 0
	this.largestSentAtLastCutback = 0
	this.lastCutbackExitedSlowStart = false
	this.cubic.Reset()
	this.congestionWindowCount = 0
	this.congestionWindow = this.initialCongestionWindow
	this.slowStartThrehold = this.initialMaxCongestionWindow
	this.maxTCPCongestionWindow = this.initialMaxCongestionWindow
}

func (this *cubicSender) SetSlowStartLargeReduction(enable bool) {
	this.slowStartLargeReduction = enable
}

func (this *cubicSender) RetransmissionDelay() time.Duration {
	if this.rtt.GetSmoothedRTT() == 0 {
		return 0
	}
	return this.rtt.GetSmoothedRTT() + this.rtt.GetMeanDeviation() * 4
}

func (this *cubicSender) BandwidthEstimate() uint64 {
	srtt := this.rtt.GetSmoothedRTT()
	if srtt == 0 {
		return 0
	}
	return this.GetCongestionWindow() * uint64(time.Second) / uint64(srtt)
}

func (this *cubicSender) GetSlowStart() *SlowStart {
	return &this.slowStart
}