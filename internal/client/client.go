package client

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/TransIRC/cesiumlib"
	"github.com/armon/go-socks5"
)

type Options struct {
	SocksHost         string
	SocksPort         int
	SocksUsername     string
	SocksPassword     string
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
}

func New(options *Options) *Client {
	return &Client{
		options: options,
	}
}

func (c *Client) Start() {

	addr := &net.UDPAddr{
		IP:   net.ParseIP(c.options.SocksHost),
		Port: c.options.SocksPort,
	}

	var auth []socks5.Authenticator
	if c.options.SocksUsername != "" && c.options.SocksPassword != "" {
		auth = []socks5.Authenticator{
			socks5.UserPassAuthenticator{
				Credentials: socks5.StaticCredentials{
					c.options.SocksUsername: c.options.SocksPassword,
				},
			},
		}
	}

	server, err := socks5.New(&socks5.Config{
		AuthMethods: auth,
		Dial:        c.handle,
	})

	if err != nil {
		panic(err)
	}

	log.Printf("Starting server on %v", addr)
	err = server.ListenAndServe("tcp", addr.String())

	if err != nil {
		panic(err)
	}
}

func (c *Client) handle(ctx context.Context, network, dest string) (net.Conn, error) {
	log.Printf("SOCKS %s %s", network, dest)

	cesiumlib.Configure(cesiumlib.Config{
		ClientRawChunkSize: c.options.ChunkSize,
		MaxRetransmits:     c.options.MaxRetransmits,
		FlowControlWindow:  c.options.FlowControlWindow,
		KeepaliveInterval:  0,
		ReadPollInterval:   time.Duration(c.options.ReadPollInterval) * time.Millisecond,
		AckTimeout:         time.Duration(c.options.AckTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(c.options.WriteTimeout) * time.Millisecond,
	})

	tunConn, err := cesiumlib.NewDnsTunnelConn("127.0.0.1:53", c.options.Domain, c.options.Password)
	if err != nil {
		log.Printf("DNS tunnel failed: %s", err)
		return nil, err
	}

	destHost, portStr, _ := net.SplitHostPort(dest)
	destPort, _ := strconv.Atoi(portStr)

	conn := &Forwarder{
		conn: tunConn,
		localAddr: &net.TCPAddr{
			IP:   net.ParseIP(c.options.SocksHost),
			Port: c.options.SocksPort,
		},
		remoteAddr: &net.TCPAddr{
			IP:   net.ParseIP(destHost),
			Port: destPort,
		},
	}

	return conn, nil
}
