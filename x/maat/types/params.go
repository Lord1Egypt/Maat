package types

import (
	"fmt"

	"github.com/Lord1Egypt/Maat/chain/token"
)

var (
	DefaultAnnualInflationBps int64 = 500
	DefaultBlocksPerYear      int64 = 15_768_000
)

type Params struct {
	AnnualInflationBps int64 `json:"annual_inflation_bps"`
	BlocksPerYear      int64 `json:"blocks_per_year"`
}

func DefaultParams() Params {
	return Params{
		AnnualInflationBps: DefaultAnnualInflationBps,
		BlocksPerYear:      DefaultBlocksPerYear,
	}
}

func (p Params) Validate() error {
	if p.AnnualInflationBps < 0 || p.AnnualInflationBps > int64(token.MaxInflationBps) {
		return fmt.Errorf("%w: got %d, max %d", ErrInvalidInflation, p.AnnualInflationBps, token.MaxInflationBps)
	}
	if p.BlocksPerYear <= 0 {
		return fmt.Errorf("%w: blocks_per_year must be positive", ErrInvalidParams)
	}
	return nil
}
