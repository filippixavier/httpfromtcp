package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
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

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\r\n"
	conn.Write([]byte(response))
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
		isClosed: atomic.Bool{},
	}

	server.isClosed.Store(false)

	go server.listen()

	return &server, nil
}
