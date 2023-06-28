package main

import "net"

type User struct {
	name string
	conn net.Conn
}

func (s *User) send(message string) error {
	_, err := s.conn.Write([]byte(message))
	s.conn.Write([]byte(CHAT_PROMPT))
	return err
}

func (s *User) close() error {
	return s.conn.Close()
}
