package client

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/armon/go-socks5"
)

type Options struct {
	HOST                 string
	PORT                 int
	DOMAINS              []string
	MIN_UPLOAD_MTU       int64
	MAX_UPLOAD_MTU       int64
	MIN_DOWNLOAD_MTU     int64
	MAX_DOWNLOAD_MTU     int64
	DNS_RESOLVER_TIMEOUT int
}

type Client struct {
	options *Options
	stop    chan bool
	mutex   sync.RWMutex
	conn    *net.UDPConn
}

func New(options *Options) *Client {
	return &Client{
		options: options,
		stop:    make(chan bool),
	}
}

func (c *Client) Start() error {
	server, err := socks5.New(&socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			log.Printf("Dial %s %s", network, addr)

			timeout := time.Duration(c.options.DNS_RESOLVER_TIMEOUT) * time.Millisecond
			resolver, err := NewResolver("127.0.0.1", 53, timeout)
			if err != nil {
				log.Fatal("Failed to create resolver:", err)
				return nil, err
			}

			domain := c.options.DOMAINS[0]
			tunnel, err := NewTunnel(domain, resolver)
			if err != nil {
				log.Fatal("Failed to create tunnel:", err)
				return nil, err
			}

			go func() {
				for {
					data, err := resolver.Receive()
					if err != nil {
						log.Println("Failed to receive:", err)
						return
					}

					log.Printf("Received data: %s", string(data))
					tunnel.read <- data
				}
			}()

			return tunnel, nil
		},
	})

	if err != nil {
		log.Printf("Failed to create socks5 server: %s", err)
		return err
	}

	addr := net.JoinHostPort(c.options.HOST, strconv.Itoa(c.options.PORT))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Failed to listen socks5 server: %s", err)
		return err
	}

	log.Printf("Listening on %s", listener.Addr())
	err = server.Serve(listener)
	if err != nil {
		log.Printf("Failed to serve socks5 server: %s", err)
		return err
	}

	return nil
}
