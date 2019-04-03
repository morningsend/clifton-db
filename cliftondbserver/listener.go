package cliftondbserver

import (
	"errors"
	"net"
	"time"
)

type StoppableListener struct {
	*net.TCPListener
	stop          chan int
	listenTimeout time.Duration
}

var StoppedErr = errors.New("listener stopped")

func NewStoppableListener(l net.Listener) (*StoppableListener, error) {
	tcpListener, ok := l.(*net.TCPListener)
	if !ok {
		return nil, errors.New("cannot wrap net.Listener, must be a net.TCPListener")
	}

	sl := &StoppableListener{
		stop:        make(chan int),
		TCPListener: tcpListener,
	}

	return sl, nil
}
func (l *StoppableListener) Accept() (net.Conn, error) {
	for {
		err := l.TCPListener.SetDeadline(time.Now().Add(l.listenTimeout))
		if err != nil {
			return nil, err
		}

		newConn, err := l.TCPListener.Accept()
		select {
		case <-l.stop:
			if err == nil {
				err = newConn.Close()
			}
			return nil, StoppedErr
		default:
		}
		if err != nil {
			netErr, ok := err.(net.Error)

			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, nil
	}
}

func (l *StoppableListener) Stop() {
	close(l.stop)
}
