package ackhandler

import (
	"../frame"
)

func stripNonRetransmittableFrames(fs []frame.IFrame) []frame.IFrame {
	ret := make([]frame.IFrame, 0, len(fs))
	for _, f := range fs {
		if IsFrameRetransmittable(f) {
			ret = append(ret, f)
		}
	}
	return ret
}

func IsFrameRetransmittable(f frame.IFrame) bool {
	switch f.GetType() {
	case frame.FRAME_TYPE_ACK:
		return false
	default:
		return true
	}
}

func HasFrameRetransmittable(fs []frame.IFrame) bool {
	for _, f := range fs {
		if IsFrameRetransmittable(f) {
			return true
		}
	}
	return false
}