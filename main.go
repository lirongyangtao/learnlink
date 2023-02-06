package main

import (
	"github.com/funny/link"
	linkCodec "github.com/funny/link/codec"
	"learnlink/server"
	"learnlink/server/codec"
	"log"
)

type AddReq struct {
	A, B int
}

type AddRsp struct {
	C int
}

type Server struct{}

func main() {
	json := codec.Json()
	json.Register(AddReq{})
	json.Register(AddRsp{})
	server := server.NewServer("127.0.0.1", 8080, new(Server), json)
	go server.Run()

	linkJson := linkCodec.Json()
	linkJson.Register(AddReq{})
	linkJson.Register(AddRsp{})
	clientSession, err := link.Dial("tcp", server.ServerAddr(), linkJson, 1024)
	checkErr(err)
	clientSessionLoop(clientSession)
}

func (*Server) HandleSession(session *server.Session) {
	for {
		req, err := session.Receive()
		checkErr(err)

		err = session.Send(&AddRsp{
			req.(*AddReq).A + req.(*AddReq).B,
		})
		checkErr(err)
	}
}

func clientSessionLoop(session *link.Session) {
	for i := 0; i < 10; i++ {
		err := session.Send(&AddReq{
			i, i,
		})
		checkErr(err)
		log.Printf("Send: %d + %d", i, i)

		rsp, err := session.Receive()
		checkErr(err)
		log.Printf("Receive: %d", rsp.(*AddRsp).C)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
