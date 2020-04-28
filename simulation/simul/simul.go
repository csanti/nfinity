package main

import (
	// Service needs to be imported here to be instantiated.

	"github.com/csanti/onet/simul"
	_ "github.com/csanti/nfinity_cothority/simulation"
)

func main() {
	simul.Start()
}