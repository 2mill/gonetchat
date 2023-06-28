package main

import "fmt"

const (
	HEAD        = "[%s]"
	DESCRIPTION = ">> %s <<"
)

type ChatChannel struct {
	name        string
	topic       string
	subscribers []User
}

func (self *ChatChannel) UserSubscribe(user User) error {
	for _, suser := range self.subscribers {
		if suser.name == user.name {
			return fmt.Errorf("%s is already subscribed to %s", user.name, self.name)
		}
	}
	subscribers := append(self.subscribers, user)
	self.subscribers = subscribers
	return nil
}

func (self ChatChannel) ChannelHeader() Message {
	header := fmt.Sprintf(HEAD, self.name)
	description := fmt.Sprintf(DESCRIPTION, self.topic)
	return Message{
		content: fmt.Sprintln(header) + fmt.Sprintln(description),
	}
}

func (self ChatChannel) FromServer(server Server) Server {
	return Server{
		name:     self.name,
		ip:       server.ip,
		port:     server.port,
		listener: nil,
		users:    make([]User, 0),
	}
}
