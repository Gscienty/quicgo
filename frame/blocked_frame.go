package frame

import (
	"../utils"
	"bytes"
)

type BlockedFrame struct {
	offset utils.VarLenIntegerStruct
}

func BlockedFrameParse (b *bytes.Reader) (*BlockedFrame, error) {
	offset, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	return &BlockedFrame { *offset }, nil
}

func (this *BlockedFrame) Serialize (b *bytes.Buffer) error {
	_, err := this.offset.Serialize (b)
	return err
}