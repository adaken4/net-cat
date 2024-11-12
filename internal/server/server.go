package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	address string
	clients map[net.Conn]bool
	mutex   sync.Mutex
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		clients: make(map[net.Conn]bool),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("Server started at", s.address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		s.mutex.Lock()
		s.clients[conn] = true
		s.mutex.Unlock()

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		s.mutex.Lock()
		delete(s.clients, conn)
		s.mutex.Unlock()
		conn.Close()
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}
		message := string(buf[:n])
		s.broadcastMessage(message, conn)
	}
}

func (s *Server) broadcastMessage(message string, sender net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}
