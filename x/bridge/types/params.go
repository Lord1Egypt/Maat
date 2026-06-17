package types

import (
	"fmt"
)

type Params struct {
	DailyCapUnits int64
	WindowBlocks  int64
	LargeTxUnits  int64
	DelayBlocks   int64
}

func DefaultParams() Params {
	return Params{
		DailyCapUnits: 50000000,
		WindowBlocks:  43200,
		LargeTxUnits:  10000000,
		DelayBlocks:   21600,
	}
}

func (p Params) Validate() error {
	if p.DailyCapUnits <= 0 {
		return fmt.Errorf("daily cap units must be positive: %d", p.DailyCapUnits)
	}
	if p.WindowBlocks <= 0 {
		return fmt.Errorf("window blocks must be positive: %d", p.WindowBlocks)
	}
	if p.LargeTxUnits <= 0 {
		return fmt.Errorf("large tx units must be positive: %d", p.LargeTxUnits)
	}
	if p.DelayBlocks <= 0 {
		return fmt.Errorf("delay blocks must be positive: %d", p.DelayBlocks)
	}
	return nil
}
