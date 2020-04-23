package service

import (
	"go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/onet"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/onet/network"
)

var Suite = bn256.NewSuite()
var G2 = Suite.G2()
var Name = "nfinity"

func init() {
	onet.RegisterNewService(Name, NewNfinityService)
}

// Dfinity service is either a beacon a notarizer or a block maker
type Nfinity struct {
	*onet.ServiceProcessor
	context 	*onet.Context
	c       	*Config
	node    	*Node
	blockChain 	*BlockChain
}

// NewDfinityService
func NewNfinityService(c *onet.Context) (onet.Service, error) {
	n := &Nfinity{
		context:          c,
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	c.RegisterProcessor(n, ConfigType)
	c.RegisterProcessor(n, BootstrapType)
	c.RegisterProcessor(n, BlockProposalType)
	c.RegisterProcessor(n, NotarizedBlockType)
	return n, nil
}

func (n *Nfinity) SetConfig(c *Config) {
	n.c = c
	if n.c.CommunicationMode == 0 {
		n.node = NewNodeProcess(n.context, c, n.broadcast, n.gossip)
	} else if n.c.CommunicationMode == 1 {
		n.node = NewNodeProcess(n.context, c, n.broadcast, n.gossip)	
	} else {
		panic("Invalid communication mode")
	}	
}

func (n *Nfinity) AttachCallback(fn func(int)) {
	// attach to something.. haha lol xd
	if n.node != nil {
		n.node.AttachCallback(fn)	
	} else {
		log.Lvl1("Could not attach callback, node is nil")
	}
}


func (n *Nfinity) Start() {
	// send a bootstrap message
	if n.node != nil {
		n.node.StartConsensus()
	} else {
		panic("that should not happen")
	}
}

// Process
func (n *Nfinity) Process(e *network.Envelope) {
	switch inner := e.Msg.(type) {
	case *Config:
		n.SetConfig(inner)
	case *Bootstrap:
		n.node.Process(e)
	case *BlockProposal:
		n.node.Process(e)
	case *NotarizedBlock:
		n.node.Process(e)
	default:
		log.Lvl1("Received unidentified message")
	}
}

// depreciated
func (n *Nfinity) getRandomPeers(numPeers int) []*network.ServerIdentity {
	var results []*network.ServerIdentity
	for i := 0; i < numPeers; {
		posPeer := n.c.Roster.RandomServerIdentity()
		if n.ServerIdentity().Equal(posPeer) {
			// selected itself
			continue
		}
		results = append(results, posPeer)
		i++
	}
	return results
}

type BroadcastFn func(sis []*network.ServerIdentity, msg interface{})

func (n *Nfinity) broadcast(sis []*network.ServerIdentity, msg interface{}) {
	for _, si := range sis {
		if n.ServerIdentity().Equal(si) {
			continue
		}
		log.Lvlf4("Broadcasting from: %s to: %s",n.ServerIdentity(), si)
		if err := n.ServiceProcessor.SendRaw(si, msg); err != nil {
			log.Lvl1("Error sending message")
			//panic(err)
		}
	}
}

func (n *Nfinity) gossip(sis []*network.ServerIdentity, msg interface{}) {
	//targets := n.getRandomPeers(n.c.GossipPeers)
	targets := n.c.Roster.RandomSubset(n.ServerIdentity(), n.c.GossipPeers).List
	for k, target := range targets {
		if k == 0 {continue}
		log.Lvlf4("Gossiping from: %s to: %s",n.ServerIdentity(), target)
		if err := n.ServiceProcessor.SendRaw(target, msg); err != nil {
			log.Lvl1("Error sending message")
		}
	}
	
}
