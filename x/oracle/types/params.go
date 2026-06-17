package types

import "fmt"

type Params struct {
	AlphaNum   int64
	AlphaDen   int64
	MaxDevBps  int64
	BreakBps   int64
	VotePeriod int64
}

func DefaultParams() Params {
	return Params{
		AlphaNum:   8,
		AlphaDen:   10,
		MaxDevBps:  500,
		BreakBps:   5000,
		VotePeriod: 1,
	}
}

func (p Params) Validate() error {
	if p.AlphaDen <= 0 {
		return fmt.Errorf("alpha denominator must be positive: %d", p.AlphaDen)
	}
	if p.AlphaNum <= 0 || p.AlphaNum > p.AlphaDen {
		return fmt.Errorf("alpha numerator must be in (0, %d]: %d", p.AlphaDen, p.AlphaNum)
	}
	if p.MaxDevBps <= 0 {
		return fmt.Errorf("max deviation bps must be positive: %d", p.MaxDevBps)
	}
	if p.BreakBps <= 0 {
		return fmt.Errorf("break bps must be positive: %d", p.BreakBps)
	}
	if p.VotePeriod <= 0 {
		return fmt.Errorf("vote period must be positive: %d", p.VotePeriod)
	}
	return nil
}
