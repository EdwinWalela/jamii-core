package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/edwinwalela/jamii-core/jcrypto"
	"github.com/edwinwalela/jamii-core/primitives"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

const (
	ON_CLIENT_VOTE         = "vote"
	ON_CLIENT_REGISTER     = "register"
	ON_CLIENT_RES_REQ      = "req"
	VOTE_ACK               = "VOTE_ACK"
	VOTE_INVALID           = "VOTE_INV"
	ON_CLIENT_LATEST_BLOCK = "latest-block"
	ON_BLOCK_HEIGHT        = "block-height"
	ON_BLOCK_AT_HEIGHT     = "block-at-height"
	PEER_BLOCK_BROADCAST   = "peer-block-broadcast"
	KEY_FILE               = "key.jpkey"
	MIN_DIFFICULTY         = 3
	MAX_BLOCK_SIZE         = 1
	BLOCK_DIR              = "/data/blocks"
)

var exit = make(chan int)

func main() {
	localPortPtr := flag.String("local", "4000", "local socket server")
	// tunnelUrlPtr := flag.String("tunnel", "", "Local tunnel URL (Ngrok)")

	flag.Parse()

	log.Println("------------------------------------------")
	log.Println("Key pair initialization")
	log.Println("------------------------------------------")
	log.Println("Looking for key pair")
	/** Key pair generation **/

	kp := &jcrypto.KeyPair{}

	_, err := ioutil.ReadFile(KEY_FILE)

	if err != nil {

		log.Println("Private key not found in directory, New KeyPair generated")

		if err := jcrypto.GenKeyPair(kp); err != nil {
			log.Println(err)
		}

		if err := ioutil.WriteFile(KEY_FILE, kp.PrivKey, 0644); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := jcrypto.ReadKeyPair(kp, KEY_FILE); err != nil {
			log.Fatal(err)
		}

		log.Println("Key pair found")
	}

	// Create blocks directory
	path := filepath.Join(".", BLOCK_DIR)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
	log.Println("------------------------------------------")
	log.Println("Initializing chain")
	log.Println("------------------------------------------")
	// Initialize chain
	diff, nonce, elapsed := jcrypto.FindDifficulty()

	// Ensure Chain's min Difficulty is met
	if diff < MIN_DIFFICULTY {
		diff = MIN_DIFFICULTY
	}

	diff = 2

	jchain := &primitives.Chain{Difficulty: diff, BlockDir: BLOCK_DIR}

	if chainInitError := jchain.Init(); chainInitError != nil {
		log.Fatal(chainInitError)
	}

	log.Printf("Chain initialized in %d seconds with: Diff:%d, Nonce:%d\n", elapsed, diff, nonce)
	log.Println("------------------------------------------")

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connection recieved from %s", c.Ip())

	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Lost connection to %s", c.Ip())
		c.Close()
	})

	// Handle registration meggage from clients
	server.On(ON_CLIENT_REGISTER, func(c *gosocketio.Channel, msg string) string {
		log.Printf("Received registration packet from %s\n", c.Ip())
		regsterStr := []byte(msg)

		var registerObj map[string]string

		if err := json.Unmarshal(regsterStr, &registerObj); err != nil {
			log.Println(err)
		}
		c.Join("clients")
		// Validate vote
		data := registerObj["data"]
		v := &primitives.Vote{}

		for i, val := range strings.Split(data, "|") {
			v.Timestamp = uint64(time.Now().Unix())
			switch i {
			case 0: // Extract hash
				v.Hash = val
			case 1: // extract pubkey
				decodedSig, sigErr := base64.StdEncoding.DecodeString(val)

				if sigErr != nil {
					log.Println(sigErr)
				}

				v.Address = decodedSig
			case 2: // Extract signature (base64 encoded)
				decodedPub, pubErr := base64.StdEncoding.DecodeString(val)
				if pubErr != nil {
					log.Println(pubErr)
				}

				v.Signature = decodedPub

			default:
			}

		}
		jchain.AddTX(*v)
		if len(jchain.PendingVotes) >= MAX_BLOCK_SIZE {

			if err := jchain.Mine(kp); err != nil {
				log.Println(err)
			}
		}
		if jchain.LatestBlock().Votes[1].Hash == v.Hash {
			c.Emit(VOTE_ACK, "")
		} else {

		}
		return "OK"
	})

	// // Handle vote message from clients
	server.On(ON_CLIENT_VOTE, func(c *gosocketio.Channel, msg string) string {
		log.Printf("Received vote packet from %s\n", c.Ip())
		voteStr := []byte(msg)

		var voteObj map[string]string

		if err := json.Unmarshal(voteStr, &voteObj); err != nil {
			log.Println(err)
		}

		// Validate vote
		data := voteObj["data"]
		log.Println("Unpacking vote")
		log.Println("------------------------------------------")
		log.Println(data)
		log.Println("------------------------------------------")

		v := &primitives.Vote{}
		for i, val := range strings.Split(data, "|") {
			switch i {
			case 0: // Extract hash
				v.Hash = val
			case 1: // extract signature
				decodedSig, sigErr := base64.StdEncoding.DecodeString(val)

				if sigErr != nil {
					log.Println(sigErr)
				}

				v.Signature = decodedSig
			case 2: // Extract publickey (base64 encoded)
				decodedPub, pubErr := base64.StdEncoding.DecodeString(val)
				if pubErr != nil {
					log.Println(pubErr)
				}

				v.Address = decodedPub
			case 3: // Extract candidate names
				for _, candidate := range strings.Split(val, ".") {

					v.Candidates = append(v.Candidates, candidate)
				}
			case 4: // Extract timestamp
				v.Timestamp, err = strconv.ParseUint(val, 10, 64)
				if err != nil {
					log.Println(err)
				}
			default:
			}
		}
		if v.IsValid() {
			log.Println("Vote valid. Added to pending tx")
			jchain.AddTX(*v)
		} else {
			log.Println("Vote invalid. Discarded")

		}
		if len(jchain.PendingVotes) >= MAX_BLOCK_SIZE {

			if err := jchain.Mine(kp); err != nil {
				log.Println(err)
			}
		}
		return "OK"
	})

	server.On(ON_CLIENT_RES_REQ, func(c *gosocketio.Channel, msg string) string {
		log.Printf("Received result query packet from %s\n", c.Ip())

		c.Emit("result", "1234")
		return "OK"
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
		log.Println("------------------------------------------")
		log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", *localPortPtr), serveMux))
	}()

	// Try to Connect to peers from server and store their connections
	// peers := []string{
	// 	"localhost:4000",
	// 	"a.com",
	// 	"b.com",
	// 	"c.com",
	// 	"d.com",
	// }

	// fmt.Scanln() // Block

	// for i := range peers { // Attempt to connect to peers from server

	// 	go func(url *string) {
	// 		p := peer.PeerConnection{Host: *url}
	// 		var err error
	// 		rawUrl := strings.Split(*url, ":")

	// 		if len(rawUrl) > 1 { // Extract port from URL
	// 			p.Port, _ = strconv.Atoi(rawUrl[1])
	// 			p.Host = rawUrl[0]
	// 		} else {
	// 			p.Port = 3000 // Set Default port
	// 		}

	// 		p.Init()
	// 		p.SetSource("peer")
	// 		c, err := p.Dial() // Attempt to connect to peer
	// 		if err != nil {
	// 			log.Printf("Peer %s not found", *url)
	// 			return
	// 		}

	// 		// Accept new block broadcast from peer, check and add to local chain
	// 		c.On(PEER_BLOCK_BROADCAST, func(h *gosocketio.Channel, args string) {
	// 			log.Println("c.onblock called", args)
	// 		})

	// 		// Accept new block from peer, check and add to local chain
	// 		c.On(PEER_BLOCK_BROADCAST, func(h *gosocketio.Channel, args string) {
	// 			log.Println("received block 1 from peer", args)
	// 		})

	// 	}(&peers[i])

	// }

	// Generate Key Pair

	// Attempt to connect to peers

	// Read Blocks from memory

	// If no blocks request from peers

	// Check chain

	// Initalize chain

	// broadcast mined block to peers
	// server.BroadcastTo("peers", "block", "here's the bloc")

	// request block from a (random) peer
	// server.List("peers")[0].Emit("block-request", "1")

	exit <- 1

}
