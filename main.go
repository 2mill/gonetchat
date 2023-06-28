package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func accept_any(l net.Listener) (net.Conn, error) {
	return l.Accept()
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

func start_server(server *Server) {

	log.Printf("Starting %s", server.name)
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
	args := os.Args
	server_name := "Olympus"
	if len(args) > 1 {
		server_name = args[1]
	}
	const (
		IP   = "127.0.0.1"
		PORT = "1337"
	)
	server := NewServer(server_name, IP, PORT)
	_ = server.open()
	_ = make(chan struct{})
	go start_server(&server)
	for {
		var console_buff string
		fmt.Scanln(&console_buff)
		if console_buff == "FIN" {
			server.teardown()
			break
		}
	}

}
