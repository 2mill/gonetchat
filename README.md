
`go run .` or `go build . && ./gonetchat` to run.
From a different terminal, connect to the server using your preferred TCP client.

`nc 127.0.0.1 1337` is a good example.

type `FIN` while the server is running to close all sessions.

# !!! Chat server is prone to DDOSing and can cause memory leaks right now. User beware.
Still need to implement graceful teardowns of user sessions.