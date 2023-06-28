package main

// Server Messages
const (
	TEARDOWN_NOTICE = "This server is closing."
	NEW_USER_MOTD   = "Welcome to %s %s"
	USER_COUNT      = "There are %v users in the chat!"
	CHAT_PROMPT     = ">"
)

type Message struct {
	content string
	user    User
}
