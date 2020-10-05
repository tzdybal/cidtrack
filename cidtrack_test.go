package cidtrack

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ipfs/go-bitswap"
	cid "github.com/ipfs/go-cid"
	u "github.com/ipfs/go-ipfs-util"
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

func TestGetAndReset(t *testing.T) {
	c := new(CIDTrack)

	err := c.Init(&plugin.Environment{Config: map[string]interface{}{"listenAddress": "127.0.0.1:33221"}})
	assert.NoError(t, err)

	err = c.Start(&core.IpfsNode{Exchange: &bitswap.Bitswap{}})
	assert.NoError(t, err)

	// give some time for server to start up
	time.Sleep(100 * time.Millisecond)
	cid1 := getCid("cid1")
	cid2 := getCid("cid2")

	c.tap.mtx.Lock()
	c.tap.stats[cid1] = 1
	c.tap.stats[cid2] = 2
	c.tap.mtx.Unlock()

	// first get - 2 records expected
	resp, err := http.Get("http://127.0.0.1:33221/get")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, err := parseMap(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, 1, data[cid1.String()])
	assert.Equal(t, 2, data[cid2.String()])

	// reset - 2 records expected
	resp, err = http.Get("http://127.0.0.1:33221/get/reset")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, err = parseMap(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
	assert.Equal(t, 1, data[cid1.String()])
	assert.Equal(t, 2, data[cid2.String()])

	// get after reset - no records expected
	resp, err = http.Get("http://127.0.0.1:33221/get")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	data, err = parseMap(resp.Body)
	assert.NoError(t, err)
	assert.Len(t, data, 0)

	c.Close()
}

func getCid(data string) cid.Cid {
	return cid.NewCidV0(u.Hash([]byte(data)))
}

func parseMap(r io.ReadCloser) (map[string]int, error) {
	body, err := ioutil.ReadAll(r)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	var data map[string]int
	err = json.Unmarshal(body, &data)
	return data, err
}
