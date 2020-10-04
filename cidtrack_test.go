package cidtrack

import (
	"net/http"
	"testing"
	"time"

	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-ipfs/core"
	plugin "github.com/ipfs/go-ipfs/plugin"
	"github.com/stretchr/testify/assert"
)

func TestPluginInterface(t *testing.T) {
	c := new(CIDTrack)

	err := c.Init(&plugin.Environment{})
	assert.NoError(t, err)

	err = c.Start(&core.IpfsNode{Exchange: &bitswap.Bitswap{}})
	assert.NoError(t, err)

	assert.Equal(t, "CIDTrack", c.Name())
	assert.NotEmpty(t, c.Version())

	err = c.Close()
	assert.NoError(t, err)
}

func TestHTTPEndpoints(t *testing.T) {
	c := new(CIDTrack)

	err := c.Init(&plugin.Environment{Config: map[string]interface{}{"listenAddress": "127.0.0.1:33221"}})
	assert.NoError(t, err)

	err = c.Start(&core.IpfsNode{Exchange: &bitswap.Bitswap{}})
	assert.NoError(t, err)

	// give some time for server to start up
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:33221/")
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get("http://127.0.0.1:33221/get")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get("http://127.0.0.1:33221/get/reset")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get("http://127.0.0.1:33221/reset")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()

	c.Close()
}
