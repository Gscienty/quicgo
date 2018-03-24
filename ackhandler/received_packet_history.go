package ackhandler

import (
	"../frame"
	"../protocol"
	"errors"
)

type packetInterval struct {
	start	protocol.PacketNumber
	end		protocol.PacketNumber
}

type packetIntervalElement struct {
	next	*packetIntervalElement
	prev	*packetIntervalElement
	list	*packetIntervalList
	Value	packetInterval
}

type packetIntervalList struct {
	root	packetIntervalElement
	len		int
}

type receivedPacketHistory struct {
	ranges							*packetIntervalList
	lowestInReceivedPacketNumbers	protocol.PacketNumber
}

func (this *packetIntervalElement) Next() *packetIntervalElement {
	if p := this.next; this.list != nil && p != &this.list.root {
		return p
	}
	return nil
}

func (this *packetIntervalElement) Prev() *packetIntervalElement {
	if p := this.prev; this.list != nil && p != &this.list.root {
		return p
	}
	return nil
}

func (this *packetIntervalList) Init() *packetIntervalList {
	this.root.next = &this.root
	this.root.prev = &this.root
	this.len = 0
	return this
}

func packetIntervalListNew() *packetIntervalList {
	return new(packetIntervalList).Init()
}

func (this *packetIntervalList) Length() int {
	return this.len
}

func (this *packetIntervalList) Front() *packetIntervalElement {
	if this.len == 0 {
		return nil
	}
	return this.root.next
}

func (this *packetIntervalList) Back() *packetIntervalElement {
	if this.len == 0 {
		return nil
	}
	return this.root.prev
}

func (this *packetIntervalList) lazyInit() {
	if this.root.next == nil {
		this.Init()
	}
}

func (this *packetIntervalList) insert(e, at *packetIntervalElement) *packetIntervalElement {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = this
	this.len++
	return e
}

func (this *packetIntervalList) insertValue(v packetInterval, at *packetIntervalElement) *packetIntervalElement {
	return this.insert(&packetIntervalElement { Value: v }, at)
}

func (this *packetIntervalList) remove(e *packetIntervalElement) *packetIntervalElement {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	this.len--
	return e
}

func (this *packetIntervalList) Remove(e *packetIntervalElement) packetInterval {
	if e.list == this {
		this.remove(e)
	}
	return e.Value
}

func (this *packetIntervalList) PushFront(v packetInterval) *packetIntervalElement {
	this.lazyInit()
	return this.insertValue(v, &this.root)
}

func (this *packetIntervalList) PushBack(v packetInterval) *packetIntervalElement {
	this.lazyInit()
	return this.insertValue(v, this.root.prev)
}

func (this *packetIntervalList) InsertBefore(v packetInterval, mark *packetIntervalElement) *packetIntervalElement {
	if (mark.list != this) {
		return nil
	}
	return this.insertValue(v, mark.prev)
}

func (this *packetIntervalList) InsertAfter(v packetInterval, mark *packetIntervalElement) *packetIntervalElement {
	if (mark.list != this) {
		return nil
	}
	return this.insertValue(v, mark)
}

func (this *packetIntervalList) MoveToFront(e *packetIntervalElement) {
	if e.list != this || this.root.next == e {
		return
	}
	this.insert(this.remove(e), &this.root)
}

func (this *packetIntervalList) MoveToBack(e *packetIntervalElement) {
	if e.list != this || this.root.prev == e {
		return
	}
	this.insert(this.remove(e), this.root.prev)
}

func (this *packetIntervalList) MoveBefore(e, mark *packetIntervalElement) {
	if e.list != this || e == mark || mark.list != this {
		return
	}
	this.insert(this.remove(e), mark.prev)
}

func (this *packetIntervalList) MoveAfter(e, mark *packetIntervalElement) {
	if e.list != this || e == mark || mark.list != this {
		return
	}
	this.insert(this.remove(e), mark)
}

func (this *packetIntervalList) PushBackList(other *packetIntervalList) {
	this.lazyInit()
	for i, e := other.Length(), other.Front(); i > 0; i, e = i - 1, e.Next() {
		this.insertValue(e.Value, this.root.prev)
	}
}

func (this *packetIntervalList) PushFrontList(other *packetIntervalList) {
	this.lazyInit()
	for i, e := other.Length(), other.Back(); i > 0; i, e = i - 1, e.Prev() {
		this.insertValue(e.Value, &this.root)
	}
}

func receivedPacketHistoryNew() *receivedPacketHistory {
	return &receivedPacketHistory { ranges: packetIntervalListNew() }
}

func (this *receivedPacketHistory) ReceivedPacket(p protocol.PacketNumber) error {
	if this.ranges.Length() >= protocol.MAX_TRACKED_RECEIVED_ACK_RANGES {
		return errors.New("too many outstanding received ack rangs")
	}

	if this.ranges.Length() == 0 {
		this.ranges.PushBack(packetInterval { start: p, end: p })
		return nil
	}

	for e := this.ranges.Back(); e != nil; e = e.Prev() {
		if p >= e.Value.start && p <= e.Value.end {
			return nil
		}

		var rangeExtended bool
		if e.Value.end == p - 1 {
			e.Value.end = p
			rangeExtended = true
		} else if e.Value.start == p + 1 {
			e.Value.start = p
			rangeExtended = true
		}

		if rangeExtended {
			prev := e.Prev()
			if prev != nil && prev.Value.end + 1 == e.Value.start {
				prev.Value.end = e.Value.end
				this.ranges.Remove(e)
				return nil
			}
			return nil
		}

		if p > e.Value.end {
			this.ranges.InsertAfter(packetInterval { start: p, end: p }, e)
			return nil
		}
	}

	this.ranges.InsertBefore(packetInterval { start: p, end: p }, this.ranges.Front())

	return nil
}

func (this *receivedPacketHistory) DeleteBelow(p protocol.PacketNumber) {
	if p <= this.lowestInReceivedPacketNumbers {
		return
	}

	this.lowestInReceivedPacketNumbers = p

	nextElement := this.ranges.Front()
	for element := this.ranges.Front(); nextElement != nil; element = nextElement {
		nextElement = element.Next()

		if p > element.Value.start && p <= element.Value.end {
			element.Value.start = p
		} else if element.Value.end < p {
			this.ranges.Remove(element)
		} else {
			return
		}
	}
}

func (this *receivedPacketHistory) GetAckRanges() []frame.AckBlock {
	if this.ranges.Length() == 0 {
		return nil
	}

	ackRanges := make([]frame.AckBlock, this.ranges.Length())
	i := 0
	for element := this.ranges.Back(); element != nil; element = element.Prev() {
		ackRanges[i] = frame.AckBlock { First: uint64(element.Value.start), Last: uint64(element.Value.end) }
		i++
	}

	return ackRanges
}

func (this *receivedPacketHistory) GetHighestAckRange() frame.AckBlock {
	ackRange := frame.AckBlock { }
	if this.ranges.Length() > 0 {
		r := this.ranges.Back().Value
		ackRange.First = uint64(r.start)
		ackRange.Last = uint64(r.end)
	}

	return ackRange
}