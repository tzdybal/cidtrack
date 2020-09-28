package cidtrack

import (
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-ipfs/core"
	plugin "github.com/ipfs/go-ipfs/plugin"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("cidtrack")

// ensure that interfaces are implemented
var _ plugin.PluginDaemonInternal = (*CIDTrack)(nil)
var _ io.Closer = (*CIDTrack)(nil)

// CIDTrack is a Per CID bandwidth tracker
type CIDTrack struct {
	t      *tracker
	btswp  *bitswap.Bitswap
	config struct {
		listenAddress string
	}
}

// Name should return unique name of the plugin
func (c *CIDTrack) Name() string {
	return "CIDtrack"
}

// Version returns current version of the plugin
func (c *CIDTrack) Version() string {
	return "0.3.0"
}

// Init is called once when the Plugin is being loaded
// The plugin is passed an environment containing the path to the
// (possibly uninitialized) IPFS repo and the plugin's config.
func (c *CIDTrack) Init(env *plugin.Environment) error {
	log.Info("CIDtrack is being initialized!")
	if env.Config == nil { // defaults
		c.config.listenAddress = ":5002"
	} else {
		confMap := env.Config.(map[string]interface{})
		c.config.listenAddress = confMap["listenAddress"].(string)
	}
	log.Info("listenAddress:", c.config.listenAddress)
	return nil
}

// Starts starts a plugin
func (c *CIDTrack) Start(node *core.IpfsNode) error {
	log.Info("CIDtrack is starting!")
	btswp, ok := node.Exchange.(*bitswap.Bitswap)
	if !ok {
		return errors.New("couldn't cast node.Exchange as *bitswap.Bitswap")
	}
	c.t = newTracker()
	c.btswp = btswp

	bitswap.EnableWireTap(NewWireTap(c.t))(btswp)

	go c.t.run()
	log.Info("CIDTrack is running!")

	mux := http.NewServeMux()
	mux.HandleFunc("/get", c.handleGet)
	mux.HandleFunc("/reset", c.handleReset)

	go http.ListenAndServe(c.config.listenAddress, mux)

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
	log.Info("CIDTrack is stopping")
	bitswap.DisableWireTap()(c.btswp)
	return c.t.stop()
}

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&CIDTrack{},
}
