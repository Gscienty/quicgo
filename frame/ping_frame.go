package frame

import (
	"../utils"
	"bytes"
)

type PingFrame struct {
	length	uint8
	data	utils.VarLenIntegerStruct
}

func PingFrameParse (b *bytes.Reader) (*PingFrame, error) {
	len, err := b.ReadByte ()
	if err != nil {
		return nil, err
	}

	data, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	return &PingFrame { len, *data }, nil
}

func (this *PingFrame) Serialize (b *bytes.Buffer) error {
	err := b.WriteByte (this.length)
	if err != nil {
		return err
	}
	_, err = this.data.Serialize (b)
	if err != nil {
		return err
	}

	return nil
}
