package types

import "net"

type client struct {
	conn net.Addr
	username string
	room *room
	commands chan <- command
}