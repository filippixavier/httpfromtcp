package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Handler func(w response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

func (s *Server) listen() {
	for {
		connection, err := s.listener.Accept()

		if s.isClosed.Load() {
			break
		}

		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
			continue
		}

		go s.handle(connection)
	}
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		return
	}

	writer := response.NewWriter(conn)

	s.handler(writer, req)
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
		isClosed: atomic.Bool{},
		handler:  handler,
	}

	server.isClosed.Store(false)

	go server.listen()

	return &server, nil
}
