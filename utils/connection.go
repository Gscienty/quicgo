package utils

import(
	"sync"
	"net"
)

type IConnection interface {
	Send([]byte) error
	Recv([]byte) (int, net.Addr, error)
	Close() error
	GetLocalAddress() net.Addr
	GetRemoteAddress() net.Addr
	SetRemoteAddress(net.Addr)
}

type ConnectionBase struct {
	remoteAddressLock	sync.RWMutex

	conn				net.PacketConn
	remoteAddress		net.Addr
}

var _ IConnection = &ConnectionBase { }

func (this *ConnectionBase) Send(b []byte) error {
	_, err := this.conn.WriteTo(b, this.remoteAddress)
	return err
}

func (this *ConnectionBase) Recv(b []byte) (int, net.Addr, error) {
	return this.conn.ReadFrom(b)
}

func (this *ConnectionBase) Close() error {
	return this.conn.Close()
}

func (this *ConnectionBase) GetLocalAddress() net.Addr {
	return this.conn.LocalAddr()
}

func (this *ConnectionBase) GetRemoteAddress() net.Addr {
	this.remoteAddressLock.RLock()
	addr := this.remoteAddress
	this.remoteAddressLock.RUnlock()
	return addr
}

func (this *ConnectionBase) SetRemoteAddress(addr net.Addr) {
	this.remoteAddressLock.Lock()
	this.remoteAddress = addr
	this.remoteAddressLock.Unlock()
}