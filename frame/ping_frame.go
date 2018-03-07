package frame

import (
	"../utils"
	"bytes"
	"errors"
)

type PingFrame struct {
	Frame
	length	uint8
	data	utils.VarLenIntegerStruct
}

func PingFrameParse (b *bytes.Reader) (*PingFrame, error) {
	frameType, err := b.ReadByte ()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_PING {
		return nil, errors.New ("PingFrameParse error: frametype not equal 0x07")
	}

	len, err := b.ReadByte ()
	if err != nil {
		return nil, err
	}

	data, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	return &PingFrame { Frame { frameType }, len, *data }, nil
}

func (this *PingFrame) Serialize (b *bytes.Buffer) error {
	err := b.WriteByte (this.frameType)
	if err != nil {
		return err
	}

	err = b.WriteByte (this.length)
	if err != nil {
		return err
	}
	_, err = this.data.Serialize (b)
	return err
}
