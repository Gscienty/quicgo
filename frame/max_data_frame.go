package frame

import(
	"errors"
	"bytes"
	"../utils"
)

type MaxDataFrame struct {
	Frame
	maximumData utils.VarLenIntegerStruct
}

func MaxDataFrameParse(b *bytes.Reader) (*MaxDataFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) != FRAME_TYPE_MAX_DATA {
		return nil, errors.New("MaxDataFrameParse error: frametype not equal 0x04")
	}

	maximumData, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}
	return &MaxDataFrame { Frame { FrameType(frameType) }, *maximumData }, nil
}

func (this *MaxDataFrame) GetType() FrameType {
	return FRAME_TYPE_MAX_DATA
}

func (this *MaxDataFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}

	_, err = this.maximumData.Serialize(b)
	return err
}