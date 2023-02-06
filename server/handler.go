package server

type Handler interface {
	HandleSession(session *Session)
}
