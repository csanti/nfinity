package service

import (
	"sync"
	//"time"

	"go.dedis.ch/onet"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/onet/network"
)

type Node struct {
	*sync.Cond
	*onet.ServiceProcessor

	// config
	c *Config
	// current round number
	round int
	// finalized Chain
	chain *BlockChain
	// information of previous rounds
	//rounds map[int]*RoundStorage
	// done callback
	done func(int) // callsback number of finalized blocks
	// store signatures received for current rounds
	tmpSigs map[int]*PartialSignature
	// store block proposals received for current rounds
	tmpBlocks map[int]*BlockProposal
	// node started the Consensus
	isGenesis bool

	// for firsts tests
	receivedBlockProposals map[int]*BlockProposal

	b BroadcastFn
}

func NewNodeProcess(c *onet.Context, conf *Config, b BroadcastFn) *Node {
	// need to create chain first
	chain := new(BlockChain)
	n := &Node {
		ServiceProcessor: onet.NewServiceProcessor(c),
		chain: chain,
		c:     conf,
    	receivedBlockProposals: make(map[int]*BlockProposal),
    	b: b,
	}
	return n
}

func (n *Node) StartConsensus() {

	log.Lvl1("Sending bootstrap message...")
	n.isGenesis = true
	/*
	packet := &Bootstrap{
		Block: GenesisBlock,
		Seed:  1234,
	}
	*/
	// send bootstrap to all nodes
	log.Lvl1("Consensus started")
}

func (n *Node) Process(e *network.Envelope) {
	n.Cond.L.Lock()
	defer n.Cond.L.Unlock()
	defer n.Cond.Broadcast()
	switch inner := e.Msg.(type) {
	case *BlockProposal:
		n.NewBlockProposal(inner)
	}
}

func (n *Node) NewRound(round int) {
	// TODO: cleanup

	go n.roundLoop(round)
}

func (n *Node) NewBlockProposal(p *BlockProposal) {
	log.Lvl3("Processing new block proposal")
	//if p.Block.Round < n.round {
	//	log.Lvl2("received too old block")
	//	return
	//}
	// TODO: check if its in storage

	// TODO: when implementing gossip, we have to check if already received the signature

	// n.tmpBlocks[p.Round] = append(n.tmpBlocks[p.Round], p)
	// if p.Signatures[0] = nil {
	//   // block proposal does not contain any signature
	//   log.Lvl2("received block proposal without signatures")
	//   return
	// }
	// n.tmpSigs[p.Round] = append(n.tmpSigs[p.Round], p.Signatures[0])
	//n.receivedBlockProposals[p.Round]++
}

func (n *Node) roundLoop(round int) {
	log.Lvl3("Round Loop")

	// wait block time to receive messages
	//time.Sleep(n.c.BlockTime * time.Millisecond)

	// wait on new inputs
	n.Cond.Wait()

}
