package protocol

import(
	"../utils"
	"errors"
	"bytes"
	"io"
)

const(
	STREAM_PERSPECTIVE_MASK		= 0x01
	STREAM_PERSPECTIVE_CLIENT 	= 0x00
	STREAM_PERSPECTIVE_SERVER 	= 0x01
	STREAM_TYPE_MASK			= 0x02
	STREAM_TYPE_BIDIRECTIONAL	= 0x00
	STREAM_TYPE_UNBIDIECTIONAL	= 0x02
)

type StreamID struct {
	perspective	uint8
	streamType	uint8
	id			uint64
}

func StreamIDNew(val uint64) *StreamID {
	return &StreamID {
		id: val >> 2,
		perspective: uint8 (val & STREAM_PERSPECTIVE_MASK),
		streamType: uint8 (val & STREAM_TYPE_MASK),
	}
}

func StreamIDParse(b io.Reader) (*StreamID, error) {
	sid, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}
	if sid.GetVal() >((uint64(1) << 62) - 1) {
		return nil, errors.New("StreamIDParse error: val too large")
	}

	return &StreamID {
		id: sid.GetVal() >> 2,
		perspective: uint8(sid.GetVal() & STREAM_PERSPECTIVE_MASK),
		streamType: uint8(sid.GetVal() & STREAM_TYPE_MASK),
	}, nil
}

func (this *StreamID) Serialize(b *bytes.Buffer) error {
	val :=(this.id << 2) | uint64(this.perspective) | uint64(this.streamType)
	utils.VarLenIntegerStructNew(val).Serialize(b)
	return nil
}

func (this *StreamID) SetID(id uint64) error {
	if id >((uint64(1) << 60) - 1) {
		return errors.New("StreamID.SetID error: id too large")
	}
	this.id = id
	return nil
}

func (this *StreamID) GetID() uint64 {
	return this.id
}

func (this *StreamID) SetPerspective(perspective uint8) error {
	if perspective != STREAM_PERSPECTIVE_CLIENT && perspective != STREAM_PERSPECTIVE_SERVER {
		return errors.New("StreamID.SetPerspective error: illegal perspective")
	}
	this.perspective = perspective
	return nil
}

func (this *StreamID) GetPerspective() uint8 {
	return this.perspective
}

func (this *StreamID) SetType(streamType uint8) error {
	if streamType != STREAM_TYPE_BIDIRECTIONAL && streamType != STREAM_TYPE_UNBIDIECTIONAL {
		return errors.New("StreamID.SetType error: illegal type")
	}
	this.streamType = streamType
	return nil
}

func (this *StreamID) GetType() uint8 {
	return this.streamType
}

func (this *StreamID) Equal(other *StreamID) bool {
	return this.id == other.id && this.perspective == other.perspective && this.streamType == other.streamType
}
