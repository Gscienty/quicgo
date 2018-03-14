package handshake

type KeyUpdateRequest uint8

const (
	KEY_UPDATE_NOT_REQUESTED	= KeyUpdateRequest(0)
	KEY_UPDATE_REQUESTED		= KeyUpdateRequest(1)
)

type KeyUpdate struct {
	RequestUpdate KeyUpdateRequest
}