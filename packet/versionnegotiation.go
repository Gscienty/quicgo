package packet

import "errors"
import "encoding/binary"

type VersionNegotiationPacket struct {
	supportedVersions []uint32
}

func (this *VersionNegotiationPacket) Parse (data []byte) (size int, err error) {
	if len (data) % 4 != 0 {
		err = errors.New("VersionNegotiationPacket.Parse: data length error")
		return
	}
	size = 0
	for size < len (data) {
		this.supportedVersions = append (
			this.supportedVersions,
			binary.BigEndian.Uint32 (data[size:]))
		size += 4
	}
	return
}

func (this *VersionNegotiationPacket) GetSerializedSize () int {
	return len (this.supportedVersions) * 4
}

func (this *VersionNegotiationPacket) Serialize (data []byte) (size int, err error) {
	if len (data) < this.GetSerializedSize () {
		err = errors.New("VersionNegotiationPacket.Parse: data length error")
		return
	}
	size = 0
	for _, v := range this.supportedVersions {
		binary.BigEndian.PutUint32 (data[size:], v)
		size += 4
	}
	return
}