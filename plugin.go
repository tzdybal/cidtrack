package main

import (
	"github.com/ipfs/go-ipfs/plugin"

	"github.com/tzdybal/cidtrack/cidtrack"
)

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&cidtrack.CIDTrack{},
}
