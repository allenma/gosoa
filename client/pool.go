package client

import (
	"container/list"
	"errors"
	"io"
	"sync"
	"time"
)

var nowFunc = time.Now // for testing

// ErrPoolExhausted is returned from a pool connection method when the
// maximum number of database connections in the pool has been reached.
var ErrPoolExhausted = errors.New("connection pool exhausted")

var (
	errPoolClosed = errors.New("connection pool closed")
	errConnClosed = errors.New("connection closed")
)

type Conn interface {
	io.Closer
	// when do some operation on this conn, if there is error
	// and you can set the error to the conn, so that the conn
	// will not be used again
	GetErr() error
}

type Pool struct {

	// Dial is an application supplied function for creating and configuring a
	// connection
	Dial func() (Conn, error)

	// TestOnBorrow is an optional application supplied function for checking
	// the health of an idle connection before the connection is used again by
	// the application. Argument t is the time that the connection was returned
	// to the pool. If the function returns an error, then the connection is
	// closed.
	TestOnBorrow func(c Conn, t time.Time) error

	// Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// If Wait is true and the pool is at the MaxIdle limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool

	// mu protects fields defined below.
	mu     sync.Mutex
	cond   *sync.Cond
	closed bool
	active int

	// Stack of idleConn with most recently used at the front.
	idle list.List
}

type idleConn struct {
	c Conn
	t time.Time
}

// NewPool creates a new pool. This function is deprecated. Applications should
// initialize the Pool fields directly as shown in example.
func NewPool(newFn func() (Conn, error), maxIdle int) *Pool {
	return &Pool{Dial: newFn, MaxIdle: maxIdle}
}

// ActiveCount returns the number of active connections in the pool.
func (p *Pool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

// Close releases the resources used by the pool.
func (p *Pool) Close() error {
	p.mu.Lock()
	idle := p.idle
	p.idle.Init()
	p.closed = true
	p.active -= idle.Len()
	if p.cond != nil {
		p.cond.Broadcast()
	}
	p.mu.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(idleConn).c.Close()
	}
	return nil
}

// release decrements the active count and signals waiters. The caller must
// hold p.mu during the call.
func (p *Pool) release() {
	p.active -= 1
	if p.cond != nil {
		p.cond.Signal()
	}
}

// Borrow prunes stale connections and returns a connection from the idle list or
// creates a new connection.
func (p *Pool) Borrow() (Conn, error) {
	p.mu.Lock()

	// Prune stale connections.
	if timeout := p.IdleTimeout; timeout > 0 {
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Back()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			if ic.t.Add(timeout).After(nowFunc()) {
				break
			}
			p.idle.Remove(e)
			p.release()
			p.mu.Unlock()
			ic.c.Close()
			p.mu.Lock()
		}
	}

	for {

		// Get idle connection.
		for i, n := 0, p.idle.Len(); i < n; i++ {
			e := p.idle.Front()
			if e == nil {
				break
			}
			ic := e.Value.(idleConn)
			p.idle.Remove(e)
			test := p.TestOnBorrow
			p.mu.Unlock()
			if test == nil || test(ic.c, ic.t) == nil {
				return ic.c, nil
			}
			ic.c.Close()
			p.mu.Lock()
			p.release()
		}

		// Check for pool closed before dialing a new connection.
		if p.closed {
			p.mu.Unlock()
			return nil, errors.New("get on closed pool")
		}

		// Dial new connection if under limit.
		if p.MaxActive == 0 || p.active < p.MaxActive {
			dial := p.Dial
			p.active += 1
			p.mu.Unlock()
			c, err := dial()
			if err != nil {
				p.mu.Lock()
				p.release()
				p.mu.Unlock()
				c = nil
			}
			return c, err
		}

		if !p.Wait {
			p.mu.Unlock()
			return nil, ErrPoolExhausted
		}

		if p.cond == nil {
			p.cond = sync.NewCond(&p.mu)
		}
		p.cond.Wait()
	}
}

// Return adds conn back to the pool, use forceClose to close the connection forcely
func (p *Pool) Return(c Conn, forceClose bool) error {
	if c == nil {
		return nil
	}
	err := c.GetErr()
	p.mu.Lock()
	if !p.closed && err == nil && !forceClose {
		p.idle.PushFront(idleConn{t: nowFunc(), c: c})
		if p.idle.Len() > p.MaxIdle {
			c = p.idle.Remove(p.idle.Back()).(idleConn).c
		} else {
			c = nil
		}
	}

	if c == nil {
		if p.cond != nil {
			p.cond.Signal()
		}
		p.mu.Unlock()
		return nil
	}

	p.release()
	p.mu.Unlock()
	return c.Close()
}
