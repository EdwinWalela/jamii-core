package peer

import (
	"net/http"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type PeerConnection struct {
	Peer        *gosocketio.Client
	Host        string
	Port        int
	SkTransport *transport.WebsocketTransport
}

func (p *PeerConnection) Init() {
	p.SkTransport = transport.GetDefaultWebsocketTransport()
	p.SkTransport.RequestHeader = http.Header{}
}

func (p *PeerConnection) SetSource(source string) {
	p.SkTransport.RequestHeader.Add("source", source)
}

func (p *PeerConnection) Dial() (*gosocketio.Client, error) {
	return gosocketio.Dial(
		gosocketio.GetUrl(p.Host, p.Port, false),
		p.SkTransport)
}
