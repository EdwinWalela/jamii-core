package peer

import (
	"fmt"
	"log"
	"net"
)

type Peer struct {
	conn      net.Conn
	Host      string
	Network   string
	Port      int
	Connected bool
}

func (peer *Peer) Bind() error {
	var err error

	peer.conn, err = net.Dial(peer.Network, peer.Host+":"+fmt.Sprint(peer.Port))

	if err != nil {
		return err
	}
	peer.Connected = true
	fmt.Println("client:Connected to remote on:" + peer.conn.RemoteAddr().String())

	return nil
}

func (peer *Peer) Write() {
	fmt.Println("Sending my block to peer")
	_, err := peer.conn.Write([]byte("hello world"))
	if err != nil {
		fmt.Println("Error", err)
	}
}

func (peer *Peer) Read() {
	var buf = make([]byte, 4096)
	var tmp = make([]byte, 256)
	n, err := peer.conn.Read(tmp)

	if err != nil {
		log.Println("Pread:", err)
	}

	buf = append(buf, tmp[:n]...)

	fmt.Println("Read from server")
}
