package frame

import (
	"bytes"
)

type FrameType uint8

const(
	FRAME_TYPE_PADDING				FrameType = 0x00
	FRAME_TYPE_RST_STREAM			FrameType = 0x01
	FRAME_TYPE_CONNECTION_CLOSE		FrameType = 0x02
	FRAME_TYPE_APPLICATION_CLOSE	FrameType = 0x03
	FRAME_TYPE_MAX_DATA				FrameType = 0x04
	FRAME_TYPE_MAX_STREAM_DATA		FrameType = 0x05
	FRAME_TYPE_MAX_STREAM_ID		FrameType = 0x06
	FRAME_TYPE_PING					FrameType = 0x07
	FRAME_TYPE_BLOCKED				FrameType = 0x08
	FRAME_TYPE_STREAM_BLOCKED		FrameType = 0x09
	FRAME_TYPE_STREAM_ID_BLOCKED	FrameType = 0x0A
	FRAME_TYPE_NEW_CONNECTION_ID	FrameType = 0x0B
	FRAME_TYPE_STOP_SENDING			FrameType = 0x0C
	FRAME_TYPE_ACK					FrameType = 0x0D
	FRAME_TYPE_PATH_CHALLENGE		FrameType = 0x0E
	FRAME_TYPE_PATH_RESPONSE		FrameType = 0x0F
	FRAME_TYPE_STREAM				FrameType = 0x10
	FRAME_TYPE_STREAM_OFF			FrameType = 0x14
	FRAME_TYPE_STREAM_LEN			FrameType = 0x12
	FRAME_TYPE_STREAM_FIN			FrameType = 0x11
)

type Frame struct {
	frameType FrameType
}

type IFrame interface {
	GetType() FrameType
	Serialize(b *bytes.Buffer) error
}