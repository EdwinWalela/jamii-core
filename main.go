package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/edwinwalela/jamii-core/net/peer"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

const (
	ON_CLIENT_VOTE         = "vote"
	ON_CLIENT_LATEST_BLOCK = "latest-block"
	ON_BLOCK_HEIGHT        = "block-height"
	ON_BLOCK_AT_HEIGHT     = "block-at-height"
	PEER_BLOCK_BROADCAST   = "peer-block-broadcast"
)

var exit = make(chan int)

func main() {

	/** Key pair generation and signing **/

	// kp := &jcrypto.KeyPair{}

	// jcrypto.GenKeyPair(kp, "") // Generate key pair

	// hash := jcrypto.SHA512("hello world") // hash data

	// signature, err := kp.Sign(hash) // sign hash

	// if err != nil {
	// 	log.Println("signing failed:", err)
	// }

	// valid := jcrypto.VerifySig(signature, []byte(hash), kp) // verify hash with signature

	// fmt.Println(valid)

	// kp := &jcrypto.KeyPair{}

	// decodedPub, err := base64.StdEncoding.DecodeString("oMIwjMspTa8oBTsc/0qJ/GUvc6Fa/MX78F2lwkNdePs=")

	// decodedSig, sErr := base64.StdEncoding.DecodeString("rgBcyHuuyeB9CgwtiG/+PG2JqRZLf3lPpsD5fSgANnQlQBrX78db874Xys5f/JnjZccXVMaUvqpzO3F/utatDw==")

	// if err != nil {
	// 	log.Println(err)
	// }
	// if sErr != nil {
	// 	log.Println(err)
	// }

	// jcrypto.PubKeyFromBytes(decodedPub, kp)

	// valid := jcrypto.VerifySig(decodedSig, []byte("hello world"), kp)

	// fmt.Println(valid)

	/** Vote Unpacking **/
	// v := &primitives.Vote{}
	// clientData := "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f|Z/Yg3ETYDRBlXcDy09p/NOyQRGycYRU8kNPsTRkwqRYpiz8ljL87u8fg/x/xzpELh65Af95kIeV2zoV++CbiDw==|rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=|rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=,rLPdtJDVywkcufc0kSQpmGnQ8sc6frPvVbMWYfwkFj0=|1623142046"
	// v.UnpackClientString(clientData)

	// fmt.Println("Address: ", v.Address)
	// fmt.Println("Candidates: ", v.Candidates)
	// fmt.Println("Signature: ", v.Signature)
	// fmt.Println("Hash:", v.Hash)
	// fmt.Println("Timestamp: ", v.Timestamp)

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	localPortPtr := flag.String("local", "3000", "local socket server")
	// tunnelUrlPtr := flag.String("tunnel", "", "Local tunnel URL (Ngrok)")

	flag.Parse()

	// if *tunnelUrlPtr == "" {
	// 	log.Fatal("Local Tunnel URL not provided")
	// }

	// Send tunnelURL to server for storage

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {

		source := c.RequestHeader().Get("source")

		switch source {
		case "peer":
			err := c.Join("peers")
			log.Println(err)
			log.Printf("%s added to peers channel\n", c.Ip())
			log.Printf("total members: %d\n", c.Amount("peers"))
		case "client":
			c.Join("clients")
			log.Printf("%s added to clients channel\n", c.Ip())
			log.Printf("total members: %d\n", c.Amount("peers"))
		}
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	// Handle vote message from clients
	server.On(ON_CLIENT_VOTE, func(c *gosocketio.Channel, msg string) string {
		log.Println("Recieved vote")
		voteStr := []byte(msg)
		var voteObj map[string]string

		if err := json.Unmarshal(voteStr, &voteObj); err != nil {
			log.Println(err)
		}

		// Validate vote
		log.Println("Vote accepted")
		// fmt.Println(voteObj["source"])

		return "OK"
		// fmt.Println(v)
	})

	// Send back latest block
	server.On(ON_CLIENT_LATEST_BLOCK, func(c *gosocketio.Channel) {

	})

	// Send back currrent block height
	server.On(ON_BLOCK_HEIGHT, func(c *gosocketio.Channel) {

	})

	// Send back requested block
	server.On(ON_BLOCK_AT_HEIGHT, func(c *gosocketio.Channel) {

	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	go func() {
		log.Printf("Starting server on port %s...\n", *localPortPtr)
		log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", *localPortPtr), serveMux))
	}()

	// Try to Connect to peers from server and store their connections
	peers := []string{
		"localhost:4000",
		"a.com",
		"b.com",
		"c.com",
		"d.com",
	}

	fmt.Scanln() // Block

	for i := range peers { // Attempt to connect to peers from server

		go func(url *string) {
			p := peer.PeerConnection{Host: *url}
			var err error
			rawUrl := strings.Split(*url, ":")

			if len(rawUrl) > 1 { // Extract port from URL
				p.Port, _ = strconv.Atoi(rawUrl[1])
				p.Host = rawUrl[0]
			} else {
				p.Port = 3000 // Set Default port
			}

			p.Init()
			p.SetSource("peer")
			c, err := p.Dial() // Attempt to connect to peer
			if err != nil {
				log.Printf("Peer %s not found", *url)
				return
			}

			// Accept new block broadcast from peer, check and add to local chain
			c.On(PEER_BLOCK_BROADCAST, func(h *gosocketio.Channel, args string) {
				log.Println("c.onblock called", args)
			})

			// Accept new block from peer, check and add to local chain
			c.On(PEER_BLOCK_BROADCAST, func(h *gosocketio.Channel, args string) {
				log.Println("received block 1 from peer", args)
			})

		}(&peers[i])

	}

	// Generate Key Pair

	// Attempt to connect to peers

	// Read Blocks from memory

	// If no blocks request from peers

	// Check chain

	// Initalize chain

	// broadcast mined block to peers
	server.BroadcastTo("peers", "block", "here's the bloc")

	// request block from a (random) peer
	server.List("peers")[0].Emit("block-request", "1")

	exit <- 1

}
