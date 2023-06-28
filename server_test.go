package main

import "testing"

func make_server(t *testing.T) {
	name := "foo"
	port := "0"
	ip := "127.0.0.1"
	server := NewServer(name, ip, port)
	 = server.open()
	t.Log(server.name)
}
