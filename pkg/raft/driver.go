package raft

import (
	"sync"
	"time"
)

type Driver interface {
	Init()
	Shutdown() <-chan struct{}
	Run()
}

type testDriver struct {
	shutdownInitiated int32
	shutdownFinished  chan struct{}
	shutdownOnce      sync.Once
}

func (d *testDriver) Init() {

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

func (d *testDriver) Run() {

}

func NewTestDriver() Driver {
	driver := testDriver{
		shutdownFinished: make(chan struct{}),
	}

	return &driver
}
