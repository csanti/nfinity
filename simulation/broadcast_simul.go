package main

import (
	"github.com/BurntSushi/toml"
	"github.com/csanti/nfinity_cothority/protocol"
	"go.dedis.ch/onet"
)

func init() {
	onet.SimulationRegister("BroadcastProtocolSimul", NewBroadcastSimulation)
}

// SimulationProtocol implements onet.Simulation.
type BroadcastSimulation struct {
	protocol.Broadcast
}

// NewSimulationProtocol is used internally to register the simulation (see the init()
// function above).
func NewBroadcastSimulation(config string) (onet.Simulation, error) {
	bs := &BroadcastSimulation{}
	_, err := toml.Decode(config, bs)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// Setup implements onet.Simulation.
func (s *BroadcastSimulation) Setup(dir string, hosts []string) (*onet.SimulationConfig, error) {

	return nil, nil
}

// Node can be used to initialize each node before it will be run
// by the server. Here we call the 'Node'-method of the
// SimulationBFTree structure which will load the roster- and the
// tree-structure to speed up the first round.
func (s *BroadcastSimulation) Node(config *onet.SimulationConfig) error {

	return nil
}

// Run implements onet.Simulation.
func (s *BroadcastSimulation) Run(config *onet.SimulationConfig) error {

	return nil
}
