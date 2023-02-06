package server

import (
	"errors"
	"sync"
	"sync/atomic"
)

var HasCloseError = errors.New("has closed")

type Session struct {
	Id        string
	mgr       *SessionManger
	codec     Codec
	sendCh    chan interface{}
	closeFlag int32
	lock      sync.Mutex
}

func (s *Session) IsClose() bool {
	return atomic.LoadInt32(&s.closeFlag) == 1
}

func (s *Session) Receive() (msg interface{}, err error) {
	s.lock.Lock()
	msg, err = s.codec.Receive()
	s.lock.Unlock()
	return
}

func (s *Session) Send(msg interface{}) error {
	s.sendCh <- msg
	return nil
}

func (s *Session) SendLoop() {
	for {
		select {
		case msg, ok := <-s.sendCh:
			if !ok || s.IsClose() {
				return
			}
			s.lock.Lock()
			err := s.codec.Send(msg)
			s.lock.Unlock()
			if err != nil {
				s.Close()
			}
		}
	}
}

func (s *Session) Close() error {
	if s.IsClose() {
		return HasCloseError
	}
	atomic.AddInt32(&s.closeFlag, 1)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.codec.Close()
	close(s.sendCh)
	return nil
}

type SessionManger struct {
	lock     sync.Mutex
	sessions map[string]*Session
}

func (mgr *SessionManger) NewSession(codec Codec) *Session {
	s := &Session{
		codec:  codec,
		sendCh: make(chan interface{}),
	}
	go s.SendLoop()
	return s
}

func NewSessionManger() *SessionManger {
	return &SessionManger{
		sessions: map[string]*Session{},
	}
}
func (mgr *SessionManger) DelSession(id string) {
	mgr.lock.Lock()
	delete(mgr.sessions, id)
	mgr.lock.Unlock()
}

func (mgr *SessionManger) Dispose() {
	mgr.lock.Lock()
	for _, v := range mgr.sessions {
		v.Close()
	}
	mgr.lock.Unlock()
}
