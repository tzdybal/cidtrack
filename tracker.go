package cidtrack

import (
	cid "github.com/ipfs/go-cid"
)

type tracker struct {
	recv   chan stat
	sent   chan stat
	stopCh chan interface{}

	stats map[cid.Cid]*bwstat
}

type bwstat struct {
	recv uint64
	sent uint64
}

func newTracker() *tracker {
	return &tracker{
		recv:   make(chan stat),
		sent:   make(chan stat),
		stopCh: make(chan interface{}),
		stats:  make(map[cid.Cid]*bwstat),
	}
}

func (t *tracker) run() {
	for {
		select {
		case <-t.stopCh:
			break
		case s := <-t.recv:
			t.increment(s.cid, uint64(s.size), 0)
		case s := <-t.sent:
			t.increment(s.cid, 0, uint64(s.size))
		}
	}
}

func (t *tracker) increment(c cid.Cid, recv, sent uint64) {
	s, ok := t.stats[c]
	if ok {
		s.sent += sent
		s.recv += recv
	} else {
		t.stats[c] = &bwstat{recv: recv, sent: sent}
	}
}

func (t *tracker) stop() error {
	close(t.stopCh)
	return nil
}

func (t *tracker) reset() {
	t.stats = make(map[cid.Cid]*bwstat)
}

func (t *tracker) recvChan() chan stat {
	return t.recv
}

func (t *tracker) sentChan() chan stat {
	return t.sent
}
