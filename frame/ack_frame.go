package frame

import(
	"errors"
	"../utils"
	"bytes"
)

type AckBlock struct {
	first	uint64
	last	uint64
}

type AckFrame struct {
	Frame
	largestAcknowledged	utils.VarLenIntegerStruct
	ackDelay			utils.VarLenIntegerStruct
	ackBlockCount		utils.VarLenIntegerStruct

	blocks				[]AckBlock
}

func AckFrameParse(b *bytes.Reader) (*AckFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if frameType != FRAME_TYPE_ACK {
		return nil, errors.New("AckFrameParse error: frametype not equal 0x0E")
	}

	largestAcknowledged, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}
	delay, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}
	blockCount, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}
	blocksCount := blockCount.GetVal()

	firstBlock, err := utils.VarLenIntegerStructParse(b)
	if err != nil {
		return nil, err
	}
	if firstBlock.GetVal() > largestAcknowledged.GetVal() {
		return nil, errors.New("AckFrameParse error: invalid first ACK range")
	}
	smallest := largestAcknowledged.GetVal() - firstBlock.GetVal()

	var blocks []AckBlock

	if blocksCount > 0 {
		blocks = append(blocks, AckBlock { smallest, largestAcknowledged.GetVal() })
	}

	for i := uint64(0); i < blocksCount; i++ {
		gap, err := utils.VarLenIntegerStructParse(b)
		if err != nil {
			return nil, err
		}
		if smallest < gap.GetVal() + 2 {
			return nil, errors.New("AckFrameParse error: invalid ack block")
		}
		largest := smallest - gap.GetVal() - 2

		block, err := utils.VarLenIntegerStructParse(b)
		if err != nil {
			return nil, err
		}
		if block.GetVal() > largest {
			return nil, errors.New("AckFrameParse error: invalid ack block")
		}
		smallest = largest - block.GetVal()

		blocks = append(blocks, AckBlock { smallest, largest })
	}

	return &AckFrame { Frame { frameType }, *largestAcknowledged, *delay, *blockCount, blocks }, nil
}

func (this *AckFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(this.frameType)
	if err != nil {
		return err
	}
	
	_, err = this.largestAcknowledged.Serialize(b)
	if err != nil {
		return err
	}

	_, err = this.ackDelay.Serialize(b)
	if err != nil {
		return err
	}

	_, err = this.ackBlockCount.Serialize(b)
	if err != nil {
		return err
	}

	utils.VarLenIntegerStructNew(this.largestAcknowledged.GetVal() - this.blocks[0].first).Serialize(b)
	var lowest uint64 = this.blocks[0].first
	for i, block := range this.blocks {
		if i == 0 {
			continue
		}

		utils.VarLenIntegerStructNew(lowest - block.last - 2).Serialize(b)
		utils.VarLenIntegerStructNew(block.last - block.first).Serialize(b)
		lowest = block.first
	}
	return nil
}