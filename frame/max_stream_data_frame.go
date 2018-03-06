package frame

import (
	"bytes"
	"../utils"
	"../protocol"
)

type MaxStreamDataFrame struct {
	streamID	protocol.StreamID
	maximumData	utils.VarLenIntegerStruct
}

func MaxStreamDataFrameParse (b *bytes.Reader) (*MaxStreamDataFrame, error) {
	streamID, err := protocol.StreamIDParse (b)
	if err != nil {
		return nil, err
	}

	maximumData, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	return &MaxStreamDataFrame { *streamID, *maximumData }, nil
}

func (this *MaxStreamDataFrame) Serialize (b *bytes.Buffer) error {
	err := this.streamID.Serialize (b)
	if err != nil {
		return err
	}
	_, err = this.maximumData.Serialize (b)
	if err != nil {
		return err
	}

	return nil
}