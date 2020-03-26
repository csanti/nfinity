package service

import (
	"testing"

	"go.dedis.ch/onet"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/pairing"
)

type networkSuite struct {
	kyber.Group
	pairing.Suite
}

func newNetworkSuite() *networkSuite {
	return &networkSuite{
		Group: Suite.G2(),
		Suite: Suite,
	}
}

func TestNfinity(t *testing.T) {
	log.Lvl1("Starting test")

	suite := newNetworkSuite()
	test := onet.NewTCPTest(suite)
	defer test.CloseAll()

	n := 6
	servers, roster, _ := test.GenTree(n, true)
	nfinities := make([]*Nfinity, n)
	for i := 0; i < n; i++ {
		c := &Config {
			Roster: roster,
			Index: i,
			N: n,
			Threshold: 4,
			GossipTime: 1000,
		}
		nfinities[i] = servers[i].Service(Name).(*Nfinity)
		nfinities[i].SetConfig(c)
	}
	
	done := make(chan bool)
	cb := func(r int) {
		if r > 10 {
			done <- true
		}
	}
	
	nfinities[0].AttachCallback(cb)
	go nfinities[0].Start()
	<-done
	log.Lvl1("finish")
	
}