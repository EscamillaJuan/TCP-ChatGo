package main

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
			case CMD_USERNAME:
				s.username(cmd.client, cmd.args)
			case CMD_JOIN:
				s.join(cmd.client, cmd.args)
			case CMD_ROOMS:
				s.listRooms(cmd.client, cmd.args)
			case CMD_MSG:
				s.msg(cmd.client, cmd.args)
			case CMD_QUIT:
				s.quit(cmd.client, cmd.args)
			default:
				log.Printf("unknown command: %v", cmd)
			}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client has connected: %s", conn.RemoteAddr().String())
	c := &client{
		conn: conn,
		username: "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func (s *server) username(c *client, args []string) {
	c.username = args[1]
	c.msg(fmt.Sprintf("username set as %s", c.username))
}

func (s *server) join(c *client, args []string) {
	r, ok := s.rooms[args[1]]
	if !ok {
		r = &room {
			name: args[1],
			members: make(map[net.Addr]*client),
		}
		s.rooms[args[1]] = r
	}

	r.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.username))
}

func (s *server) listRooms(c *client, args []string) {}

func (s *server) msg(c *client, args []string) {}

func (s *server) quit(c *client, args []string) {}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.username))
	}
}