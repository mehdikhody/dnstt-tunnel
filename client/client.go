package client

import (
	"net"
	"sync"
)

type Options struct {
	HOST                 string
	PORT                 int
	DOMAINS              []string
	MIN_UPLOAD_MTU       int
	MAX_UPLOAD_MTU       int
	MIN_DOWNLOAD_MTU     int
	MAX_DOWNLOAD_MTU     int
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

func (c *Client) Start() {

}
