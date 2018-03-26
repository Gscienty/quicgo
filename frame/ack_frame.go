package frame

import(
	"time"
	"errors"
	"../utils"
	"bytes"
)

type AckBlock struct {
	First	uint64
	Last	uint64
}

type AckFrame struct {
	Frame
	LargestAcknowledged	utils.VarLenIntegerStruct
	LowAcknowledged		uint64
	PacketReceivedTime	time.Time
	ackDelay			utils.VarLenIntegerStruct
	ackBlockCount		utils.VarLenIntegerStruct

	Blocks				[]AckBlock
}

func AckFrameParse(b *bytes.Reader) (*AckFrame, error) {
	frameType, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	if FrameType(frameType) != FRAME_TYPE_ACK {
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

	return &AckFrame { Frame { FrameType(frameType) }, *largestAcknowledged, smallest, time.Time { }, *delay, *blockCount, blocks }, nil
}

func (this *AckFrame) GetType() FrameType {
	return FRAME_TYPE_ACK
}


func (this *AckFrame) Serialize(b *bytes.Buffer) error {
	err := b.WriteByte(uint8(this.frameType))
	if err != nil {
		return err
	}
	
	_, err = this.LargestAcknowledged.Serialize(b)
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

	utils.VarLenIntegerStructNew(this.LargestAcknowledged.GetVal() - this.Blocks[0].First).Serialize(b)
	var lowest uint64 = this.Blocks[0].First
	for i, block := range this.Blocks {
		if i == 0 {
			continue
		}

		utils.VarLenIntegerStructNew(lowest - block.Last - 2).Serialize(b)
		utils.VarLenIntegerStructNew(block.Last - block.First).Serialize(b)
		lowest = block.First
	}
	return nil
}

func (this *AckFrame) HasMissingRanges() bool {
	return len(this.Blocks) > 0
}

func (this *AckFrame) AcksPacket(p uint64) bool {
	if p < this.LowAcknowledged || p > this.LargestAcknowledged.GetVal() {
		return false
	}

	if this.HasMissingRanges() {
		for _, block := range this.Blocks {
			if block.First <= p && p <= block.Last {
				return true
			}
		}
		return false
	}

	return (this.LowAcknowledged <= p && p <= this.LargestAcknowledged.GetVal())
}