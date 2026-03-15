package server

import (
	"io"
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
		KeepaliveInterval: 0,
		AckTimeout:        time.Duration(s.options.AckTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(s.options.WriteTimeout) * time.Millisecond,
	})

	err = cesiumlib.AcceptServerDnsTunnelConns(
		udpConn,
		s.options.Domain,
		s.options.Password,
		s.handle,
	)

	if err != nil {
		panic(err)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	lenBuffer := make([]byte, 1)
	_, err := io.ReadFull(conn, lenBuffer)
	if err != nil {
		log.Printf("Error reading length: %v", err)
		return
	}

	destLen := int(lenBuffer[0])
	if destLen <= 0 || destLen > 255 {
		log.Printf("Invalid destination length: %d", destLen)
		return
	}

	destBuffer := make([]byte, destLen)
	_, err = io.ReadFull(conn, destBuffer)
	if err != nil {
		log.Printf("Error reading destination: %v", err)
		return
	}

	dest := string(destBuffer)
	log.Printf("Received destination: %s\n", dest)

	remote, err := net.Dial("tcp", dest)
	if err != nil {
		log.Printf("Error connecting to server: %v", err)
		return
	}

	defer remote.Close()

	go io.Copy(conn, remote)
	io.Copy(remote, conn)
}
