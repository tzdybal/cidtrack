package cidtrack

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ipfs/go-bitswap/message"
	blocks "github.com/ipfs/go-block-format"
	"github.com/test-go/testify/assert"
)

func TestWiretap(t *testing.T) {
	tap := newWireTap()

	go tap.run()
	defer tap.stop()

	for i := 0; i < 10; i++ {
		tap.MessageSent("", getMessageWith10RandomBlocks())
	}

	// Give some time to process thru channel
	time.Sleep(100 * time.Millisecond)

	tap.mtx.RLock()
	assert.Len(t, tap.stats, 10*10)
	tap.mtx.RUnlock()
}

func getMessageWith10RandomBlocks() message.BitSwapMessage {
	msg := message.New(false)
	for i := 0; i < 10; i++ {
		msg.AddBlock(blocks.NewBlock([]byte(fmt.Sprintf("%d:%d", i, rand.Int31()))))
	}

	return msg
}
