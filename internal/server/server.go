package server

import (
	"log"
	"net"
	"sync"
)

type Options struct {
	HOST    string
	PORT    int
	DOMAINS []string
}

type Server struct {
	options *Options
	mutex   sync.RWMutex
	conn    *net.UDPConn
	close   chan bool
}

func New(options *Options) *Server {
	return &Server{
		options: options,
		close:   make(chan bool),
	}
}

func (s *Server) Start() {
	addr := net.UDPAddr{
		IP:   net.ParseIP(s.options.HOST),
		Port: s.options.PORT,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s.conn = conn
	defer conn.Close()

	buf := make([]byte, 65535)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		go s.handlePacket(data, clientAddr)
	}
}

func (s *Server) Close() error {
	close(s.close)
	return s.conn.Close()
}
