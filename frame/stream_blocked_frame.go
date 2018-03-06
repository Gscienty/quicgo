package frame

import (
	"../utils"
	"../protocol"
	"bytes"
)

type StreamBlockedFrame struct {
	streamID	protocol.StreamID
	offset		utils.VarLenIntegerStruct
}

func StreamBlockedFrameParse (b *bytes.Reader) (*StreamBlockedFrame, error) {
	streamID, err := protocol.StreamIDParse (b)
	if err != nil {
		return nil, err
	}

	offset, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	return &StreamBlockedFrame { *streamID, *offset }, nil
}

func (this *StreamBlockedFrame) Serialize (b *bytes.Buffer) error {
	err := this.streamID.Serialize (b)
	if err != nil {
		return err
	}
	_, err = this.offset.Serialize (b)
	return err
}