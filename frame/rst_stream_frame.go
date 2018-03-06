package frame

import (
	"../utils"
	"../protocol"
	"bytes"
)

type RstStreamFrame struct {
	streamID	protocol.StreamID
	errorCode	uint16
	finalOffset	utils.VarLenIntegerStruct
}

func RstStreamFrameParse (b *bytes.Reader) (*RstStreamFrame, error) {
	streamID, err := protocol.StreamIDParse (b)
	if err != nil {
		return nil, err
	}
	errorCode, err := utils.BigEndian.ReadUInt (b, 2)
	if err != nil {
		return nil, err
	}
	finalOffset, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	return &RstStreamFrame { *streamID, uint16 (errorCode), *finalOffset }, nil
}

func (this *RstStreamFrame) Serialize (b *bytes.Buffer) error {
	this.streamID.Serialize (b)
	utils.BigEndian.WriteUInt (b, uint64 (this.errorCode), 2)
	this.finalOffset.Serialize (b)

	return nil
}