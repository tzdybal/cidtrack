package cidtrack

import (
	cid "github.com/ipfs/go-cid"
)

type tracker struct {
	recv chan stat
	sent chan stat

	recvStats map[cid.Cid]uint64
	sentStats map[cid.Cid]uint64
}

func newTracker() *tracker {
	return &tracker{
		recv:      make(chan stat),
		sent:      make(chan stat),
		recvStats: make(map[cid.Cid]uint64),
		sentStats: make(map[cid.Cid]uint64),
	}
}

func (t *tracker) run() {
	for {
		// TODO(tzdybal): gently exit
		select {
		case s := <-t.recv:
			log.Debugf("recv cid=%s\tsize=%d\n", s.cid, s.size)
			t.recvStats[s.cid] += uint64(s.size)
		case s := <-t.sent:
			log.Debugf("sent cid=%s\tsize=%d\n", s.cid, s.size)
			t.sentStats[s.cid] += uint64(s.size)
		}
	}
}

func (t *tracker) recvChan() chan stat {
	return t.recv
}

func (t *tracker) sentChan() chan stat {
	return t.sent
}
