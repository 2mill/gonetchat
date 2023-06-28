package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func accept_any(l net.Listener) (net.Conn, error) {
	return l.Accept()
}

type User struct {
	name string
	conn net.Conn
}

func (s *User) send(message string) error {
	_, err := s.conn.Write([]byte(message))
	s.conn.Write([]byte(CHAT_PROMPT))
	return err
}

func (s *User) end_session() error {
	return s.conn.Close()
}

func (server *Server) session_user(user User) {

	buff := make([]byte, 1024)
	for {
		n, err := user.conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Printf("%s err: %s", user.name, err)
			}
			continue
		}
		log.Printf("%s sent %s", user.name, buff[:n])
		message := Message{
			content: fmt.Sprintf("\n[%s]: %s", user.name, buff[:n]),
			user:    user,
		}
		server.broadcast(message)
	}

}

func start_server(server Server) {
	for {
		conn, _ := server.listener.Accept()
		username := server.get_random_username()
		log.Printf("%s connected to %s", username, server.name)
		user := User{
			name: username,
			conn: conn,
		}
		server.register_user(user)
		go server.session_user(user)
	}
}

func main() {
	const (
		SERVER_NAME = "Olympus"
		IP          = "127.0.0.1"
		PORT        = "1337"
	)
	server := chat_server(SERVER_NAME, IP, PORT)
	_ = server.open()
	defer func() {
		err := server.teardown()
		if err != nil {
			log.Fatal(err)
		}
	}()
	_ = make(chan struct{})
	go func(server *Server) {
		for {
			conn, _ := server.listener.Accept()
			username := server.get_random_username()
			log.Printf("%s connected to %s", username, server.name)
			user := User{
				name: username,
				conn: conn,
			}
			server.register_user(user)
			go server.session_user(user)
		}
	}(&server)

	for {
		var console_buff string
		fmt.Scanln(&console_buff)
		if console_buff == "FIN" {
			break
		}
	}

}

type Message struct {
	content string
	user    User
}

// Server Messages
const (
	TEARDOWN_NOTICE = "This server is closing.\n"
	NEW_USER_MOTD   = "Welcome to %s %s\n"
	USER_COUNT      = "There are %v users in the chat!\n"
	CHAT_PROMPT     = ">"
)

type Server struct {
	name     string
	ip       string
	port     string
	listener net.Listener
	users    []User
}

func (self *Server) teardown() error {
	for _, user := range self.users {
		user.send(TEARDOWN_NOTICE + "\n")
		// TODO fix this
		user.conn.Close()
	}
	return self.listener.Close()
}

func (self Server) get_random_username() string {
	return fmt.Sprintf("God%v", len(self.users))
}
func (self *Server) register_user(user User) {

	user.send(
		fmt.Sprint(
			fmt.Sprintf(NEW_USER_MOTD, self.name, user.name),
			fmt.Sprintf(USER_COUNT, len(self.users)),
			CHAT_PROMPT,
		),
	)
	self.users = append(self.users, user)
}

func (self *Server) broadcast(message Message) {
	for _, user := range self.users {
		if user.name == message.user.name {
			continue
		}
		err := user.send(message.content)
		if err != nil {
			log.Println(err)
		}
	}
}

func (self *Server) open() error {
	if self.listener != nil {
		return fmt.Errorf("%s has an active connection", self.name)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", self.ip, self.port))
	if err != nil {
		return err
	}
	self.listener = listener

	return nil
}

func chat_server(name string, ip string, port string) Server {
	return Server{
		name:     name,
		ip:       ip,
		port:     port,
		listener: nil,
		users:    make([]User, 0),
	}
}
