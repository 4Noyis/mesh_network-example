package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// peer struct
type Peer struct {
	Address   string              // the peer's own address (IP:port)
	Peers     map[string]net.Conn // Connections to other peers  {key -> address : value -> conn(TCP connection)}
	PeersLock sync.Mutex          // Mutex to protect access to peer list
}

// create new peer
func NewPeer(address string) *Peer {
	return &Peer{
		Address: address,
		Peers:   make(map[string]net.Conn),
	}
}

// Add a new peer connection
func (p *Peer) AddPeer(address string) {
	// ensure that only one goroutine can access the Peers map at time. Making it safe for concurrent use.
	p.PeersLock.Lock()
	defer p.PeersLock.Unlock()

	// If connection already exists (exists is true), the function exits early, preventing duplicate connections.

	if _, exists := p.Peers[address]; exists {
		return
	}

	// Esteblish a new connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Failed to connect to peer %s: %v\n", address, err)
		return
	}

	p.Peers[address] = conn
	fmt.Printf("Connected to peer %s\n", address)

	// Send own adress to identify this peer
	_, _ = conn.Write([]byte(p.Address + "\n"))

	p.DisplayConnectios()

}

// Accept incoming connections
func (p *Peer) AcceptConnection() {
	listener, err := net.Listen("tcp", p.Address)
	if err != nil {
		fmt.Printf("Failed to start listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Accepting Connections...\n")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go p.HandleConnection(conn)
	}
}

func (p *Peer) DisplayConnectios() {
	fmt.Println("Current connections to peers:")
	for address := range p.Peers {
		fmt.Printf(" - %s\n", address)
	}
}

// handle incoming connections
func (p *Peer) HandleConnection(conn net.Conn) {
	defer conn.Close()

	// read the addres of the new peer
	reader := bufio.NewReader(conn)
	peerAddr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to get peer address:", err)
		return
	}

	// moves any leading or trailing whitespace (including newline characters) from peerAddr. (\n)
	peerAddr = strings.TrimSpace(peerAddr)

	// Add the new peer
	p.AddPeer(peerAddr)

	// Broadcast the new peer's address to all connected peers
	p.PeersLock.Lock()
	for addr, peerConn := range p.Peers {
		if addr != peerAddr { // Avoid broadcasting to the same peer twice
			_, _ = peerConn.Write([]byte(peerAddr + "\n"))
		}
	}
	p.PeersLock.Unlock()
}

// CLI-based connection to other peers
func (p *Peer) ConnectToPeers(addresses []string) {
	for _, addr := range addresses {

		p.AddPeer(addr)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <own_address> [<other_peer_adress>...]")
		return
	}

	ownAdress := os.Args[1]
	peer := NewPeer(ownAdress)

	time.Sleep(time.Second)

	// start accepting connections
	go peer.AcceptConnection()

	// Connect to other peers
	if len(os.Args) > 2 {
		peersAdresses := os.Args[2:] // including 2 and the ones after that
		peer.ConnectToPeers(peersAdresses)
	}

	// keep program running its waiting forever and because of that main function working until we stop it manually
	select {}

}
