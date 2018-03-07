package frame

import (
	"errors"
	"bytes"
	"../utils"
	"../protocol"
)

type NewConnectionIDFrame struct {
	Frame
	sequence		utils.VarLenIntegerStruct
	connectionID	protocol.ConnectionID
	token			[]byte
}

func NewConnectionIDFrameParse (b *bytes.Reader) (*NewConnectionIDFrame, error) {
	frameType, err := b.ReadByte ()
	if err != nil {
		return nil, err
	}
	if frameType != 0x0B {
		return nil, errors.New ("NewConnectionIDFrameParse error: frametype not equal 0x0B")
	}

	sequence, err := utils.VarLenIntegerStructParse (b)
	if err != nil {
		return nil, err
	}

	connectionID, err := utils.BigEndian.ReadUInt (b, 8)
	if err != nil {
		return nil, err
	}

	token := make ([]byte, 16)
	len, err := b.Read (token)
	if err != nil {
		return nil, err
	}
	if len != 16 {
		return nil, errors.New ("NewConnectionIDFrameParse error: cannot read fully token")
	}

	return &NewConnectionIDFrame { Frame { frameType }, *sequence, protocol.ConnectionID (connectionID), token }, nil
}

func (this *NewConnectionIDFrame) Serialize (b *bytes.Buffer) error {
	err := b.WriteByte (this.frameType)
	if err != nil {
		return err
	}

	_, err = this.sequence.Serialize (b)
	if err != nil {
		return err
	}

	utils.BigEndian.WriteUInt (b, uint64 (this.connectionID), 8)

	len, err := b.Write (this.token)
	if err != nil {
		return err
	}
	if len != 16 {
		return errors.New ("NewConnectionIDFrame.Serialize error: token length error")
	}
	return nil
}