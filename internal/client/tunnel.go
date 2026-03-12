package client

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Tunnel struct {
	resolver *Resolver
	domain   string
	streamId uint16
	packetId uint16
	read     chan []byte
	close    chan bool
}

func NewTunnel(domain string, resolver *Resolver) (*Tunnel, error) {
	tunnel := &Tunnel{
		resolver: resolver,
		domain:   domain,
		streamId: 0,
		packetId: 0,
		read:     make(chan []byte),
		close:    make(chan bool),
	}

	streamId, err := tunnel.requestStreamId()
	if streamId == 0 || err != nil {
		return nil, errors.New("failed to request stream id")
	}

	tunnel.streamId = streamId
	return tunnel, nil
}

func (t *Tunnel) requestStreamId() (uint16, error) {
	payload := []byte{0x00, 0x00, 0x00, 0x00}
	err := t.resolver.Send(t.domain, payload)
	if err != nil {
		log.Printf("Failed to request stream id: %v", err)
		return 0, err
	}

	resp, err := t.resolver.Receive()
	if err != nil {
		log.Printf("Failed to request stream id: %v", err)
		return 0, err
	}

	if len(resp) < 2 {
		return 0, errors.New("failed to request stream id")
	}

	return binary.BigEndian.Uint16(resp[:2]), nil
}

func (t *Tunnel) Write(b []byte) (int, error) {
	offset := 0

	for offset < len(b) {
		chunkSize := 1024
		if offset+chunkSize > len(b) {
			chunkSize = len(b) - offset
		}

		t.packetId++

		buf := bytes.Buffer{}
		binary.Write(&buf, binary.BigEndian, t.streamId)
		binary.Write(&buf, binary.BigEndian, t.packetId)
		buf.Write(b[offset : offset+chunkSize])

		if err := t.resolver.Send(t.domain, buf.Bytes()); err != nil {
			return offset, err
		}

		offset += chunkSize
	}

	return len(b), nil
}

func (t *Tunnel) Read(b []byte) (int, error) {
	select {
	case data := <-t.read:
		n := copy(b, data)
		return n, nil

	case <-t.close:
		return 0, fmt.Errorf("connection closed")
	}
}

func (t *Tunnel) Close() error {
	close(t.close)
	return t.resolver.Close()
}

func (t *Tunnel) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4zero, Port: 0}
}

func (t *Tunnel) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4zero, Port: 0}
}

func (t *Tunnel) SetDeadline(_ time.Time) error {
	return nil
}

func (t *Tunnel) SetReadDeadline(_ time.Time) error {
	return nil
}

func (t *Tunnel) SetWriteDeadline(_ time.Time) error {
	return nil
}
