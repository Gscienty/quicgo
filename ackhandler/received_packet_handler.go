package ackhandler

import (
	"time"
	"../protocol"
	"../utils"
	"../frame"
)

const (
	RECEIVED_PACKET_HANDLER_ACK_SEND_DELAY							= 25 * time.Millisecond
	RECEIVED_PACKET_HANDLER_MIN_RECEIVED_BEFORE_ACK_DECIMATION		= 100
	RECEIVED_PACKET_HANDLER_RETRANSMITTABLE_PACKETS_BEFORE_ACK		= 10
	RECEIVED_PACKET_HANDLER_INIT_RETRANSMITTABLE_PACKETS_BEFORE_ACK	= 2
	RECEIVED_PACKET_HANDLER_ACK_DECIMATION_DELAY					= 1.0 / 4
	RECEIVED_PACKET_HANDLER_SHORT_ACK_DECIMATION_DELAY				= 1.0 / 8
	RECEIVED_PACKET_HANDLER_MAX_PACKETS_AFTER_NEW_MISSING			= 4

)

type receivedPacketHandler struct {
	largestObserved				protocol.PacketNumber
	ignoreBelow					protocol.PacketNumber
	largestObservedReceivedTime	time.Time

	packetHistory				*receivedPacketHistory
	
	ackSendDelay				time.Duration
	rttStats					*utils.RTTStat

	packetsReceivedSinceLaskAck					int
	retransmittablePacketsReceivedSinceLaskAck	int
	ackQueued									bool
	ackAlarm									time.Time
	lastAckFrame								*frame.AckFrame

	version	protocol.Version
}

func ReceivedPacketHandlerNew(rtt *utils.RTTStat, version protocol.Version) IReceivedPacketHandler {
	return &receivedPacketHandler {
		packetHistory:	receivedPacketHistoryNew(),
		ackSendDelay:	RECEIVED_PACKET_HANDLER_ACK_SEND_DELAY,
		rttStats:		rtt,
		version:		version,
	}
}

func (this *receivedPacketHandler) isMissing(p protocol.PacketNumber) bool {
	if this.lastAckFrame == nil {
		return false
	}

	return uint64(p) < this.lastAckFrame.LargestAcknowledged.GetVal() && this.lastAckFrame.AcksPacket(uint64(p))
}

func (this *receivedPacketHandler) hasNewMissingPackets() bool {
	if this.lastAckFrame == nil {
		return false
	}
	highestRange := this.packetHistory.GetHighestAckRange()
	return highestRange.First >= this.lastAckFrame.LargestAcknowledged.GetVal() && (highestRange.Last - highestRange.First + 1) <= RECEIVED_PACKET_HANDLER_MAX_PACKETS_AFTER_NEW_MISSING
}

func (this *receivedPacketHandler) maybeQueueAck(packetNumber protocol.PacketNumber, recvTime time.Time, shouldInstigateAck, wasMissing bool) {
	this.packetsReceivedSinceLaskAck++

	if this.lastAckFrame == nil {
		this.ackQueued = true
		return
	}

	if wasMissing {
		this.ackQueued = true
	}

	if !this.ackQueued && shouldInstigateAck {
		this.retransmittablePacketsReceivedSinceLaskAck++

		if packetNumber > RECEIVED_PACKET_HANDLER_MIN_RECEIVED_BEFORE_ACK_DECIMATION {
			if this.retransmittablePacketsReceivedSinceLaskAck >= RECEIVED_PACKET_HANDLER_RETRANSMITTABLE_PACKETS_BEFORE_ACK {
				this.ackQueued = true
			} else if this.ackAlarm.IsZero() {
				ackDelay := time.Duration(float64(this.rttStats.GetMinRTT()) * float64(RECEIVED_PACKET_HANDLER_ACK_DECIMATION_DELAY))
				if ackDelay > RECEIVED_PACKET_HANDLER_ACK_SEND_DELAY {
					ackDelay = RECEIVED_PACKET_HANDLER_ACK_SEND_DELAY
				}

				this.ackAlarm = recvTime.Add(ackDelay)
			}
		} else {
			if this.retransmittablePacketsReceivedSinceLaskAck >= RECEIVED_PACKET_HANDLER_INIT_RETRANSMITTABLE_PACKETS_BEFORE_ACK {
				this.ackQueued = true
			} else if this.ackAlarm.IsZero() {
				this.ackAlarm = recvTime.Add(RECEIVED_PACKET_HANDLER_ACK_SEND_DELAY)
			}
		}

		if this.hasNewMissingPackets() {
			ackDelay := float64(this.rttStats.GetMinRTT()) * float64(RECEIVED_PACKET_HANDLER_SHORT_ACK_DECIMATION_DELAY)
			ackTime := recvTime.Add(time.Duration(ackDelay))
			if this.ackAlarm.IsZero() || this.ackAlarm.After(ackTime) {
				this.ackAlarm = ackTime
			}
		}
	}

	if this.ackQueued {
		this.ackAlarm = time.Time { }
	}
}

func (this *receivedPacketHandler) ReceivedPacket(packetNumber protocol.PacketNumber, recvTime time.Time, shouldInstigateAck bool) error {
	if packetNumber < this.ignoreBelow {
		return nil
	}

	isMissing := this.isMissing(packetNumber)
	if packetNumber > this.largestObserved {
		this.largestObserved = packetNumber
		this.largestObservedReceivedTime = recvTime
	}

	if err := this.packetHistory.ReceivedPacket(packetNumber); err != nil {
		return err
	}

	this.maybeQueueAck(packetNumber, recvTime, shouldInstigateAck, isMissing)
	return nil
}

func (this *receivedPacketHandler) IgnoreBelow(p protocol.PacketNumber) {
	this.ignoreBelow = p
	this.packetHistory.DeleteBelow(p)
}

func (this *receivedPacketHandler) GetAlarmTimeout() time.Time {
	return this.ackAlarm
}

func (this *receivedPacketHandler) GetAckFrame() *frame.AckFrame {
	if !this.ackQueued && (this.ackAlarm.IsZero() || this.ackAlarm.After(time.Now())) {
		return nil
	}

	ackRanges := this.packetHistory.GetAckRanges()
	ack := &frame.AckFrame {
		LargestAcknowledged:	*utils.VarLenIntegerStructNew(uint64(this.largestObserved)),
		LowAcknowledged:		ackRanges[len(ackRanges) - 1].First,
		PacketReceivedTime:		this.largestObservedReceivedTime,
	}

	if len(ackRanges) > 1 {
		ack.Blocks = ackRanges
	}

	this.lastAckFrame = ack
	this.ackAlarm = time.Time { }
	this.ackQueued = false
	this.packetsReceivedSinceLaskAck = 0
	this.retransmittablePacketsReceivedSinceLaskAck = 0
	return ack
}