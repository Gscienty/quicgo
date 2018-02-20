package protocol

import (
	"encoding/binary"
)

const (
	TAG_CHLO	= uint32 (('C') + ('H' << 8) + ('L' << 16) + ('O' << 24))
	TAG_SHLO	= uint32 (('S') + ('H' << 8) + ('L' << 16) + ('O' << 24))
	TAG_REJ		= uint32 (('R') + ('E' << 8) + ('J' << 16))
	TAG_SCUP	= uint32 (('S') + ('C' << 8) + ('U' << 16) + ('P' << 24))
	TAG_PRST	= uint32 (('P') + ('R' << 8) + ('S' << 16) + ('T' << 24))
	TAG_VERS	= uint32 (('V') + ('E' << 8) + ('R' << 16) + ('S' << 24))
	TAG_PAD		= uint32 (('P') + ('A' << 8) + ('D' << 16))
	TAG_STK		= uint32 (('S') + ('T' << 8) + ('K' << 16))
	TAG_SNI		= uint32 (('S') + ('N' << 8) + ('I' << 16))
	TAG_PDMD	= uint32 (('P') + ('D' << 8) + ('M' << 16) + ('D' << 24))
	TAG_X509	= uint32 (('X') + ('5' << 8) + ('0' << 16) + ('9' << 24))
	TAG_X59R	= uint32 (('X') + ('5' << 8) + ('9' << 16) + ('R' << 24))
	TAG_CCS		= uint32 (('C') + ('C' << 8) + ('S' << 16))
	TAG_CCRT	= uint32 (('C') + ('C' << 8) + ('R' << 16) + ('T' << 24))
	TAG_SCFG	= uint32 (('S') + ('C' << 8) + ('F' << 16) + ('G' << 24))
	TAG_SNO		= uint32 (('S') + ('N' << 8) + ('O' << 16))
	TAG_CRT		= uint32 (('C') + ('R' << 8) + ('T' << 16))
	TAG_PROF	= uint32 (('P') + ('R' << 8) + ('O' << 16) + ('F' << 24))
	TAG_SCID	= uint32 (('S') + ('C' << 8) + ('I' << 16) + ('D' << 24))
	TAG_KEXS	= uint32 (('K') + ('E' << 8) + ('X' << 16) + ('S' << 24))
	TAG_C255	= uint32 (('C') + ('2' << 8) + ('5' << 16) + ('5' << 24))
	TAG_P256	= uint32 (('p') + ('2' << 8) + ('5' << 16) + ('6' << 24))
	TAG_PUBS	= uint32 (('P') + ('U' << 8) + ('B' << 16) + ('S' << 24))
	TAG_AEAD	= uint32 (('A') + ('E' << 8) + ('A' << 16) + ('D' << 24))
	TAG_NULL	= uint32 (('N') + ('U' << 8) + ('L' << 16) + ('L' << 24))
	TAG_AESG	= uint32 (('A') + ('E' << 8) + ('S' << 16) + ('G' << 24))
	TAG_520P	= uint32 (('5') + ('2' << 8) + ('0' << 16) + ('P' << 24))
	TAG_ORBT	= uint32 (('O') + ('R' << 8) + ('B' << 16) + ('T' << 24))
	TAG_EXPY	= uint32 (('E') + ('X' << 8) + ('P' << 16) + ('Y' << 24))
	TAG_NONC	= uint32 (('N') + ('O' << 8) + ('N' << 16) + ('C' << 24))
	TAG_CETV	= uint32 (('C') + ('E' << 8) + ('T' << 16) + ('V' << 24))
	TAG_CIDK	= uint32 (('C') + ('I' << 8) + ('D' << 16) + ('K' << 24))
	TAG_CIDS	= uint32 (('C') + ('I' << 8) + ('D' << 16) + ('S' << 24))
	TAG_RREJ	= uint32 (('R') + ('R' << 8) + ('E' << 16) + ('J' << 24))
	TAG_CADR	= uint32 (('C') + ('A' << 8) + ('D' << 16) + ('R' << 24))
	TAG_RNON	= uint32 (('R') + ('N' << 8) + ('O' << 16) + ('N' << 24))
	TAG_RSEQ	= uint32 (('R') + ('S' << 8) + ('E' << 16) + ('Q' << 24))
	TAG_COPT	= uint32 (('C') + ('O' << 8) + ('P' << 16) + ('T' << 24))
	TAG_ICSL	= uint32 (('I') + ('C' << 8) + ('S' << 16) + ('L' << 24))
	TAG_SCLS	= uint32 (('S') + ('C' << 8) + ('L' << 16) + ('S' << 24))
	TAG_MSPC	= uint32 (('M') + ('S' << 8) + ('P' << 16) + ('C' << 24))
	TAG_IRTT	= uint32 (('I') + ('R' << 8) + ('T' << 16) + ('T' << 24))
	TAG_SWND	= uint32 (('S') + ('W' << 8) + ('N' << 16) + ('D' << 24))
	TAG_SFCW	= uint32 (('S') + ('F' << 8) + ('C' << 16) + ('W' << 24))
	TAG_CFCW	= uint32 (('C') + ('F' << 8) + ('C' << 16) + ('W' << 24))
)

type Message struct {
	tag uint32
	tags []uint32
	values [][]byte
}

func NewMessage (tag uint32) *Message {
	switch tag {
	case TAG_CHLO, TAG_REJ, TAG_SHLO, TAG_SCUP, TAG_PRST:
		return &Message { tag: tag }
	}
	return nil
}

func (this *Message) GetTag () uint32 {
	return this.tag
}

func (this *Message) EqualTag (tag uint32) bool {
	return this.tag == tag
}

func (this *Message) CountTags () int {
	return len (this.tags)
}

func (this *Message) ContainsTag (tag uint32) (bool, []byte) {
	for i, v := range this.tags {
		if v == tag {
			return true, this.values[i]
		}
	}
	return false, nil
}

func (this *Message) UpdateTag (tag uint32, value []byte) bool {
	for i, v := range this.tags {
		if v == tag {
			this.values[i] = value
			return true
		}
	}
	return false
}

func (this *Message) AppendTag (tag uint32, value []byte) bool {
	if res, _ := this.ContainsTag (tag); res {
		return false
	}

	this.tags = append (this.tags, tag)
	this.values = append (this.values, value)
	return true
}

func (this *Message) GetSerializedSize () uint32 {
	var size uint32 = 0

	for _, v := range this.values {
		size += uint32 (len (v))
	}
	return uint32 (len (this.tags)) * 8 + size + 8
}

func (this *Message) Serialize () []byte {
	message := make ([]byte, this.GetSerializedSize ())

	if len (this.tags) > 1 {
		haveSwap := true
		for haveSwap {
			haveSwap = false
			for i := 0; i < len (this.tags) - 1; i++ {
				if this.tags[i] > this.tags[i + 1] {
					t := this.tags[i]
					v := this.values[i]
					this.tags[i] = this.tags[i + 1]
					this.values[i] = this.values[i + 1]
					this.tags[i + 1] = t
					this.values[i + 1] = v

					haveSwap = true
				}
			}
		}
	}
	var off uint32 = 0
	binary.LittleEndian.PutUint32 (message, uint32 (this.tag))
	off += 4
	binary.LittleEndian.PutUint16 (message[off:], uint16 (len (this.tags)))
	off += 2
	// padding
	off += 2

	var endOff uint32 = 0
	// tags
	for i, t := range this.tags {
		binary.LittleEndian.PutUint32 (message[off:], uint32 (t))
		off += 4
		endOff += uint32 (len (this.values[i]))
		binary.LittleEndian.PutUint32 (message[off:], uint32 (endOff))
		off += 4
	}
	// values
	for _, v := range this.values {
		copy (message[off:], v)
		off += uint32 (len (v))
	}

	return message
}




