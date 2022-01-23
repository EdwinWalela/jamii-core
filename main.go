package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/edwinwalela/jamii-core/jcrypto"
	"github.com/edwinwalela/jamii-core/net/peer"
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
	DATA_SYNC              = "data-sync"
	ON_CLIENT_LATEST_BLOCK = "latest-block"
	ON_BLOCK_HEIGHT        = "block-height"
	ON_BLOCK_AT_HEIGHT     = "block-at-height"
	PEER_BLOCK_BROADCAST   = "peer-block-broadcast"
	KEY_FILE               = "key.jpkey"
	MIN_DIFFICULTY         = 3
	MAX_BLOCK_SIZE         = 1
	BLOCK_DIR              = "/data/blocks"
	PEER_NETWORK           = "https://jamii-peer.herokuapp.com/peers"
)

var peerChannel *gosocketio.Client
var jchain *primitives.Chain
var exit = make(chan int)

// register tunnel url

func registerPeer(tunnelUrl string) error {
	log.Println("registering peer")
	data := url.Values{
		"url": {tunnelUrl},
	}

	_, err := http.PostForm(PEER_NETWORK, data)

	return err
}

// deregister tunnel url

func deregisterPeer(tunnelUrl string) error {
	log.Println("deregistering peer")
	data := url.Values{
		"url": {tunnelUrl},
	}

	_, err := http.PostForm(PEER_NETWORK+"/delete", data)

	return err
}

func main() {
	localPortPtr := flag.String("local", "4000", "local socket server")
	tunnelUrlPtr := flag.String("tunnel", "", "Local tunnel URL (Ngrok)")

	flag.Parse()
	// send tunnel url

	if err := registerPeer(*tunnelUrlPtr); err != nil {
		log.Fatal("Unable to register peer")
	}
	defer deregisterPeer(*tunnelUrlPtr)

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

	jchain = &primitives.Chain{Difficulty: diff, BlockDir: BLOCK_DIR}

	if chainInitError := jchain.Init(); chainInitError != nil {
		log.Fatal(chainInitError)
	}

	log.Printf("Chain initialized in %d seconds with: Diff:%d, Nonce:%d\n", elapsed, diff, nonce)
	log.Println("------------------------------------------")

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connection recieved from %s", c.Ip())
	})

	server.On("data-req", func(c *gosocketio.Channel, msg string) string {
		log.Println("Data sync initiated")
		c.Emit(DATA_SYNC, "block1,block2,block3")
		return "OK"
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
		server.BroadcastToAll(ON_CLIENT_REGISTER, data)
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

	// Handle vote message from clients
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
			server.BroadcastToAll(ON_CLIENT_VOTE, data)

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

	// Results request
	server.On(ON_CLIENT_RES_REQ, func(c *gosocketio.Channel, msg string) string {
		log.Printf("Received result query packet from %s\n", c.Ip())
		res := jchain.Result()
		log.Printf("result:%s", res)
		c.Emit("result", res)
		return "OK"
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	go func() {
		log.Printf("Starting server on port %s...\n", *localPortPtr)
		log.Println("------------------------------------------")
		log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", *localPortPtr), serveMux))
	}()

	// Try to Connect to peers from server and store their connections
	peers := []string{
		"localhost:3000",
	}
	if *localPortPtr == "4000" {

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
				peerChannel, err = p.Dial() // Attempt to connect to peer
				if err != nil {
					log.Printf("Peer %s not found", *url)
					return
				}

				peerChannel.Emit("data-req", "")

				// Accept new block broadcast from peer, check and add to local chain
				peerChannel.On(ON_CLIENT_VOTE, func(h *gosocketio.Channel, data string) {
					log.Println("recived vote packet from peer")
					log.Println("------------------------------------------")
					log.Println("Unpacking data")
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
				})

				peerChannel.On(ON_CLIENT_REGISTER, func(h *gosocketio.Channel, data string) {
					log.Println("recived registration packet from peer")
					log.Println("------------------------------------------")
					log.Println("Unpacking data")
					log.Println("------------------------------------------")
					log.Println(data)
					log.Println("------------------------------------------")
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

				})

				// Accept new block from peer, check and add to local chain
				peerChannel.On(DATA_SYNC, func(h *gosocketio.Channel, args string) {
					log.Println("Recieved block from remote", args)

				})

			}(&peers[i])

		}
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Println(sig.String())
			deregisterPeer(*tunnelUrlPtr)
			os.Exit(1)
		}
	}()
	exit <- 1
	// deregister tunnel url

}
