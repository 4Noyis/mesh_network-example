# P2P Networking in Go

## Overview

This project demonstrates a simple Peer-to-Peer (P2P) networking implementation in Go. Each peer can connect to multiple other peers, send its address to newly connected peers, and broadcast new peer addresses to all connected peers. The application uses TCP for communication between peers.

## Features

- Establish TCP connections between multiple peers.
- Broadcast new peer connections to all existing peers.
- Handle incoming connections from peers.

## Requirements

- Go 1.16 or higher

## Installation

- Clone the repository:

   ```bash
   git clone https://github.com/4Noyis/p2p_network-example.git
   cd p2p_network-example

## Usage

To run the P2P networking application, you need to start multiple instances of the program, each acting as a peer. Each peer will listen on a specified port and can connect to other peers.

### Starting a Peer

Run the following command in separate terminal windows for each peer you want to start. Replace `<address>` with the desired address and port:

```bash
go run main.go <own_address> [<other_peer_adress>...]
