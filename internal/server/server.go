package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode response.StatusCode
	ErrorMsg   []byte
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h HandlerError) print(w io.Writer) error {
	err := response.WriteStatusLine(w, h.StatusCode)

	if err != nil {
		return err
	}

	headers := response.GetDefaultHeaders(len(h.ErrorMsg))

	err = response.WriteHeaders(w, headers)

	if err != nil {
		return nil
	}

	_, err = fmt.Fprintf(w, "\r\n%s", h.ErrorMsg)

	return err
}

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

	buf := new(bytes.Buffer)

	reqError := s.handler(buf, req)

	if reqError != nil {
		reqError.print(conn)
		return
	}

	// //Header part
	response.WriteStatusLine(conn, response.Ok)
	response.WriteHeaders(conn, response.GetDefaultHeaders(buf.Len()))
	// CRLF => end of headers.
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, buf)
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
