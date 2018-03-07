package frame

import (
	"../utils"
	"../protocol"
	"bytes"
	"errors"
)

type RstStreamFrame struct {
	Frame
	streamID	protocol.StreamID
	errorCode	uint16
	finalOffset	utils.VarLenIntegerStruct
}

func RstStreamFrameParse (b *bytes.Reader) (*RstStreamFrame, error) {
	frameType,err := b.ReadByte ()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_RST_STREAM {
		return nil, errors.New ("RstStreamFrameParse error: frametype not equal 0x01")
	}

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

	return &RstStreamFrame { Frame { frameType }, *streamID, uint16 (errorCode), *finalOffset }, nil
}

func (this *RstStreamFrame) Serialize (b *bytes.Buffer) error {
	err := b.WriteByte (this.frameType)
	if err != nil {
		return err
	}

	this.streamID.Serialize (b)
	utils.BigEndian.WriteUInt (b, uint64 (this.errorCode), 2)
	this.finalOffset.Serialize (b)

	return nil
}