package handshake

import (
	"io"
	"sync"
	"../protocol"
)

type cryptoSetupClient struct {
	mutex				sync.RWMutex

	hostName			string
	connectionID		protocol.ConnectionID
	version				protocol.Version
	initialVersion		protocol.Version
	negotiatedVersions	protocol.Version

	cryptoStream		io.ReadWriter
}