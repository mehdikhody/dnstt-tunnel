package client

import (
	"net"
	"time"
)

type Forwarder struct {
	conn       net.Conn
	localAddr  net.Addr
	remoteAddr net.Addr
	headerSent bool
}

func NewForwarder(tun net.Conn) {
	tun.LocalAddr()
}

func (t *Forwarder) Write(p []byte) (int, error) {
	if !t.headerSent {
		destBytes := []byte(t.remoteAddr.String())
		header := append([]byte{byte(len(destBytes))}, destBytes...)
		_, err := t.conn.Write(header)
		if err != nil {
			return 0, err
		}

		t.headerSent = true
	}

	return t.conn.Write(p)
}

func (t *Forwarder) Read(p []byte) (int, error)          { return t.conn.Read(p) }
func (t *Forwarder) LocalAddr() net.Addr                 { return t.localAddr }
func (t *Forwarder) Close() error                        { return t.conn.Close() }
func (t *Forwarder) RemoteAddr() net.Addr                { return t.remoteAddr }
func (t *Forwarder) SetDeadline(ti time.Time) error      { return t.conn.SetDeadline(ti) }
func (t *Forwarder) SetReadDeadline(ti time.Time) error  { return t.conn.SetReadDeadline(ti) }
func (t *Forwarder) SetWriteDeadline(ti time.Time) error { return t.conn.SetWriteDeadline(ti) }
