package service

import (
	"testing"

	"go.dedis.ch/onet"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/pairing"
	"go.dedis.ch/kyber/share"
	"go.dedis.ch/kyber/util/random"
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
	shares, public := dkg(3, 6)
	_, commits := public.Info()
	nfinities := make([]*Nfinity, n)
	for i := 0; i < n; i++ {
		c := &Config {
			Roster: roster,
			Index: i,
			N: n,
			Threshold: 4,
			GossipTime: 1000,
			Public: commits,
			Share: shares[i], // i have to check this..
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

func dkg(t, n int) ([]*share.PriShare, *share.PubPoly) {
	allShares := make([][]*share.PriShare, n)
	var public *share.PubPoly
	for i := 0; i < n; i++ {
		priPoly := share.NewPriPoly(G2, t, nil, random.New())
		allShares[i] = priPoly.Shares(n)
		if public == nil {
			public = priPoly.Commit(G2.Point().Base())
			continue
		}
		public, _ = public.Add(priPoly.Commit(G2.Point().Base()))
	}
	shares := make([]*share.PriShare, n)
	for i := 0; i < n; i++ {
		v := G2.Scalar().Zero()
		for j := 0; j < n; j++ {
			v = v.Add(v, allShares[j][i].V)
		}
		shares[i] = &share.PriShare{I: i, V: v}
	}
	return shares, public

}