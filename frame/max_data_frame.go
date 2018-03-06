package frame

import (
	"bytes"
	"../utils"
)

type MaxDataFrame struct {
	maximumData utils.VarLenIntegerStruct
}

func MaxDataFrameParse (b *bytes.Reader) (*MaxDataFrame, error) {
	maximumData, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}
	return &MaxDataFrame { maximumData: *maximumData }, nil
}

func (this *MaxDataFrame) Serialize (b *bytes.Buffer) error {
	_, err := this.maximumData.Serialize (b)
	if err != nil {
		return err
	}
	return nil
}