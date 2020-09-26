package cidtrack

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-ipfs/core"
	plugin "github.com/ipfs/go-ipfs/plugin"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("cidtrack")
var _ plugin.PluginDaemonInternal = (*CIDTrack)(nil)
var _ io.Closer = (*CIDTrack)(nil)

// CIDTrack is a Per CID bandwidth tracker
type CIDTrack struct {
	t *tracker
}

// Name should return unique name of the plugin
func (c *CIDTrack) Name() string {
	return "cidtrack"
}

// Version returns current version of the plugin
func (c *CIDTrack) Version() string {
	return "0.1.0"
}

// Init is called once when the Plugin is being loaded
// The plugin is passed an environment containing the path to the
// (possibly uninitialized) IPFS repo and the plugin's config.
func (c *CIDTrack) Init(env *plugin.Environment) error {
	log.Info("cidtrack is being initialized!")
	return nil
}

// Starts starts a plugin
func (c *CIDTrack) Start(node *core.IpfsNode) error {
	log.Info("cidtrack is starting!")
	btswp, ok := node.Exchange.(*bitswap.Bitswap)
	if !ok {
		return errors.New("couldn't cast node.Exchange as *bitswap.Bitswap")
	}

	c.t = newTracker()
	bitswap.EnableWireTap(NewWireTap(c.t))(btswp)

	go c.t.run()
	log.Info("CIDTrack is running!")

	mux := http.NewServeMux()
	mux.HandleFunc("/get", c.handleGet)
	mux.HandleFunc("/reset", c.handleReset)

	go http.ListenAndServe(":5002", mux)

	return nil
}

func (c *CIDTrack) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	writer := csv.NewWriter(w)
	for cid, s := range c.t.stats {
		writer.Write([]string{cid.String(), strconv.FormatUint(s.recv, 10), strconv.FormatUint(s.sent, 10)})
	}
	writer.Flush()

	if r.URL.Path == "/get/reset" {
		c.t.reset()
	}
}

func (c *CIDTrack) handleReset(w http.ResponseWriter, r *http.Request) {
	c.t.reset()
}

func (c *CIDTrack) Close() error {
	fmt.Println("CIDTrack is stopping")
	return c.t.stop()
}

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&CIDTrack{},
}
