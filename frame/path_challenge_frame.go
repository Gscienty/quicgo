package frame

import(
	"errors"
	"bytes"
)

type PathChallengeFrame struct {
	Frame
	data uint8
}

func PathChallengeFrameParse(b *bytes.Reader) (*PathChallengeFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) != FRAME_TYPE_PATH_CHALLENGE {
		return nil, errors.New("PathChallengeFrameParse error: frametype not equal 0x0e")
	}
	data, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	return &PathChallengeFrame { Frame { FrameType(frameType) }, uint8(data) }, nil
}

func (this *PathChallengeFrame) GetType() FrameType {
	return FRAME_TYPE_PATH_CHALLENGE
}

func (this *PathChallengeFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}
	return b.WriteByte(this.data)
}