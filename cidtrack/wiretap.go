package cidtrack

import (
	"sync"

	bsmsg "github.com/ipfs/go-bitswap/message"
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-peer"
)

// WireTap implements go-bitswap WireTap interface
type wireTap struct {
	stats   map[cid.Cid]uint64
	statsCh chan cid.Cid
	mtx     sync.RWMutex
}

func newWireTap() *wireTap {
	return &wireTap{
		stats:   make(map[cid.Cid]uint64),
		statsCh: make(chan cid.Cid),
	}
}

func (t *wireTap) run() {
	for c := range t.statsCh {
		t.mtx.Lock()
		t.stats[c]++
		t.mtx.Unlock()
	}
}

func (t *wireTap) reset() {
	t.mtx.Lock()
	t.stats = make(map[cid.Cid]uint64)
	t.mtx.Unlock()
}

func (t *wireTap) stop() {
	close(t.statsCh)
}

func (t *wireTap) MessageReceived(p peer.ID, msg bsmsg.BitSwapMessage) {
	// We're not interested in received messages
}

func (t *wireTap) MessageSent(p peer.ID, msg bsmsg.BitSwapMessage) {
	for _, block := range msg.Blocks() {
		t.statsCh <- block.Cid()
	}
}
