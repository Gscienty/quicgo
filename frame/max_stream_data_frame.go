package frame

import(
	"errors"
	"bytes"
	"../utils"
	"../protocol"
)

type MaxStreamDataFrame struct {
	Frame
	streamID	protocol.StreamID
	maximumData	utils.VarLenIntegerStruct
}

func MaxStreamDataFrameParse(b *bytes.Reader) (*MaxStreamDataFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_MAX_STREAM_DATA {
		return nil, errors.New("MaxStreamDataFrameParse error: frametype not equal 0x05")
	}

	streamID, err := protocol.StreamIDParse(b)
	if err != nil {
		return nil, err
	}

	maximumData, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}

	return &MaxStreamDataFrame { Frame { frameType }, *streamID, *maximumData }, nil
}

func (this *MaxStreamDataFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(this.frameType)
	if err != nil {
		return err
	}
	
	err = this.streamID.Serialize(b)
	if err != nil {
		return err
	}
	_, err = this.maximumData.Serialize(b)
	return err
}