package server

import (
	"fmt"
	"net"
)

type Server struct {
	conn    net.Listener
	Host    string
	Network string
	Port    int
	Clients []net.Conn
}

func (server *Server) Init() error {
	var err error
	server.conn, err = net.Listen(server.Network, server.Host+":"+fmt.Sprint(server.Port))

	if err != nil {
		return err
	}
	fmt.Println("server:Waiting for client connections on port:", server.Port)
	return nil
}

func (server *Server) Accept() error {
	for {
		client, err := server.conn.Accept()

		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return err
		}
		fmt.Println("server:Client connected.")
		fmt.Println("server:Client " + client.RemoteAddr().String() + " connected.")
		client.Write([]byte("hello"))
		server.Clients = append(server.Clients, client)
	}
}

func (server *Server) Read() {
	for _, client := range server.Clients {
		fmt.Println("Recieved block from peer")
		client.RemoteAddr()
	}
}
func (server *Server) Connections() int {
	return len(server.Clients)
}

func (server *Server) HandleConnection() {

}
