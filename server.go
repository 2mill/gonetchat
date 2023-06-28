package main

import (
	"fmt"
	"log"
	"net"
)

func SpinUpServer(server Server) (Server, error) {
	if server.listener != nil {
		return server, fmt.Errorf("Server %s is already listening.", server.name)
	}
	server.open()
	return server, nil
}

type Server struct {
	name     string
	ip       string
	port     string
	listener net.Listener
	users    []User
}

func (self *Server) teardown() error {
	return self.listener.Close()
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

func NewServer(name string, ip string, port string) Server {
	return Server{
		name:     name,
		ip:       ip,
		port:     port,
		listener: nil,
		users:    make([]User, 0),
	}

}

func (self *Server) get_random_username() string {
	return fmt.Sprintf("God%v", len(self.users))
}
