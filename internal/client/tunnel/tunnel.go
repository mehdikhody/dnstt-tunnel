package tunnel

import (
	"net"
	"strconv"
	"time"

	"github.com/TransIRC/cesiumlib"
)

type TunnelOptions struct {
	IP       string
	Port     int
	Domain   string
	Password string
}

type Tunnel struct {
	IP        string
	Port      int
	Domain    string
	Password  string
	Conn      net.Conn
	IsActive  bool
	Success   int64
	Fail      int64
	Ping      int
	UpdatedAt time.Time
}

func NewTunnel(options *TunnelOptions) *Tunnel {
	t := &Tunnel{
		IP:        options.IP,
		Port:      options.Port,
		Domain:    options.Domain,
		Password:  options.Password,
		IsActive:  false,
		Success:   0,
		Fail:      0,
		Ping:      0,
		UpdatedAt: time.Now(),
	}

	addr := net.JoinHostPort(t.IP, strconv.Itoa(t.Port))
	conn, err := cesiumlib.NewDnsTunnelConn(addr, c.options.Domain, c.options.Password

	return tunnel
}
