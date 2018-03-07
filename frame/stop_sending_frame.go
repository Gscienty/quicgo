package frame

import (
	"bytes"
	"../protocol"
	"../utils"
)

type StopSendingFrame struct {
	streamID	protocol.StreamID
	errorCode	uint16
}

func StopSendingFrameParse (b *bytes.Reader) (*StopSendingFrame, error) {
	streamID, err := protocol.StreamIDParse (b)
	if err != nil {
		return nil, err
	}

	errorCode, err := utils.BigEndian.ReadUInt (b, 2)
	if err != nil {
		return nil, err
	}

	return &StopSendingFrame { *streamID, uint16 (errorCode) }, nil
}

func (this *StopSendingFrame) Serialize (b *bytes.Buffer) error {
	err := this.streamID.Serialize (b)
	if err != nil {
		return err
	}
	utils.BigEndian.WriteUInt (b, uint64 (this.errorCode), 2)
	return nil
}