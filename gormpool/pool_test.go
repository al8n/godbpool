package gormpool

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMySQLNewPool(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		DBType:          "mysql",
		DBArgs:          []interface{}{"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True"},
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 5 * time.Second, i)
	}
	wg.Wait()
	canc()
}

func TestMySQLClose(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		DBType:          "mysql",
		DBArgs:          []interface{}{"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True"},
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 2 * time.Second, i)
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 8 * time.Second, i + 2)
	}
	time.Sleep(4 * time.Second)
	canc()
	wg.Wait()
}

func mockJob(wg *sync.WaitGroup, p *Pool, duration time.Duration, idx int)  {
	fmt.Printf("Routine %d begin\n", idx)
	format(idx, p.IdleConn(), "Idle")
	format(idx, p.BusyConn(), "Busy")
	conn, err := p.Get()
	if err != nil {
		fmt.Println(err)
	} else {
		format(idx, p.IdleConn(), "Idle")
		format(idx, p.BusyConn(), "Busy")
		time.Sleep(duration)
		p.Put(conn)
		format(idx, p.IdleConn(), "Idle")
		format(idx, p.BusyConn(), "Busy")
	}
	wg.Done()
}

func format(idx int, value uint64, typ string)  {
	fmt.Printf("Routine %d: %s: %d\n", idx, typ, value)
}
