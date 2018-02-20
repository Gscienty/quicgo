package protocol

import "errors"
import "encoding/binary"

type VersionNegotiationPacket struct {
	versions []uint32
}

func (this *VersionNegotiationPacket) ContainVersion (version uint32) bool {
	for _, v := range this.versions {
		if v == version {
			return true
		}
	}

	return false
}

func (this *VersionNegotiationPacket) AppendVersion (version uint32) bool {
	if this.ContainVersion (version) {
		return false
	}

	this.versions = append (this.versions, version)
	return true
}

func (this *VersionNegotiationPacket) GetSerializedSize () int {
	return len (this.versions) * 8
}

func (this *VersionNegotiationPacket) Serialize (data []byte) (size int, err error) {
	if this.GetSerializedSize () > len (data) {
		err = errors.New ("Version Negotiation Packet error: data too small")
		return
	}

	size = 0
	for _, v := range (this.versions) {
		binary.LittleEndian.PutUint32 (data[size:], v)
		size += 4
	}

	if size != this.GetSerializedSize () {
		err = errors.New ("Version Negotiation Packet error: internal error size different")
	}

	return
}

func (this *VersionNegotiationPacket) Parse (data []byte) (size int, err error) {
	size = 0
	if len (data) % 4 != 0 {
		err = errors.New ("Version Negotiation Packet error: parse data error")
		return
	}

	count := len (data) / 4
	for i := 0; i < count; i++ {
		this.AppendVersion (binary.LittleEndian.Uint32 (data[size:]))
		size += 4
	}

	return
}