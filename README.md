[![Go](https://github.com/tzdybal/cidtrack/workflows/Go/badge.svg)](https://github.com/tzdybal/cidtrack/actions/)
[![Go Report Card](https://goreportcard.com/badge/github.com/tzdybal/cidtrack)](https://goreportcard.com/report/github.com/tzdybal/cidtrack)

# cidtrack
Per CID bandwidth tracking for IPFS daemon.

CIDTrack is a ipfs daemon plugin that counts how many times a block was sent to the network.
Currently it's a low level utility that counts raw block usage, without taking care about subblocks (see https://docs.ipfs.io/how-to/work-with-blocks/). 

## Building
Building is done out of tree (of go-ipfs). Currently you have to checkout go-ipfs and manually update go-bitswap dependency, because change adding `WireTap` is merged into master but not yet released.

Daemon:
1. `git clone https://github.com/ipfs/go-ipfs.git`
1. `cd go-ipfs`
1. `go get -v github.com/ipfs/go-bitswap@bc3df6b`
1. `go build`

Plugin:
1. `git clone https://github.com/tzdybal/cidtrack.git`
1. `cd cidtrack`
1. `make go.mod IPFS_VERSION=/ABSOLUTE/PATH/TO/go-ipfs`
1. `make install`

## Configuration
By default CIDTracker listen on ":5002" (port 5002 on all addresses).
This can be changed using `listenAddress` configuration, for example:
```
  "Plugins": {
    "Plugins": {
      "CIDtrack": {
        "Config": {
          "listenAddress": "127.0.0.1:8888"
        }
      }
    }
  }
```

## Usage
Plugin can be accessed via HTTP endpoint, listening on `listenAddress`, default address (`http://127.0.0.1:5002`) will be used in examples.
### get
`http://127.0.0.1:5002/get` returns all statistics gathered by CIDTracker.
Data is returned as JSON as associative array (object), where key is a CID of block, and valuue is a number of times given block was sent.
For example:
```
curl localhost:5002/get
{"Qma5RSy8wpWUpnXfegzNz6iJnLSwWUQdrpmxEar3sZT5GX":1,"QmapJQeFtp3rtZs3N1nKPxKgcRhGkRjxjNjYkafjoQXJNf":1}
```

`http://127.0.0.1:5002/get/reset` returns all stats and then resets/clears all collected data.

### reset
`http://127.0.0.1:5002/reset` resets/clears all collected data (without returning any value).
