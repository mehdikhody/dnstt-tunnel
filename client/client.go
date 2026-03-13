package client

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
	Domain            string
	Password          string
	Resolvers         []string
	ChunkSize         int
	MaxRetransmits    int
	FlowControlWindow int
	KeepaliveInterval int
	ReadPollInterval  int
	AckTimeout        int
	WriteTimeout      int
}

type Client struct {
	options *Options
	mutex   sync.RWMutex
	conn    *net.UDPConn
	stop    chan bool
}

func New(options *Options) *Client {
	return &Client{
		options: options,
		stop:    make(chan bool),
	}
}

func (c *Client) Start() {
	resolver := "127.0.0.1:53"
	domain := c.options.Domain
	password := c.options.Password

	cesiumlib.Configure(cesiumlib.Config{
		ClientRawChunkSize: c.options.ChunkSize,
		MaxRetransmits:     c.options.MaxRetransmits,
		FlowControlWindow:  c.options.FlowControlWindow,
		KeepaliveInterval:  time.Duration(c.options.KeepaliveInterval) * time.Millisecond,
		ReadPollInterval:   time.Duration(c.options.ReadPollInterval) * time.Millisecond,
		AckTimeout:         time.Duration(c.options.AckTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.options.WriteTimeout) * time.Millisecond,
	})

	conn, err := cesiumlib.NewDnsTunnelConn(resolver, domain, password)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	conn.Write([]byte("ping"))
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	log.Printf("Got: %s", buffer[:n])
}
