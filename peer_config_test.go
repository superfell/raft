package raft

import (
	"testing"
)

func TestPeerConfigurations_Read(t *testing.T) {
	pc := peerConfigurations{}
	pc.initialize()
	if pc.readCh() == nil {
		t.Errorf("no inflight peer change, should get a real channel to read changes from")
	}
	pc.markInflight(40)
	if pc.readCh() != nil {
		t.Errorf("inflight peer change, readCh() should return nil as its not safe to start processing another peer change")
	}
	pc.onCommit(35)
	if pc.readCh() != nil {
		t.Errorf("inflight peer change, readCh() should return nil as its not safe to start processing another peer change")
	}
	pc.onCommit(40)
	if pc.readCh() == nil {
		t.Errorf("no inflight peer change, should get a real channel to read changes from")
	}
	pc.onCommit(44)
	if pc.readCh() == nil {
		t.Errorf("no inflight peer change, should get a real channel to read changes from")
	}
}