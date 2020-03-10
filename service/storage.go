package service

// roundStorage keeps tracks of all received valid blocks for a given round and
// if one has been notarized yet.
type RoundStorage struct {
	c       *Config                      // the config used to verify block signatures
	Round   int                          // the round number
	blocks  map[string]*blockStorage     // round blocks mapped from their hash
	tmpSigs map[int][]*SignatureProposal // all tmp signatures

	randomness int64
	// the finalizer
	finalizer *Finalizer
}

// newRoundStorage returns a new round storage for the given round
func NewRoundStorage(c *Config, round int, randomness int64, f *Finalizer) *roundStorage {
	return &roundStorage{
		c:          c,
		Round:      round,
		blocks:     make(map[string]*blockStorage),
		tmpSigs:    make(map[int][]*SignatureProposal),
		randomness: randomness,
		finalizer:  f,
	}
}
