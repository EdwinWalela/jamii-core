package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/edwinwalela/jamii-core/jcrypto"
	"github.com/edwinwalela/jamii-core/net/peer"
)

var clear map[string]func()     //create a map for storing OS specific clear funcs
var connectedPeers *[]peer.Peer // Store connections to connected peers
var kp *jcrypto.KeyPair         // Node's KeyPair

func MainMenu(peerList *[]peer.Peer) {
	connectedPeers = peerList
	fmt.Println("----------------\n Jamii Core v0.01\n----------------")
	GenPair()       // generate pair or read from file
	DiscoverPeers() // attempt to connect to peers

}

func GenPair() {
	fmt.Println("Generating key pair...")

	kp = &jcrypto.KeyPair{}

	jcrypto.GenKeyPair(kp, "")

	fmt.Println("Key Pair Generated...")
}

func DiscoverPeers() {
	fmt.Printf("Retrieving peers from Server...\n")

	time.Sleep(time.Second * 4) // Mock http request to server

	peers := []*peer.Peer{ // Get Peers from Server
		{Host: "localhost", Network: "tcp", Port: 4000},
		{Host: "http:a.com", Network: "tcp", Port: 5000},
		{Host: "http:b.com", Network: "tcp", Port: 5000},
	}

	fmt.Printf("Found %d peers...\nAttempting to establish connection...\n\n", len(peers))

	// Attempt to establish connection with peer(s)
	for _, peer := range peers {
		if err := peer.Bind(); err != nil {
			// fmt.Println("Unable to connect to host:", err)
		} else {
			*connectedPeers = append(*connectedPeers, *peer)
		}
	}
	if len(*connectedPeers) != 0 {

		fmt.Printf("Connected to %d peers...\n", len(*connectedPeers))
	} else {
		fmt.Println("No peers online.")
	}
}

func BroadcastToPeers() {
	for _, peer := range *connectedPeers {
		peer.Write()
	}
}

func initClear() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func callClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
