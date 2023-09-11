package framework

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	DefaultTimeout = time.Minute
)

// ReservedPort represents a port that has been reserved but not yet used
type ReservedPort struct {
	port     int
	listener net.Listener
	isClosed bool
}

// Port returns the reserved port number as a string
func (p *ReservedPort) Port() string {
	return strconv.Itoa(p.port)
}

// FindAvailablePort searches for an available port between the specified 'from' and 'to' range
func FindAvailablePort(from, to int) *ReservedPort {
	for port := from; port < to; port++ {
		addr := fmt.Sprintf("localhost:%d", port)
		if l, err := net.Listen("tcp", addr); err == nil {
			return &ReservedPort{port: port, listener: l}
		}
	}

	return nil
}

// Close closes the reserved port and sets the isClosed flag to true
func (p *ReservedPort) Close() error {
	if p.isClosed {
		return nil
	}

	err := p.listener.Close()
	p.isClosed = true

	return err
}
