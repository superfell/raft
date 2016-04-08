package raft

// peerConfigurations keeps track of any inflight peer change, and while
// its inflight stops any other peer change from being processed.
//
// it is assumed that this is only interacted with on the primary goroutine
//
type peerConfigurations struct {
	changeCh      chan *logFuture
	inflightIndex uint64
}

func (pc *peerConfigurations) initialize() {
	pc.changeCh = make(chan *logFuture)
}

// readCh returns a channel that peer change requests can be read from.
// if its not currently safe to work on a peer change [because a previous
// one is still in progress], this returns nil [and so you'd better be
// reading this channel with a select statement]
func (pc *peerConfigurations) readCh() <-chan *logFuture {
	if pc.inflightIndex > 0 {
		return nil
	}
	return pc.changeCh
}

// keep track of the inflight peer config change that was started, while
// this is inflight, we can't start processing any more config changes
func (pc *peerConfigurations) markInflight(logIndex uint64) {
	pc.inflightIndex = logIndex
}

// this should be called when the commit index changes, we notice if our
// inflight config change is committed and if so mark it as no longer inflight
func (pc *peerConfigurations) onCommit(commitIndex uint64) {
	if pc.inflightIndex > 0 && commitIndex >= pc.inflightIndex {
		pc.inflightIndex = 0
	}
}
