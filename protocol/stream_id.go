package protocol

import (
	"io"
	"errors"
	"encoding/binary"
	"bytes"
)

const (
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

func (StreamID) Parse (b io.Reader) (*StreamID, error) {
	buf := make([]byte, 4)
	l, err := b.Read (buf)
	if err != nil {
		return nil, err
	}
	if l != 4 {
		return nil, errors.New ("StreamID.Parse error: len error")
	}
	var retval *StreamID = &StreamID { }
	if (buf[3] & STREAM_PERSPECTIVE_MASK) == STREAM_PERSPECTIVE_CLIENT {
		retval.perspective = STREAM_PERSPECTIVE_CLIENT
	} else {
		retval.perspective = STREAM_PERSPECTIVE_SERVER
	}
	if (buf[3] & STREAM_TYPE_MASK) == STREAM_TYPE_BIDIRECTIONAL {
		retval.streamType = STREAM_TYPE_BIDIRECTIONAL
	} else {
		retval.streamType = STREAM_TYPE_UNBIDIECTIONAL
	}
	buf[3] &= 0xFC
	retval.id = binary.BigEndian.Uint64 (buf) >> 2
	return retval, nil
}

func (this *StreamID) Serialize (b bytes.Buffer) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint64 (buf, this.id)
	buf[3] &= (this.perspective | this.streamType)
	_, err := b.Write (buf)
	if err != nil {
		return err
	}
	return nil
}

func (this *StreamID) SetID (id uint64) error {
	if id > 0xFFFFFFFFFFFFFC {
		return errors.New ("StreamID.SetID error: id too large")
	}
	this.id = id
	return nil
}

func (this *StreamID) SetPerspective (perspective uint8) error {
	if perspective != STREAM_PERSPECTIVE_CLIENT && perspective != STREAM_PERSPECTIVE_SERVER {
		return errors.New ("StreamID.SetPerspective error: illegal perspective")
	}
	this.perspective = perspective
	return nil
}

func (this *StreamID) SetType (streamType uint8) error {
	if streamType != STREAM_TYPE_BIDIRECTIONAL && streamType != STREAM_TYPE_UNBIDIECTIONAL {
		return errors.New ("StreamID.SetType error: illegal type")
	}
	this.streamType = streamType
	return nil
}

