package server

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/edwinwalela/jamii-core/net/peer"
)

type Server struct {
	conn    net.Listener
	Host    string
	Network string
	Port    int
	Clients []net.Conn // Mobile Clients
	Peers   []net.Conn // Node Peers
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

func (server *Server) Accept(peerList *[]peer.Peer) error {

	for {
		client, err := server.conn.Accept()
		var source string
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return err
		}
		fmt.Println("server:Client " + client.RemoteAddr().String() + " connected.")
		tmp := make([]byte, 1024)
		client.SetReadDeadline(time.Now().Add(time.Second * 2))
		i, e := client.Read(tmp)
		fmt.Println(i, ":", e)
		decoded := strings.Split(BytesToString(tmp), "\n")
		if len(decoded) <= 1 {
			*peerList = append(*peerList, peer.Peer{Conn: client})
			continue
		}
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
		if source == "Client" {
			server.Clients = append(server.Clients, client)
		}

	}
}

func (server *Server) Read() {
	for _, client := range server.Clients {
		fmt.Println("Recieved block from peer")
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
