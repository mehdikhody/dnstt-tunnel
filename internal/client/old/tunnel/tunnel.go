package tunnel

import (
	"dnstt-tunnel/internal/client/old/resolver"
)

type Tunnel struct {
	resolver *resolver.Resolver
	streamID uint16
	packetId uint16
	read     chan []byte
	close    chan bool
}
