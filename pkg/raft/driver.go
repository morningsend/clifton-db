package raft

import (
	"context"
	"log"
	"sync"
	"time"
)

type RaftDriver interface {
	Init(fsm RaftFSM, tickInterval time.Duration)
	Shutdown() <-chan struct{}
	Run(ctx context.Context)
}

type testDriver struct {
	shutdownInitiated int32
	shutdownFinished  chan struct{}
	shutdownOnce      sync.Once
	tickInterval      time.Duration
	fsm               RaftFSM
}

func (d *testDriver) Init(fsm RaftFSM, tickInterval time.Duration) {
	d.fsm = fsm
	d.tickInterval = tickInterval
}

func (d *testDriver) Shutdown() <-chan struct{} {
	go func() {
		d.shutdownOnce.Do(func() {
			time.Sleep(500 * time.Millisecond)
			d.shutdownFinished <- struct{}{}
		})
	}()
	return d.shutdownFinished
}

func (d *testDriver) Run(ctx context.Context) {
	ticker := time.NewTicker(d.tickInterval)
	defer ticker.Stop()
	tickChan := ticker.C

	log.Println(d.fsm.Id(), "tick interval", d.tickInterval)
	for {
		select {
		case <-ctx.Done():
			tickChan = nil
			log.Println(d.fsm.Id(), "done")
			return
		case <-tickChan:
			_ = d.fsm.Tick()
		case msg := <-d.fsm.GetComms().Reply():
			d.fsm.ReceiveMsg(msg)
		}
	}
}

func NewTestDriver() RaftDriver {
	driver := testDriver{
		shutdownFinished: make(chan struct{}),
	}

	return &driver
}
