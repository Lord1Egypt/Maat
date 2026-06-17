package types

import "fmt"

type Params struct {
	BaseBps    int64
	VolMultBps int64
	SkewMaxBps int64
	MaxBps     int64
}

func DefaultParams() Params {
	return Params{
		BaseBps:    15,
		VolMultBps: 15,
		SkewMaxBps: 100,
		MaxBps:     300,
	}
}

func (p Params) Validate() error {
	if p.BaseBps <= 0 {
		return fmt.Errorf("base bps must be positive: %d", p.BaseBps)
	}
	if p.VolMultBps < 0 {
		return fmt.Errorf("vol mult bps must be non-negative: %d", p.VolMultBps)
	}
	if p.SkewMaxBps < 0 {
		return fmt.Errorf("skew max bps must be non-negative: %d", p.SkewMaxBps)
	}
	if p.MaxBps < p.BaseBps {
		return fmt.Errorf("max bps (%d) must be >= base bps (%d)", p.MaxBps, p.BaseBps)
	}
	return nil
}
