package main

import (
	uniquepkgname "github.com/ipfs/go-ipfs/plugin/plugins/cidtrack"
)

var Plugins = uniquepkgname.Plugins //nolint

func main() {
	panic("this is a plugin, build it as a plugin, this is here as for go#20312")
}
