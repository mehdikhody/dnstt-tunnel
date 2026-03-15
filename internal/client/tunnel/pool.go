package tunnel

import (
	"errors"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	StrategyRandom int = iota
	StrategyRoundRobin
	StrategyLeastLoss
	StrategyLowestLatency
)

type PoolOptions struct {
	Domain        string
	Password      string
	Resolvers     []string
	Strategy      int
	CheckInterval time.Duration
}

type Pool struct {
	options   *PoolOptions
	actives   []*Tunnel
	inactives []*Tunnel
	mutex     sync.RWMutex
}

func NewPool(options *PoolOptions) *Pool {
	p := &Pool{
		options:   options,
		actives:   make([]*Tunnel, 10),
		inactives: make([]*Tunnel, 10),
	}

	for i := 0; i < len(options.Resolvers); i++ {
		addr := options.Resolvers[i]
		if !strings.Contains(addr, ":") {
			addr = net.JoinHostPort(addr, "53") // default port
		}

		udpAddr, err := net.ResolveUDPAddr("tcp", addr)
		if err != nil {
			log.Fatal(err)
			continue
		}

		t, err := NewTunnel(TunnelOptions{
			IP:       udpAddr.IP.String(),
			Port:     udpAddr.Port,
			Domain:   options.Domain,
			Password: options.Password,
		})

		if err != nil {
			log.Fatal(err)
			continue
		}

		p.inactives = append(p.inactives, t)
	}

	p.scheduleChecker()
	return p
}

func (p *Pool) Length() int {
	return len(p.actives)
}

func (p *Pool) Ready() bool {
	return len(p.actives) > 0
}

func (p *Pool) Get() (*Tunnel, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if len(p.actives) == 0 {
		return nil, errors.New("no active tunnels in pool")
	}

	i := rand.Intn(len(p.actives))
	return p.actives[i], nil
}

func (p *Pool) scheduleChecker() {
	ticker := time.NewTimer(p.options.CheckInterval)
	go func() {
		for range ticker.C {
			p.mutex.Lock()
		}
	}()
}
