package cidtrack

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

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
	tap    *wireTap
	btswp  *bitswap.Bitswap
	srv    *http.Server
	config struct {
		listenAddress string
	}
}

// Name should return unique name of the plugin
func (c *CIDTrack) Name() string {
	return "CIDTrack"
}

// Version returns current version of the plugin
func (c *CIDTrack) Version() string {
	return "0.5.0"
}

// Init is called once when the Plugin is being loaded
// The plugin is passed an environment containing the path to the
// (possibly uninitialized) IPFS repo and the plugin's config.
func (c *CIDTrack) Init(env *plugin.Environment) error {
	log.Info("CIDtrack is being initialized")
	if env.Config == nil { // defaults
		c.config.listenAddress = ":5002"
	} else {
		confMap := env.Config.(map[string]interface{})
		c.config.listenAddress = confMap["listenAddress"].(string)
	}
	log.Info("listenAddress:", c.config.listenAddress)
	return nil
}

// Start starts a plugin
func (c *CIDTrack) Start(node *core.IpfsNode) error {
	log.Info("CIDtrack is starting")
	btswp, ok := node.Exchange.(*bitswap.Bitswap)
	if !ok {
		return errors.New("couldn't cast node.Exchange as *bitswap.Bitswap")
	}
	c.btswp = btswp
	c.tap = newWireTap()

	bitswap.EnableWireTap(c.tap)(btswp)

	go c.tap.run()
	log.Info("CIDTrack is running")

	mux := http.NewServeMux()
	mux.HandleFunc("/get", c.handleGet)
	mux.HandleFunc("/get/", c.handleGet)
	mux.HandleFunc("/reset", c.handleReset)

	c.srv = &http.Server{Addr: c.config.listenAddress, Handler: mux}

	go func() {
		log.Infof("Starting CIDTrack HTTP endpoint on: '%s'", c.config.listenAddress)
		if err := c.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("CIDTrack HTTP server failed", err)
		}
	}()

	return nil
}

// Close does all the necessary cleanup
func (c *CIDTrack) Close() error {
	log.Info("CIDTrack is stopping")
	bitswap.DisableWireTap()(c.btswp)
	c.tap.stop()
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	c.srv.Shutdown(ctx)
	return nil
}

func (c *CIDTrack) handleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	c.tap.mtx.RLock()
	json.NewEncoder(w).Encode(c.tap.stats)
	c.tap.mtx.RUnlock()

	if r.URL.Path == "/get/reset" {
		c.tap.reset()
	}
}

func (c *CIDTrack) handleReset(w http.ResponseWriter, r *http.Request) {
	c.tap.reset()
}

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&CIDTrack{},
}
