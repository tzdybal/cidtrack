package cidtrack

import (
	bsmsg "github.com/ipfs/go-bitswap/message"
	cid "github.com/ipfs/go-cid"
	peer "github.com/libp2p/go-libp2p-peer"
)

type WireTap struct {
	consumer statConsumer
}

type stat struct {
	cid  cid.Cid
	size int
}

type statConsumer interface {
	recvChan() chan stat
	sentChan() chan stat
}

func NewWireTap(c statConsumer) *WireTap {
	return &WireTap{consumer: c}
}

func (t *WireTap) MessageReceived(p peer.ID, msg bsmsg.BitSwapMessage) {
	for _, block := range msg.Blocks() {
		t.consumer.recvChan() <- stat{block.Cid(), len(block.RawData())}
	}
}

func (t *WireTap) MessageSent(p peer.ID, msg bsmsg.BitSwapMessage) {
	for _, block := range msg.Blocks() {
		t.consumer.sentChan() <- stat{block.Cid(), len(block.RawData())}
	}
}
