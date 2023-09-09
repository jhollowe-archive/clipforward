package cmd

import "time"

const (
	CLIPBOARD_INTERVAL = time.Second
	BUFF_SIZE          = 4096
	CLIPBOARD_CHARS    = 2048
	BUFF_TIMEOUT       = time.Second
)

const (
	DISPLAY = "DISPLAY"
	PING    = "PING"
	PONG    = "PONG"
	BYE     = "BYE"
	CLEAR   = ""
	ACK     = "ACK"
)

const (
	MISC_ERR         = 1
	SERVER_NOT_FOUND = 2
)
