package old

import (
	"log"
	"net"
	"time"
)

func sendAndReceiveUDP(resolver string, port int, query []byte, timeout time.Duration) ([]byte, error) {
	addr := &net.UDPAddr{
		IP:   net.ParseIP(resolver),
		Port: port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Failed to create UDP stream: %v", err)
		return nil, err
	}

	defer conn.Close()

	err = conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		log.Printf("Failed to set write deadline: %v", err)
		return nil, err
	}

	_, err = conn.Write(query)
	if err != nil {
		log.Printf("Failed to write: %v", err)
		return nil, err
	}

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		log.Printf("Failed to set read deadline: %v", err)
		return nil, err
	}

	buffer := make([]byte, 4096)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("Failed to read from UDP: %v", err)
		return nil, err
	}

	return buffer[:n], nil
}
