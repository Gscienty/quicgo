package frame

import(
	"../utils"
	"../protocol"
	"bytes"
	"errors"
)

type StreamBlockedFrame struct {
	Frame
	streamID	protocol.StreamID
	offset		utils.VarLenIntegerStruct
}

func StreamBlockedFrameParse(b *bytes.Reader) (*StreamBlockedFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_STREAM_BLOCKED {
		return nil, errors.New("StreamBlockedFrameParse error: frametype not equal 0x09")
	}

	streamID, err := protocol.StreamIDParse(b)
	if err != nil {
		return nil, err
	}

	offset, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}

	return &StreamBlockedFrame { Frame { frameType }, *streamID, *offset }, nil
}

func (this *StreamBlockedFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(this.frameType)
	if err != nil {
		return err
	}

	err = this.streamID.Serialize(b)
	if err != nil {
		return err
	}
	_, err = this.offset.Serialize(b)
	return err
}