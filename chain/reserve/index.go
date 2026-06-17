package reserve

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidWeights = errors.New("reserve: index weights must sum to 10000 bps")
	ErrRebalanceUnderway = errors.New("reserve: rebalance already underway")
)

type IndexWeight struct {
	Denom      string
	WeightBps  int64
}

type ReserveIndex struct {
	Weights       []IndexWeight
	LastFeeHeight int64
	FeeAccruedUSD int64
}

func NewReserveIndex(weights []IndexWeight, height int64) (*ReserveIndex, error) {
	if err := validateWeights(weights); err != nil {
		return nil, err
	}
	return &ReserveIndex{
		Weights:       weights,
		LastFeeHeight: height,
	}, nil
}

func (ri *ReserveIndex) AccrueManagementFee(currentHeight int64, totalNavUSD int64, feeBpsPerYear int64, blocksPerYear int64) (int64, error) {
	if currentHeight < ri.LastFeeHeight {
		return 0, fmt.Errorf("invalid height: current %d < last %d", currentHeight, ri.LastFeeHeight)
	}
	if blocksPerYear <= 0 {
		return 0, fmt.Errorf("invalid blocks per year: %d", blocksPerYear)
	}

	elapsed := currentHeight - ri.LastFeeHeight
	if elapsed == 0 {
		return 0, nil
	}

	// fee = NAV * (feeBps/10000) * (elapsed/blocksPerYear)
	feeNumerator := totalNavUSD * feeBpsPerYear * elapsed
	feeDenominator := int64(10000) * blocksPerYear

	fee := feeNumerator / feeDenominator
	ri.FeeAccruedUSD += fee
	ri.LastFeeHeight = currentHeight

	return fee, nil
}

func (ri *ReserveIndex) Rebalance(newWeights []IndexWeight) error {
	if err := validateWeights(newWeights); err != nil {
		return err
	}
	ri.Weights = newWeights
	return nil
}

func validateWeights(weights []IndexWeight) error {
	var sum int64
	for _, w := range weights {
		if w.WeightBps <= 0 {
			return ErrInvalidWeights
		}
		sum += w.WeightBps
	}
	if sum != 10000 {
		return ErrInvalidWeights
	}
	return nil
}
