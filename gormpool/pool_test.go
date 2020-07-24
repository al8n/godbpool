package gormpool

import (
	"context"
	"fmt"
	godbpool "github.com/ALiuGuanyan/go-db-pool"
	"sync"
	"testing"
	"time"
)

func TestMySQLNewPool(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
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
		go mockJob(&wg, p, 5*time.Second)
	}
	wg.Wait()
	canc()
}

func TestErrorPool(t *testing.T) {
	var (
		err error
	)

	ctx := context.Background()
	opts1 := Options{
		Type:            godbpool.MySQL,
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}
	_, err = NewPool(ctx, opts1)
	if err == nil {
		t.Error()
	}

	opts2 := Options{
		Type:            godbpool.MySQL,
		KeepConn:        7,
		Capacity:        5,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	_, err = NewPool(ctx, opts2)
	if err == nil {
		t.Error()
	}

	opts3 := Options{
		Type:            godbpool.MySQL,
		KeepConn:        0,
		Capacity:        0,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	_, err = NewPool(ctx, opts3)
	if err == nil {
		t.Error()
	}

	opts4 := Options{
		Type:            godbpool.MySQL,
		KeepConn:        0,
		Capacity:        0,
		Args:            "root:1234568910@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	_, err = NewPool(ctx, opts4)
	if err == nil {
		t.Error()
	}
}



func TestGetFromClosedPool(t *testing.T) {
	ctx := context.Background()
	opts := Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 1*time.Second)
	}
	p.Close()
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 6*time.Second)
	}
	time.Sleep(4 * time.Second)
	wg.Wait()
	if !p.Status().Closed {
		t.Error()
	}
}

func TestMySQLClose(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 1*time.Second)
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 6*time.Second)
	}
	time.Sleep(4 * time.Second)
	canc()
	wg.Wait()
}

func TestWaitingGet(t *testing.T) {
	ctx, canc := context.WithCancel(context.Background())
	opts := Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		KeepConn:        1,
		Capacity:        2,
		MaxWaitDuration: 2000 * time.Millisecond,
	}
	p, err := NewPool(ctx, opts)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 2*time.Second)
	}
	time.Sleep(time.Second)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 3*time.Second)
	}
	wg.Wait()
	canc()
}

func mockJob(wg *sync.WaitGroup, p *Pool, duration time.Duration) {
	conn, err := p.Get()
	if err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(duration)
		p.Put(conn)
	}
	wg.Done()
}
