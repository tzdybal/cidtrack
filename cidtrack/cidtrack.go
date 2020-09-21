package cidtrack

import (
	"errors"
	"fmt"

	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-ipfs/core"
	plugin "github.com/ipfs/go-ipfs/plugin"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("cidtrack")
var _ plugin.PluginDaemonInternal = (*CIDTrack)(nil)

// CIDTrack is a Per CID bandwidth tracker
type CIDTrack struct{}

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

	t := newTracker()
	bitswap.EnableWireTap(NewWireTap(t))(btswp)

	t.run()
	fmt.Println("CIDTrack is running!")

	return nil
}
