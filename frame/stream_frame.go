package frame

import(
	"bytes"
	"errors"
	"io"
	"../protocol"
	"../utils"
)

type StreamFrame struct {
	Frame
	streamID	protocol.StreamID
	offset		utils.VarLenIntegerStruct
	length		utils.VarLenIntegerStruct

	data		[]byte
}

func StreamFrameParse(b *bytes.Reader) (*StreamFrame, error) {
	ret := &StreamFrame { }

	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) & FRAME_TYPE_STREAM != FRAME_TYPE_STREAM {
		return nil, errors.New("StreamFrameParse error: frametype error")
	}
	ret.frameType = FrameType(frameType)

	streamID, err := protocol.StreamIDParse(b)
	if err != nil {
		return nil, err
	}
	ret.streamID = *streamID

	if ret.frameType & FRAME_TYPE_STREAM_OFF == FRAME_TYPE_STREAM_OFF {
		offset, err := utils.VarLenIntegerStructParse(b)
		if err != nil {
			return nil, err
		}
		ret.offset = *offset
	}

	if ret.frameType & FRAME_TYPE_STREAM_LEN == FRAME_TYPE_STREAM_LEN {
		length, err := utils.VarLenIntegerStructParse(b)
		if err != nil {
			return nil, err
		}
		ret.length = *length
		if length.GetVal() > uint64(b.Len()) {
			return nil, io.EOF
		}
	} else {
		ret.length = *utils.VarLenIntegerStructNew(uint64(b.Len()))
	}

	if ret.length.GetVal() != 0 {
		ret.data = make([]byte, ret.length.GetVal())
		if _, err := io.ReadFull(b, ret.data); err != nil {
			return nil, err
		}
	}

	if(ret.frameType & FRAME_TYPE_STREAM_FIN != FRAME_TYPE_STREAM_FIN) && len(ret.data) == 0 {
		return nil, errors.New("StreamFrameParse error: empty stream")
	}

	return ret, nil
}

func (this *StreamFrame) GetType() FrameType {
	return FRAME_TYPE_STREAM
}

func (this *StreamFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}

	err = this.streamID.Serialize(b)
	if err != nil {
		return err
	}

	if this.frameType & FRAME_TYPE_STREAM_OFF == FRAME_TYPE_STREAM_OFF {
		_, err = this.offset.Serialize(b)
		if err != nil {
			return err
		}
	}

	if this.frameType & FRAME_TYPE_STREAM_LEN == FRAME_TYPE_STREAM_LEN {
		_, err = this.offset.Serialize(b)
		if err != nil {
			return err
		}
	}

	_, err = b.Write(this.data)
	return err
}
