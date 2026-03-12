package old

import (
	"bytes"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type Connection struct {
	readBuffer bytes.Buffer
	readMutex  sync.Mutex
	writeMutex sync.Mutex
	closed     bool
}

func (c *Connection) Read(b []byte) (int, error) {
	log.Printf("Reading %x\n", len(b))

	c.readMutex.Lock()
	defer c.readMutex.Unlock()

	for c.readBuffer.Len() == 0 {
		if c.closed {
			return 0, errors.New("connection closed")
		}

		time.Sleep(time.Millisecond * 10)
	}

	return c.readBuffer.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	log.Printf("Writing %x\n", len(b))

	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	return len(b), nil
}

func (c *Connection) Close() error {
	c.readMutex.Lock()
	defer c.readMutex.Unlock()
	c.closed = true
	return nil
}

func (c *Connection) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (c *Connection) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (c *Connection) SetDeadline(t time.Time) error {
	return nil
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}
