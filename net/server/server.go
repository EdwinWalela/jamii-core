package server

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"unsafe"
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
		fmt.Println("server:Client " + client.RemoteAddr().String() + " connected.")
		tmp := make([]byte, 1024)
		client.Read(tmp)

		decoded := strings.Split(BytesToString(tmp), "\n")
		decoded = decoded[1 : len(decoded)-1]
		for _, v := range decoded {
			headerData := strings.Split(v, ":")

			headerExtract := strings.Replace(strings.Join(headerData, ","), " ", "", -1)
			headerList := strings.Split(headerExtract, ",")
			for i, val := range headerList {
				switch val {
				case "Source":
					fmt.Println("source:" + headerList[i+1])
				case "Type":
					fmt.Println("msg-type:" + headerList[i+1])
				case "Data":
					fmt.Println("data:" + headerList[i+1])
				}
			}

		}
		server.Clients = append(server.Clients, client)

	}
}

func (server *Server) Read() {
	for _, client := range server.Clients {
		// fmt.Println("Recieved block from peer")
		tmp := make([]byte, 64)
		client.Read(tmp)
		fmt.Println(tmp)
	}
}
func (server *Server) Connections() int {
	return len(server.Clients)
}

func (server *Server) HandleConnection() {

}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
