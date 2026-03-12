package client

import (
	"context"
	"dnstt-tunnel/internal/client/old/stream"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/armon/go-socks5"
)

type Options struct {
	HOST             string
	PORT             int
	DOMAINS          []string
	MIN_UPLOAD_MTU   int64
	MAX_UPLOAD_MTU   int64
	MIN_DOWNLOAD_MTU int64
	MAX_DOWNLOAD_MTU int64
}

type Client struct {
	Options *Options
	Signal  chan string
	Mutex   sync.RWMutex
	Conn    *net.UDPConn
	Streams map[string]map[uint8]*stream.Stream
}

func New(options *Options) *Client {
	return &Client{
		Options: options,
		Signal:  make(chan string),
		Streams: make(map[string]map[uint8]*stream.Stream),
	}
}

func (c *Client) Start() {
	addr := net.JoinHostPort(c.Options.HOST, strconv.Itoa(c.Options.PORT))

	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			fmt.Printf("socks5 dial %s %s\n", network, addr)
			conn, err := net.Dial(network, addr)
			if err != nil {
				log.Printf("Failed to dial tunnel for %s: %v", addr, err)
			}

			return conn, err
		},
	}

	server, err := socks5.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	log.Printf("SOCKS5 server listening on %s", addr)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("SOCKS5 server error: %v", err)
	}
}
