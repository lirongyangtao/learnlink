package server

import "io"

type Protocol interface {
	NewCodec(rw io.ReadWriter) (Codec, error)
}

type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}
