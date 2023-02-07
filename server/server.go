package server

import (
	"fmt"
	"net"
)

type Server struct {
	handler  Handler
	protocol Protocol
	mgr      *SessionManger
	Listener net.Listener
	Addr     string
	Port     int
}

func NewServer(Addr string, port int, handler Handler, protocol Protocol) *Server {
	return &Server{
		Addr:     Addr,
		Port:     port,
		mgr:      NewSessionManger(),
		handler:  handler,
		protocol: protocol,
	}
}
func (server *Server) ServerAddr() string {
	return fmt.Sprintf("%v:%v", server.Addr, server.Port)
}
func (server *Server) Run() {
	listener, err := net.Listen("tcp", server.ServerAddr())
	if err != nil {
		panic(err)
	}
	server.Listener = listener
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			codec, err := server.protocol.NewCodec(conn)
			if err != nil {
				codec.Close()
			}
			session := server.mgr.NewSession("", codec)
			server.handler.HandleSession(session)
		}()
	}
}

func (server *Server) GetSession(id string) {

}

func (server *Server) Stop() {
	server.mgr.Dispose()
	server.Listener.Close()
}
