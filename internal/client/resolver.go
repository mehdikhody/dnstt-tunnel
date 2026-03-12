package client

import (
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Resolver struct {
	ip      string
	port    int
	domain  string
	timeout time.Duration
	remote  *net.UDPAddr
	conn    *net.UDPConn
}

func NewResolver(ip string, port int, timeout time.Duration) (*Resolver, error) {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	addr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		log.Printf("Failed to resolve UDP address: %s", err)
		return nil, errors.New("failed to resolve UDP address")
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Failed to dial UDP: %s", err)
		return nil, errors.New("failed to dial UDP")
	}

	resolver := &Resolver{
		ip:      ip,
		port:    port,
		timeout: timeout,
		remote:  addr,
		conn:    conn,
	}

	return resolver, nil
}

func (r *Resolver) Send(domain string, payload []byte) error {
	err := r.conn.SetWriteDeadline(time.Now().Add(r.timeout))
	if err != nil {
		log.Printf("Failed to set write deadline: %s", err)
		return err
	}

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	data := fmt.Sprintf("%s.%s", encoder.EncodeToString(payload), domain)
	_, err = r.conn.Write([]byte(data))
	return err
}

func (r *Resolver) Receive() ([]byte, error) {
	err := r.conn.SetReadDeadline(time.Now().Add(r.timeout))
	if err != nil {
		log.Printf("Failed to set read deadline: %s", err)
		return nil, err
	}

	buffer := make([]byte, 4096)
	n, _, err := r.conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("Failed to read from UDP: %s", err)
		return nil, err
	}

	return buffer[:n], nil
}

func (r *Resolver) Close() error {
	return r.conn.Close()
}
