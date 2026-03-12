package old

import (
	"bytes"
	"context"
	"dnstt-tunnel/internal/client/old/resolver"
	"dnstt-tunnel/internal/client/old/stream"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

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

func (c *Client) Start() error {
	stream, err := resolver.New("127.0.0.1", 53, 3*time.Second)
	if err != nil {
		panic(err)
	}

	defer stream.Close()

	buf := bytes.Buffer{}

	binary.Write(&buf, binary.BigEndian, uint16(65535))
	binary.Write(&buf, binary.BigEndian, uint16(65535))
	buf.Write([]byte("Hello World"))
	payload := buf.Bytes()

	log.Printf("Buffer %v", buf)
	log.Printf("Payload %v", payload)

	if err := stream.Send(payload); err != nil {
		fmt.Println("send error:", err)
		return err
	}

	resp, err := stream.Receive()
	if err != nil {
		fmt.Println("receive error:", err)
		return err
	}

	fmt.Printf("received %d bytes\n", len(resp))

	server, err := socks5.New(&socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			log.Printf("Dial %s %s", strings.ToUpper(network), addr)
			conn := &Connection{}
			return conn, nil
		},
	})

	if err != nil {
		log.Printf("Failed to create socks5 server: %s", err)
		return err
	}

	addr := net.JoinHostPort(c.Options.HOST, strconv.Itoa(c.Options.PORT))
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
