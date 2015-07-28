package client

import (
	//	"fmt"
	"errors"
	"sync"
	"testing"
	"time"
)

type poolTestConn struct {
	d                *poolDialer
	err              error
	actionSuccessCnt int
	actionFailedCnt  int
}

func (c *poolTestConn) Close() error  { c.d.open -= 1; return nil }
func (c *poolTestConn) GetErr() error { return c.err }
func (c *poolTestConn) action(err error) error {
	if c.err != nil {
		return c.err
	}

	if err != nil {
		c.actionFailedCnt++
		c.err = err
	} else {
		c.actionSuccessCnt++
	}
	return nil
}

type poolDialer struct {
	t       *testing.T
	dialed  int
	open    int
	dialErr error
}

func (d *poolDialer) dial() (Conn, error) {
	d.dialed += 1
	if d.dialErr != nil {
		return nil, d.dialErr
	}
	d.open += 1
	return &poolTestConn{d: d}, nil
}

func (d *poolDialer) check(message string, p *Pool, dialed, open int) {
	if d.dialed != dialed {
		d.t.Errorf("%s: dialed=%d, want %d", message, d.dialed, dialed)
	}
	if d.open != open {
		d.t.Errorf("%s: open=%d, want %d", message, d.open, open)
	}
	if active := p.ActiveCount(); active != open {
		d.t.Errorf("%s: active=%d, want %d", message, active, open)
	}
}

func TestBorrow(t *testing.T) {
	d := poolDialer{t: t}
	p := &Pool{
		MaxIdle:   10,
		MaxActive: 10,
		Dial:      d.dial,
		Wait:      true,
	}
	wg := new(sync.WaitGroup)
	wg.Add(50)
	for i := 0; i < 50; i++ {
		go run(p, wg, t)
	}
	wg.Wait()
	d.check("before close", p, 10, 10)
	p.Close()
	d.check("after close", p, 10, 0)
}

func TestConnErr(t *testing.T) {
	d := poolDialer{t: t}
	p := &Pool{
		MaxIdle:   10,
		MaxActive: 10,
		Dial:      d.dial,
		Wait:      true,
	}
	conns := make([]*poolTestConn, 0, 10)
	for i := 0; i < 10; i++ {
		c1, _ := p.Borrow()
		pc1 := c1.(*poolTestConn)
		if i%3 == 0 {
			pc1.action(errors.New("action error"))
		} else {
			pc1.action(nil)
		}
		conns = append(conns, pc1)
	}
	for _, conn := range conns {
		p.Return(conn, false)
	}
	d.check("when conn error", p, 10, 6)

	conns = make([]*poolTestConn, 0, 10)
	for i := 0; i < 10; i++ {
		c1, _ := p.Borrow()
		pc1 := c1.(*poolTestConn)
		pc1.action(nil)
		conns = append(conns, pc1)
	}
	for _, conn := range conns {
		p.Return(conn, false)
	}
	d.check("try again after conn error", p, 14, 10)
}

func TestIdle(t *testing.T) {
	d := poolDialer{t: t}
	p := &Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Dial:      d.dial,
		Wait:      true,
	}
	conns := make([]*poolTestConn, 0, 10)
	for i := 0; i < 10; i++ {
		c1, _ := p.Borrow()
		pc1 := c1.(*poolTestConn)
		pc1.action(nil)
		conns = append(conns, pc1)
	}
	for _, conn := range conns {
		p.Return(conn, false)
	}
	d.check("idle test", p, 10, 5)
}

func TestDialErr(t *testing.T) {
	d := poolDialer{t: t, dialErr: errors.New("connect error")}
	p := &Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Dial:      d.dial,
		Wait:      true,
	}
	for i := 0; i < 10; i++ {
		c1, _ := p.Borrow()
		p.Return(c1, false)
	}
	d.check("idle test", p, 10, 0)
}

func run(p *Pool, wg *(sync.WaitGroup), t *testing.T) {
	c1, err := p.Borrow()
	if err != nil {
		t.Error("error:", err.Error())
		return
	}
	pc1 := c1.(*poolTestConn)
	pc1.action(nil)
	time.Sleep(100 * time.Millisecond)
	p.Return(pc1, false)
	wg.Done()
}
