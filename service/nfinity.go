package service

import (
	"go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/onet"
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
	context *onet.Context
	c       *Config
	node    *Node
}

// NewDfinityService
func NewNfinityService(c *onet.Context) (onet.Service, error) {
	d := &Nfinity{
		context:          c,
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	c.RegisterProcessor(d, ConfigType)
	return d, nil
}

func (n *Nfinity) SetConfig(c *Config) {
	n.c = c
	n.node = NewNodeProcess(n.context, c, n.broadcast)
}

func (d *Dfinity) AttachCallback(fn func(int)) {
	chain := new(Chain)
	d.fin = NewFinalizer(d.c, chain, fn)
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
    n.Process(inner)
  case *BlockProposal
    n.Process(inner)
}

type BroadcastFn func(sis []*network.ServerIdentity, msg interface{})

func (n *Nfinity) broadcast(sis []*network.ServerIdentity, msg interface{}) {
	for _, si := range sis {
		if d.ServerIdentity().Equal(si) {
			continue
		}
		if err := d.ServiceProcessor.SendRaw(si, msg); err != nil {
			panic(err)
		}
	}
}
