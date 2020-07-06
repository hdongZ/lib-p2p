package p2p

import (
	p2p_peer "github.com/libp2p/go-libp2p-core/peer"
	_ "github.com/multiformats/go-multiaddr"
)

type Peer struct {
	Addrs []Multiaddr
}