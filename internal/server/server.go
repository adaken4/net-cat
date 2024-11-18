package server

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var welcomeMsg string = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    `.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n"

type Client struct {
	Conn net.Conn
	Name string
}

type Server struct {
	address string
	clients map[net.Conn]Client
	mutex   sync.Mutex
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		clients: make(map[net.Conn]Client),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Server started at %s\n", s.address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		fmt.Println("New client connected")

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

	conn.Write([]byte(welcomeMsg))

	conn.Write([]byte("[ENTER YOUR NAME]: "))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read name:", err)
		return
	}
	name := string(buf[:n])
	name = strings.TrimSpace(name)

	s.mutex.Lock()
	for _, client := range s.clients {
		if client.Name == name {
			conn.Write([]byte("Name already taken. Disconnecting...\n"))
			s.mutex.Unlock()
			return
		}
	}
	s.clients[conn] = Client{Conn: conn, Name: name}
	s.mutex.Unlock()

	conn.Write([]byte(fmt.Sprintf("Welcome, %s!\n", name)))
	fmt.Printf("%s joined the chat.\n", name)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Connection closed:", err)
			return
		}
		layout := "2006-01-02 15:04:05"
		currentTime := time.Now()
		formattedTime := currentTime.Format(layout)
		message := string(buf[:n])
		fmt.Printf("[%s]: %s", name, message)
		s.broadcastMessage(fmt.Sprintf("[%s][%s]:%s", formattedTime, name, message), conn)
	}
}

func (s *Server) broadcastMessage(message string, sender net.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn, client := range s.clients {
		if conn != sender {
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("Failed to send message to %s: %v\n", client.Name, err)
			}
		}
	}
}
