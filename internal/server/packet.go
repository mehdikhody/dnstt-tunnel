package server

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func (s *Server) handlePacket(data []byte, addr *net.UDPAddr) {
	fmt.Printf("packet received from %s\n", addr.String())
	fmt.Printf("data: %s\n", string(data))

	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	payload, err := decoder.DecodeString(string(data))
	if err != nil {
		fmt.Printf("error decoding payload: %v\n", err)
		return
	}

	if len(payload) < 4 {
		fmt.Printf("error decoding payload: %v\n", err)
		return
	}

	streamID := binary.BigEndian.Uint16(payload[0:2])
	packetID := binary.BigEndian.Uint16(payload[2:4])
	body := payload[4:]

	log.Printf("packet received from %s\n", addr.String())
	log.Printf("Stream ID: %d, Packet ID: %d\n", streamID, packetID)
	log.Printf("data: %s\n", string(body))

	if streamID == 0 && packetID == 0 && len(body) == 0 {

	}

	reply := []byte("ACK")
	_, err = s.conn.WriteToUDP(reply, addr)
	if err != nil {
		fmt.Printf("failed to send reply: %v\n", err)
	}
}
