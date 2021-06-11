package localsocket

import (
	"fmt"
	"log"

	gosocketio "github.com/graarh/golang-socketio"
)

type LocalServer struct {
	Server *gosocketio.Server
}

func (s *LocalServer) HandleNewConnection(c *gosocketio.Channel) {
	log.Println("Connected")
	fmt.Println(c.RequestHeader())
}

func (s *LocalServer) HandleDisconnection(c *gosocketio.Channel) {
	log.Println("Disconnected")
}
