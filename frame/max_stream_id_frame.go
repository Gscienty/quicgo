package frame

import (
	"../protocol"
	"bytes"
)

type MaxStreamIDFrame struct {
	streamID protocol.StreamID
}

func MaxStreamIDFrameParse (b *bytes.Reader) (*MaxStreamIDFrame, error) {
	streamID, err := protocol.StreamIDParse (b)
	if err != nil {
		return nil, err
	}

	return &MaxStreamIDFrame { *streamID }, nil
}

func (this *MaxStreamIDFrame) Serialize (b *bytes.Buffer) error {
	return this.streamID.Serialize (b)
}