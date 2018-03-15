package handshake

import (
	"encoding/binary"
	"bytes"
	"../protocol"
	"../utils"
)

type encryptedExtension struct {
	NegotiatedVersion	protocol.Version
	SupportedVersions	[]protocol.Version
	Parameters			[]TransportParameter
}

func (this *encryptedExtension) Serialize(b *bytes.Buffer) error {
	utils.BigEndian.WriteUInt(b, uint64(this.NegotiatedVersion), 4)
	
	svBuffer := bytes.NewBuffer([]byte { })
	for _, v := range this.SupportedVersions {
		utils.BigEndian.WriteUInt(svBuffer, uint64(v), 4)
	}
	utils.BigEndian.WriteUInt(b, uint64(svBuffer.Len()), 1)
	_, err := b.Write(svBuffer.Bytes())
	if err != nil {
		return err
	}

	parametersBuffer := bytes.NewBuffer([]byte { })
	for _, p := range this.Parameters {
		err := p.Serialize(parametersBuffer)
		if err != nil {
			return err
		}
	}
	utils.BigEndian.WriteUInt(b, uint64(parametersBuffer.Len()), 2)
	_, err = b.Write(parametersBuffer.Bytes())
	return err
}

func (this *encryptedExtension) Parse(b *bytes.Reader) (int, error) {
	version, err := utils.BigEndian.ReadUInt(b, 4)
	if err != nil {
		return 0, err
	}
	svSize, err := utils.BigEndian.ReadUInt(b, 1)
	if err != nil {
		return 0, err
	}
	svBuf := make([]byte, svSize)
	_, err = b.Read(svBuf)
	if err != nil {
		return 0, err
	}
	var supportedVersions []protocol.Version
	var vcount int = int(svSize) / 4
	for i := 0; i < vcount; i++ {
		supportedVersions = append(supportedVersions, protocol.Version(binary.BigEndian.Uint32(svBuf[i * 4:])))
	}
	paramSize, err := utils.BigEndian.ReadUInt(b, 2)
	paramBuf := make([]byte, paramSize)
	_, err = b.Read(paramBuf)
	if err != nil {
		return 0, err
	}
	var parameters []TransportParameter
	paramReader := bytes.NewReader(paramBuf)
	var off int = 0
	for off < int(paramSize) {
		parameter := &TransportParameter { }
		size, err := parameter.Parse(paramReader)
		if err != nil {
			return 0, err
		}
		off += size
		parameters = append(parameters, *parameter)
	}
	
	this.NegotiatedVersion = protocol.Version(version)
	this.SupportedVersions = supportedVersions
	this.Parameters = parameters

	return int(7 + paramSize + svSize), nil
}