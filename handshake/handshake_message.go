package handshake

import (
	"sort"
	"encoding/binary"
	"bytes"
	"../protocol"
	"errors"
	"../utils"
)

type HandshakeMessage struct {
	Tag		HandshakeTag
	Data	map[HandshakeTag][]byte
}

func HandshakeMessageParse(b *bytes.Reader) (HandshakeMessage, error) {
	buf := make([]byte, 4)
	if len, err := b.Read(buf); err != nil && len != 4 {
		return HandshakeMessage { }, err
	}

	tag := HandshakeTag(binary.LittleEndian.Uint32(buf))

	if len, err := b.Read(buf); err != nil && len != 4 {
		return HandshakeMessage { }, err
	}
	
	nPairs := binary.LittleEndian.Uint32(buf)
	if nPairs > protocol.CRYPTO_MAX_PARAMS {
		return HandshakeMessage { }, errors.New("Crypto Too Many Entries")
	}

	index := make([]byte, nPairs * 8)
	if len, err := b.Read(index); err != nil && uint32(len) != nPairs * 8 {
		return HandshakeMessage { }, err
	}

	resultMap := map[HandshakeTag][]byte { }

	var dataStart uint32
	for indexPos := 0; indexPos < int(nPairs) * 8; indexPos += 8 {
		tag := HandshakeTag(binary.LittleEndian.Uint32(index[indexPos : indexPos + 4]))
		dataEnd := binary.LittleEndian.Uint32(index[indexPos + 4 : indexPos + 8])

		dataLen := dataEnd - dataStart
		if dataLen > protocol.CRYPTO_PARAMETER_MAX_LENGTH {
			return HandshakeMessage { }, errors.New("value too lang")
		}

		data := make([]byte, dataLen)
		if len, err := b.Read(data); err != nil && uint32(len) != dataLen {
			return HandshakeMessage { }, err
		}

		resultMap[tag] = data
		dataStart = dataEnd
	}

	return HandshakeMessage { tag, resultMap }, nil
}

func (this *HandshakeMessage) sortedHandshakeTags() []HandshakeTag {
	tags := make([]HandshakeTag, len(this.Data))
	i := 0
	for t := range this.Data {
		tags[i] = t
		i++
	}

	sort.Slice(tags, func(i int, j int) bool {
		return tags[i] < tags[j]
	})

	return tags
}

func (this HandshakeMessage) Serialize(b *bytes.Buffer) {
	utils.LittleEndian.WriteUInt(b, uint64(this.Tag), 4)
	utils.LittleEndian.WriteUInt(b, uint64(len(this.Data)), 2)
	utils.LittleEndian.WriteUInt(b, 0, 2)

	indexStart := b.Len()
	indexData := make([]byte, 8 * len(this.Data))
	b.Write(indexData)

	offset := uint32(0)
	for i, t := range this.sortedHandshakeTags() {
		v := this.Data[t]
		b.Write(v)
		offset += uint32(len(v))
		binary.LittleEndian.PutUint32(indexData[i * 8:], uint32(t))
		binary.LittleEndian.PutUint32(indexData[i * 8 + 4:], offset)
	}

	copy(b.Bytes()[indexStart:], indexData)
}