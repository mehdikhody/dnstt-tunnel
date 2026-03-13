package server

import (
	"net"
	"sync"
)

type Options struct {
	HOST     string
	PORT     int
	PASSWORD string
	DOMAIN   string
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
		IP:   net.ParseIP(s.options.HOST),
		Port: s.options.PORT,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

}
