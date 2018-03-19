package ackhandler

type PacketList struct {
	root	PacketElement
	len		int
}

type PacketElement struct {
	next	*PacketElement
	prev	*PacketElement
	list	*PacketList

	Value	Packet
}

func (this *PacketElement) Next() *PacketElement {
	if p := this.next; this.list != nil && p != &this.list.root {
		return p
	}
	return nil
}

func (this *PacketElement) Prev() *PacketElement {
	if p := this.prev; this.list != nil && p != &this.list.root {
		return p
	}
	return nil
}

func (this *PacketList) Init() *PacketList {
	this.root.next = &this.root
	this.root.prev = &this.root
	this.len = 0
	return this
}

func PacketListNew() *PacketList {
	return new(PacketList).Init()
}

func (this *PacketList) Length() int {
	return this.len
}

func (this *PacketList) Front() *PacketElement {
	if this.len == 0 {
		return nil
	}
	return this.root.next
}

func (this *PacketList) Back() *PacketElement {
	if this.len == 0 {
		return nil
	}
	return this.root.prev
}

func (this *PacketList) lazyInit() {
	if this.root.next == nil {
		this.Init()
	}
}

func (this *PacketList) insert(e, at *PacketElement) *PacketElement {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = this
	this.len++
	return e
}

func (this *PacketList) insertValue(v Packet, at *PacketElement) *PacketElement {
	return this.insert(&PacketElement { Value: v }, at)
}

func (this *PacketList) remove(e *PacketElement) *PacketElement {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	this.len--
	return e
}

func (this *PacketList) Remove(e *PacketElement) Packet {
	if e.list == this {
		this.remove(e)
	}
	return e.Value
}

func (this *PacketList) PushFront(v Packet) *PacketElement {
	this.lazyInit()
	return this.insertValue(v, &this.root)
}

func (this *PacketList) PushBack(v Packet) *PacketElement {
	this.lazyInit()
	return this.insertValue(v, this.root.prev);
}

func (this *PacketList) InsertAfter(v Packet, mark *PacketElement) *PacketElement {
	if mark.list != this {
		return nil
	}
	return this.insertValue(v, mark)
}

func (this *PacketList) MoveToFront(e *PacketElement) {
	if e.list != this || this.root.next == e {
		return
	}
	this.insert(this.remove(e), &this.root)
}

func (this *PacketList) MoveToBack(e *PacketElement) {
	if e.list != this || this.root.prev == e {
		return
	}
	this.insert(this.remove(e), this.root.prev)
}

func (this *PacketList) MoveBefore(e, mark *PacketElement) {
	if e.list != this || e == mark || mark.list != this {
		return
	}
	this.insert(this.remove(e), mark.prev)
}

func (this *PacketList) MoveAfter(e, mark *PacketElement) {
	if e.list != this || e == mark || mark.list != this {
		return
	}
	this.insert(this.remove(e), mark)
}

func (this *PacketList) PushBackList(other *PacketList) {
	this.lazyInit()
	for i, e := other.Length(), other.Front(); i > 0; i, e = i - 1, e.Next() {
		this.insertValue(e.Value, this.root.prev)
	}
}

func (this *PacketList) PushFrontList(other *PacketList) {
	this.lazyInit()
	for i, e := other.Length(), other.Back(); i > 0; i, e = i - 1, e.Prev() {
		this.insertValue(e.Value, &this.root)
	}
}