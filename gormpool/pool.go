package gormpool

import (
	"context"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)


var (
	ErrGetFromClosedPool = errors.New("pool: get from closed pool")
	ErrExceedingMaxWaitingDuration = errors.New("pool: exceeding the maximum waiting duration")
)

type Options struct {
	DBType string

	DBArgs []interface{}

	// how many idle conn to keep when there are no work to do
	// this field should smaller than Capacity
	KeepConn uint64

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	Capacity uint64

	MaxWaitDuration time.Duration
}

type Pool struct {
	DBType string

	DBArgs []interface{}

	// how many idle conn to keep when there are no work to do
	// this field should smaller than Capacity
	KeepConn uint64

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	Capacity uint64

	MaxWaitDuration time.Duration

	DBErrChan chan error // DB errors will be sent in this channel

	mu     sync.Mutex // mu protects the following fields

	idleConn *conns

	busyConn *conns

	closed bool // set to true when the pool is closed.
	ch     chan struct{} // limits open connections when p.Wait is true
	currentWaitCount    uint64         // current number of connections waited for.
	totalWaitCount    uint64         // total number of connections waited for.
	waitDuration time.Duration // total time waited for new connections.
	droppedGetCount uint64

	ctx    context.Context
}

// called when do not know DBType and DBArgs are valid
func (p *Pool) initConn() error   {
	db, err := gorm.Open(p.DBType, p.DBArgs...)
	if err != nil {
		return err
	}
	key := uuid.New().String()
	conn := &Conn{
		DB:      db,
		Key:     key,
		created: time.Now(),
		updated: time.Now(),
	}
	p.ch <- struct{}{}
	p.idleConn.put(conn)
	return nil
}

// called when know DBType and DBArgs are valid
func (p *Pool) newConn() {
	db, err := gorm.Open(p.DBType, p.DBArgs...)
	if err != nil {
		p.DBErrChan <- err
		return
	}
	key := uuid.New().String()
	conn := &Conn{
		DB:      db,
		Key:     key,
		created: time.Now(),
		updated: time.Now(),
	}
	p.ch <- struct{}{}
	p.idleConn.put(conn)
}

func NewPool(ctx context.Context, opts Options) (p *Pool, err error) {
	p = &Pool{
		DBType:           opts.DBType,
		DBArgs:           opts.DBArgs,
		KeepConn:         opts.KeepConn,
		Capacity:         opts.Capacity,
		MaxWaitDuration:  opts.MaxWaitDuration,
		mu:               sync.Mutex{},
		idleConn: 		  newConns(),
		busyConn: 		  newConns(),
		closed:           false,
		ch:               make(chan struct{}, opts.Capacity),
		currentWaitCount: 0,
		totalWaitCount:   0,
		waitDuration:     0,
		droppedGetCount:  0,
		ctx:              ctx,
	}

	err = p.initConn()
	if err != nil {
		return nil, err
	}

	for i := uint64(0); i < p.KeepConn - 1; i++ {
		p.newConn()
	}

	return p, nil
}

// Get the number of current total connections of the pool
func (p *Pool) Size() uint64  {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.busyConn.size + p.idleConn.size
}

// Get the number of current idle connections of the pool
func (p *Pool) IdleConn() uint64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.idleConn.size
}

// Get the number of current busy connections of the pool
func (p *Pool) BusyConn() uint64  {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.busyConn.size
}

// Get a SQL connection from the pool
func (p *Pool) Get() (conn *Conn, err error) {
	select {
	case <-p.ctx.Done():
		p.mu.Lock()
		p.droppedGetCount++
		p.mu.Unlock()
		return nil, ErrGetFromClosedPool
	default:
		p.mu.Lock()
		if p.closed {
			p.droppedGetCount++
			p.mu.Unlock()
			return nil, ErrGetFromClosedPool
		}

		if p.idleConn.size > 0 {
			conn = p.get()
			p.mu.Unlock()
			return conn, nil
		}

		if p.busyConn.size < p.Capacity {
			p.newConn()
			conn = p.get()
			p.mu.Unlock()
			return conn, nil
		}

		timer := time.NewTimer(p.MaxWaitDuration)
		start := time.Now()
		p.currentWaitCount++
		p.totalWaitCount++
		select {
		case <-p.ctx.Done():
			return nil, ErrGetFromClosedPool
		case <-p.ch:
			conn = p.get()
			p.waitDuration += time.Since(start)
			p.currentWaitCount--
			p.mu.Unlock()
			return conn, nil
		case <-timer.C:
			p.waitDuration += time.Since(start)
			p.droppedGetCount++
			p.currentWaitCount--
			p.mu.Unlock()
			return nil, ErrExceedingMaxWaitingDuration
		}
	}
}

func (p *Pool) get() (conn *Conn)  {
	conn = p.idleConn.get()
	p.busyConn.put(conn)
	<-p.ch
	return conn
}

// Put back a connection in the pool
func (p *Pool) Put(conn *Conn)  {
	p.mu.Lock()
	p.busyConn.deleteByKey(conn.Key)
	if p.idleConn.size < p.KeepConn && !p.closed {
		p.idleConn.put(conn)
		p.ch <- struct{}{}
	} else {
		conn.DB.Close()
	}
	p.mu.Unlock()
}


func (p *Pool) Close()  {
	p.mu.Lock()
	p.closed = true
	p.idleConn.close()
	p.mu.Unlock()
}

type conns struct {
	keys []string
	conns map[string]*Conn
	size uint64
}

type Conn struct {
	DB *gorm.DB
	Key string
	created, updated time.Time
}

func newConns() *conns  {
	return &conns{
		keys:  []string{},
		conns: map[string]*Conn{},
		size:  0,
	}
}

func (l *conns) get() (conn *Conn) {
	key := l.keys[0]
	l.keys = l.keys[1:]
	conn = l.conns[key]
	delete(l.conns, key)
	l.size--
	return conn
}

func (l *conns) put(conn *Conn) {
	l.keys = append(l.keys, conn.Key)
	l.conns[conn.Key] = conn
	l.size++
}

func (l *conns) deleteByKey(key string)  {
	keys := make([]string, l.size - 1)
	for _, val := range l.keys {
		if val != key {
			keys = append(keys, val)
		}
	}
	l.keys = keys
	l.size--
	delete(l.conns, key)
}

func (l *conns) close()  {
	for _, conn := range l.conns {
		conn.DB.Close()
	}
	for _, key := range l.keys {
		delete(l.conns, key)
	}
}