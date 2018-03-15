package handshake

import (
	"encoding/binary"
	"errors"
	"bytes"
	"time"
	"../protocol"
	"../utils"
)

type transportParameterID uint16

const (
	initialMaxStreamDataParameterID		transportParameterID = 0x00
	initialMaxDataParameterID			transportParameterID = 0x01
	initialMaxStreamIDBiDiParameterID	transportParameterID = 0x02
	idleTimeoutParameterID				transportParameterID = 0x03
	omitConnectionIDParameterID			transportParameterID = 0x04
	maxPacketSizeParameterID			transportParameterID = 0x05
	statelessResetTokenParameterID		transportParameterID = 0x06
	initialMaxStreamIDUniParameterID	transportParameterID = 0x08
)

type TransportParameter struct {
	Parameter	transportParameterID
	Value		[]byte
}

func (this *TransportParameter) Serialize(b *bytes.Buffer) error {
	utils.BigEndian.WriteUInt(b, uint64(this.Parameter), 2)
	utils.BigEndian.WriteUInt(b, uint64(len(this.Value)), 2)
	writedLen, err := b.Write(this.Value)
	if err != nil {
		return err
	}
	if writedLen != len(this.Value) {
		return errors.New("writed length not equal")
	}
	return nil
}

func (this *TransportParameter) Parse(b *bytes.Reader) (int, error) {
	parameterID, err := utils.BigEndian.ReadUInt(b, 2)
	if err != nil {
		return 0, err
	}
	valueLength, err := utils.BigEndian.ReadUInt(b, 2)
	if err != nil {
		return 0, err
	}

	value := make([]byte, valueLength)
	readedValueLength, err := b.Read(value)
	if err != nil {
		return 0, err
	}
	if readedValueLength != int(valueLength) {
		return 0, errors.New("value length error")
	}

	this.Parameter = transportParameterID(parameterID)
	this.Value = value

	return 4 + int(valueLength), nil
}

type TransportParameters struct {
	StreamFlowControlWindow		uint64
	ConnectionFlowControlWindow	uint64

	MaxBidiStreamID				protocol.StreamID
	MaxUniStreamID				protocol.StreamID
	MaxStreams					uint32

	OmitConnectionID			bool
	IdleTimeout					time.Duration
}

func transportParameterToTransportParameters(params []TransportParameter) (*TransportParameters, error) {
	ret := &TransportParameters { }

	var existInitialMaxStreamData bool
	var existInitialMaxData bool
	var existIdleTimeout bool

	for _, p := range params {
		switch p.Parameter {
		case initialMaxStreamDataParameterID:
			existInitialMaxStreamData = true
			if len(p.Value) != 4 {
				return nil, errors.New("wrong length")
			}
			ret.StreamFlowControlWindow = uint64(binary.BigEndian.Uint32(p.Value))
		case initialMaxDataParameterID:
			existInitialMaxData = true
			if len(p.Value) != 4 {
				return nil, errors.New("wrong length")
			}
			ret.ConnectionFlowControlWindow = uint64(binary.BigEndian.Uint32(p.Value))
		case initialMaxStreamIDBiDiParameterID:
			if len(p.Value) != 4 {
				return nil, errors.New("wrong length")
			}
			ret.MaxBidiStreamID = *protocol.StreamIDNew(uint64(binary.BigEndian.Uint32(p.Value)))
		case initialMaxStreamIDUniParameterID:
			if len(p.Value) != 4 {
				return nil, errors.New("wrong length")
			}
			ret.MaxUniStreamID = *protocol.StreamIDNew(uint64(binary.BigEndian.Uint32(p.Value)))
		case idleTimeoutParameterID:
			existIdleTimeout = true
			if len(p.Value) != 2 {
				return nil, errors.New("wrong length")
			}
			ret.IdleTimeout = protocol.MIN_REMOTE_IDLE_TIMEOUT
			t := time.Duration(binary.BigEndian.Uint16(p.Value)) * time.Second
			if t > ret.IdleTimeout {
				ret.IdleTimeout = t
			}
		case omitConnectionIDParameterID:
			if len(p.Value) != 0 {
				return nil, errors.New("wrong length")
			}
			ret.OmitConnectionID = true
		}
	}

	if !(existInitialMaxStreamData && existInitialMaxData && existIdleTimeout) {
		return nil, errors.New("missing parameter")
	}
	return ret, nil
}

func helloMapTransferToTransportParameters(tags map[HandshakeTag][]byte) (*TransportParameters, error) {
	ret := &TransportParameters { }
	if value, ok := tags[TAG_TCID]; ok {
		v, err := utils.LittleEndian.ReadUInt(bytes.NewBuffer(value), 4)
		if err != nil {
			return nil, err
		}
		ret.OmitConnectionID = (v == 0)
	}
	if value, ok := tags[TAG_MIDS]; ok {
		v, err := utils.LittleEndian.ReadUInt(bytes.NewBuffer(value), 4)
		if err != nil {
			return nil, err
		}
		ret.MaxStreams = uint32(v)
	}
	if value, ok := tags[TAG_ICSL]; ok {
		v, err := utils.LittleEndian.ReadUInt(bytes.NewBuffer(value), 4)
		if err != nil {
			return nil, err
		}
		ret.IdleTimeout = protocol.MIN_REMOTE_IDLE_TIMEOUT
		t := time.Duration(v) * time.Second
		if t > ret.IdleTimeout {
			ret.IdleTimeout = t
		}
	}
	if value, ok := tags[TAG_SFCW]; ok {
		v, err := utils.LittleEndian.ReadUInt(bytes.NewBuffer(value), 4)
		if err != nil {
			return nil, err
		}
		ret.StreamFlowControlWindow = v
	}
	if value, ok := tags[TAG_CFCW]; ok {
		v, err := utils.LittleEndian.ReadUInt(bytes.NewBuffer(value), 4)
		if err != nil {
			return nil, err
		}
		ret.ConnectionFlowControlWindow = v
	}
	return ret, nil
}

func (this *TransportParameters) transferToHelloMap() map[HandshakeTag][]byte {
	sfcw := bytes.NewBuffer([]byte { })
	utils.LittleEndian.WriteUInt(sfcw, uint64(this.StreamFlowControlWindow), 4)
	cfcw := bytes.NewBuffer([]byte { })
	utils.LittleEndian.WriteUInt(cfcw, uint64(this.ConnectionFlowControlWindow), 4)
	mids := bytes.NewBuffer([]byte { })
	utils.LittleEndian.WriteUInt(mids, uint64(this.MaxStreams), 4)
	icsl := bytes.NewBuffer([]byte { })
	utils.LittleEndian.WriteUInt(icsl, uint64(this.IdleTimeout / time.Second), 4)

	tags := map[HandshakeTag][]byte {
		TAG_ICSL: icsl.Bytes(),
		TAG_MIDS: mids.Bytes(),
		TAG_CFCW: cfcw.Bytes(),
		TAG_SFCW: sfcw.Bytes(),
	}
	if this.OmitConnectionID {
		tags[TAG_TCID] = []byte { 0, 0, 0, 0 }
	}
	return tags
}

func (this *TransportParameters) transferTransportParameter() []TransportParameter {
	initialMaxStreamData := make([]byte, 4)
	binary.BigEndian.PutUint32(initialMaxStreamData, uint32(this.StreamFlowControlWindow))
	initialMaxData := make([]byte, 4)
	binary.BigEndian.PutUint32(initialMaxData, uint32(this.ConnectionFlowControlWindow))
	initialBiDIStreamID := make([]byte, 4)
	binary.BigEndian.PutUint32(initialBiDIStreamID, uint32(this.MaxBidiStreamID.GetID()))
	initialUniStreamID := make([]byte, 4)
	binary.BigEndian.PutUint32(initialUniStreamID, uint32(this.MaxUniStreamID.GetID()))
	idleTimeout := make([]byte, 2)
	binary.BigEndian.PutUint16(idleTimeout, uint16(this.IdleTimeout / time.Second))
	maxPacketSize := make([]byte, 2)
	binary.BigEndian.PutUint16(maxPacketSize, uint16(protocol.MAX_RECEIVE_PACKET_SIZE))

	ret := []TransportParameter {
		{ initialMaxStreamDataParameterID, initialMaxStreamData },
		{ initialMaxDataParameterID, initialMaxData },
		{ initialMaxStreamIDBiDiParameterID, initialBiDIStreamID },
		{ initialMaxStreamIDUniParameterID, initialUniStreamID },
		{ idleTimeoutParameterID, idleTimeout },
		{ maxPacketSizeParameterID, maxPacketSize },
	}
	if this.OmitConnectionID {
		ret = append(ret, TransportParameter { omitConnectionIDParameterID, []byte { } })
	}
	return ret
}