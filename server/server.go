package server

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/TransIRC/cesiumlib"
)

type Options struct {
	Host              string
	Port              int
	Password          string
	Domain            string
	MaxRetransmits    int
	FlowControlWindow int
	KeepaliveInterval int
	AckTimeout        int
	WriteTimeout      int
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
	addr := &net.UDPAddr{
		IP:   net.ParseIP(s.options.Host),
		Port: s.options.Port,
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on %v\n", addr)
	log.Printf("Tunnel Domain: %s\n", s.options.Domain)
	log.Printf("Tunnel Password: %s\n", s.options.Password)

	cesiumlib.Configure(cesiumlib.Config{
		MaxRetransmits:    s.options.MaxRetransmits,
		FlowControlWindow: s.options.FlowControlWindow,
		KeepaliveInterval: time.Duration(s.options.KeepaliveInterval) * time.Millisecond,
		AckTimeout:        time.Duration(s.options.AckTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(s.options.WriteTimeout) * time.Millisecond,
	})

	err = cesiumlib.AcceptServerDnsTunnelConns(
		udpConn,
		s.options.Domain,
		s.options.Password,
		func(conn net.Conn) {
			defer conn.Close()
			buffer := make([]byte, 4096)
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					return
				}

				log.Printf("Received: %s", buffer[:n])
				conn.Write([]byte("pong"))
			}
		},
	)

	if err != nil {
		panic(err)
	}
}
