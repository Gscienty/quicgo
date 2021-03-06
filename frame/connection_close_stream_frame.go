package frame

import(
	"bytes"
	"errors"
	"../utils"
	"encoding/binary"
)

type ConnectionCloseStreamFrame struct {
	Frame
	errorCode	uint16
	reason		string
}

func ConnectionCloseStreamFrameParse(b *bytes.Reader) (*ConnectionCloseStreamFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) != FRAME_TYPE_CONNECTION_CLOSE {
		return nil, errors.New("ConnectionCloseStreamFrameParse error: frametype not equal 0x02")
	}

	errcodeBuf := make([]byte, 2)
	readedLen, err := b.Read(errcodeBuf)
	if err != nil {
		return nil, err
	}
	if readedLen != 2 {
		return nil, errors.New("ConnectionCloseStreamFrameParse error: buffer too small")
	}
	errCode := binary.BigEndian.Uint16(errcodeBuf)

	reasonLen, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}

	reasonBuf := make([]byte, reasonLen.GetVal())
	len, err := b.Read(reasonBuf)
	if err != nil {
		return nil, err
	}
	if len != int(reasonLen.GetVal()) {
		return nil, errors.New("ConnectionCloseStreamFrameParse error: reason length error")
	}

	return &ConnectionCloseStreamFrame { Frame { FrameType(frameType) }, errCode, string(reasonBuf) }, nil
}

func (this *ConnectionCloseStreamFrame) GetType() FrameType {
	return FRAME_TYPE_CONNECTION_CLOSE
}

func (this *ConnectionCloseStreamFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}
	
	utils.BigEndian.WriteUInt(b, uint64(this.errorCode), 2)
	reasonBuf := []byte(this.reason)
	reasonLen := utils.VarLenIntegerStructNew(uint64(len(reasonBuf)))
	_, err = reasonLen.Serialize(b)
	if err != nil {
		return err
	}
	_, err = b.Write(reasonBuf)
	return err
}