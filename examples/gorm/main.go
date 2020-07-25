package main

import (
	"context"
	"fmt"
	"github.com/ALiuGuanyan/godbpool"
	"github.com/ALiuGuanyan/godbpool/gormpool"
	"log"
	"sync"
	"time"
)

func main() {
	// config options
	opts := gormpool.Options{
		Type:            godbpool.MySQL,
		Args:            "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True",
		KeepConn:        2,
		Capacity:        5,
		MaxWaitDuration: 2000 * time.Millisecond,
	}

	// create a pool
	ctx, canc := context.WithCancel(context.Background())
	p, err := gormpool.NewPool(ctx, opts)
	if err != nil {
		log.Println(err)
		return
	}

	// mock to do some CURD jobs
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go mockJob(&wg, p, 5*time.Second, i)
	}
	wg.Wait()
	p.Close()
	canc()
}

func mockJob(wg *sync.WaitGroup, p *gormpool.Pool, duration time.Duration, idx int) {
	fmt.Printf("Routine %d begin\n", idx)
	conn, err := p.Get()
	if err != nil {
		fmt.Println(err)
	} else {
		time.Sleep(duration)
		// Put conn back to pool
		p.Put(conn)
	}
	wg.Done()
}

func format(idx int, value uint64, typ string) {
	fmt.Printf("Routine %d: %s: %d\n", idx, typ, value)
}
