package frame

import(
	"../utils"
	"bytes"
	"errors"
)

type BlockedFrame struct {
	Frame
	offset utils.VarLenIntegerStruct
}

func BlockedFrameParse(b *bytes.Reader) (*BlockedFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) != FRAME_TYPE_BLOCKED {
		return nil, errors.New("BlockedFrameParse error: frametype not equal 0x08")
	}

	offset, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}

	return &BlockedFrame { Frame { FrameType(frameType) }, *offset }, nil
}

func (this *BlockedFrame) GetType() FrameType {
	return FRAME_TYPE_BLOCKED
}

func (this *BlockedFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}
	_, err = this.offset.Serialize(b)
	return err
}