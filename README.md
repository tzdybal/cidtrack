# cidtrack
Per CID bandwidth tracking for IPFS daemon.

## Building
Currently this plugin can be compiled "in source tree".
Required steps:
1. `cd /path/to/go-ipfs/plugin/plugins`
1. `git clone github.com:tzdybal/cidtrack.git`
1. `make -C cidtrack`

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
Plugin can be accessed via HTTP endpoint, listening on `listenAddress`.
### get
`{listenAddress}/get` returns all statistics gathered by CIDTracker.
Data is returned as CSV, with following columns:
| CID | bytes received | bytes sent |
|-----|----------------|------------|

`{listenAddress}/get/reset` returns all stats and then resets/clears all collected data.

### reset
`{listenAddress}/reset` resets/clears all collected data.
