package frame

import (
	"errors"
	"bytes"
	"../protocol"
	"../utils"
)

type StopSendingFrame struct {
	Frame
	streamID	protocol.StreamID
	errorCode	uint16
}

func StopSendingFrameParse (b *bytes.Reader) (*StopSendingFrame, error) {
	frameType, err := b.ReadByte ()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_STOP_SENDING {
		return nil, errors.New ("StopSendingFrameParse error: frametype not equal 0x0C")
	}

	streamID, err := protocol.StreamIDParse (b)
	if err != nil {
		return nil, err
	}

	errorCode, err := utils.BigEndian.ReadUInt (b, 2)
	if err != nil {
		return nil, err
	}

	return &StopSendingFrame { Frame { frameType }, *streamID, uint16 (errorCode) }, nil
}

func (this *StopSendingFrame) Serialize (b *bytes.Buffer) error {
	err := b.WriteByte (this.frameType)
	if err != nil {
		return nil
	}

	err = this.streamID.Serialize (b)
	if err != nil {
		return err
	}
	utils.BigEndian.WriteUInt (b, uint64 (this.errorCode), 2)
	return nil
}