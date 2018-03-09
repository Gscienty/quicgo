package frame

import(
	"errors"
	"bytes"
)

type PathResponseFrame struct {
	Frame
	data uint8
}

func PathResponseFrameParse(b *bytes.Reader) (*PathResponseFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_PATH_RESPONSE {
		return nil, errors.New("PathResponseFrameParse error: frametype not equal 0x0e")
	}
	data, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	return &PathResponseFrame { Frame { frameType }, uint8(data) }, nil
}

func (this *PathResponseFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(this.frameType)
	if err != nil {
		return err
	}
	return b.WriteByte(this.data)
}