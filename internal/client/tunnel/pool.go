package tunnel

import (
	"errors"
	"math/rand"
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

func NewTunnelPool(options *PoolOptions) *Pool {
	pool := &Pool{
		options:   options,
		actives:   make([]*Tunnel, 10),
		inactives: make([]*Tunnel, 10),
	}

	return pool
}

func (p *Pool) Length() int {
	return len(p.actives)
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
