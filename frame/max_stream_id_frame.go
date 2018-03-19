package frame

import(
	"../protocol"
	"bytes"
	"errors"
)

type MaxStreamIDFrame struct {
	Frame
	streamID protocol.StreamID
}

func MaxStreamIDFrameParse(b *bytes.Reader) (*MaxStreamIDFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) != FRAME_TYPE_MAX_STREAM_ID {
		return nil, errors.New("MaxStreamIDFrameParse error: frametype not equal 0x06")
	}

	streamID, err := protocol.StreamIDParse(b)
	if err != nil {
		return nil, err
	}

	return &MaxStreamIDFrame { Frame { FrameType(frameType) }, *streamID }, nil
}

func (this *MaxStreamIDFrame) GetType() FrameType {
	return FRAME_TYPE_MAX_STREAM_ID
}

func (this *MaxStreamIDFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}

	return this.streamID.Serialize(b)
}