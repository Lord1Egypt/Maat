package treasury

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPrincipal  = errors.New("bonds: principal must be positive")
	ErrInvalidMaturity   = errors.New("bonds: maturity block must be in the future")
	ErrBondNotMatured    = errors.New("bonds: bond has not reached maturity height")
	ErrAlreadyRedeemed   = errors.New("bonds: bond already redeemed")
)

type Bond struct {
	ID            uint64
	Owner         string
	PrincipalUSD  int64
	CouponBps     int64 // annual yield bps
	StartHeight   int64
	MatureHeight  int64
	YieldPaidUSD  int64
	LastPayHeight int64
	Redeemed      bool
}

type BondController struct {
	nextBondID uint64
	bonds      map[uint64]*Bond
}

func NewBondController() *BondController {
	return &BondController{
		nextBondID: 1,
		bonds:      make(map[uint64]*Bond),
	}
}

func (bc *BondController) Issue(owner string, principal int64, couponBps int64, startHeight int64, durationBlocks int64) (*Bond, error) {
	if principal <= 0 {
		return nil, ErrInvalidPrincipal
	}
	if durationBlocks <= 0 {
		return nil, ErrInvalidMaturity
	}

	bond := &Bond{
		ID:            bc.nextBondID,
		Owner:         owner,
		PrincipalUSD:  principal,
		CouponBps:     couponBps,
		StartHeight:   startHeight,
		MatureHeight:  startHeight + durationBlocks,
		LastPayHeight: startHeight,
	}

	bc.bonds[bond.ID] = bond
	bc.nextBondID++
	return bond, nil
}

func (bc *BondController) AccrueYield(bondID int64, currentHeight int64, blocksPerYear int64) (int64, error) {
	bond, ok := bc.bonds[uint64(bondID)]
	if !ok {
		return 0, fmt.Errorf("bond %d not found", bondID)
	}
	if bond.Redeemed {
		return 0, ErrAlreadyRedeemed
	}
	if currentHeight <= bond.LastPayHeight {
		return 0, nil
	}

	// Yield = Principal * (CouponBps/10000) * (elapsed/blocksPerYear)
	elapsed := currentHeight - bond.LastPayHeight
	// Limit payout to maturity height
	if bond.LastPayHeight >= bond.MatureHeight {
		return 0, nil
	}
	if currentHeight > bond.MatureHeight {
		elapsed = bond.MatureHeight - bond.LastPayHeight
	}

	yieldNum := bond.PrincipalUSD * bond.CouponBps * elapsed
	yieldDen := int64(10000) * blocksPerYear

	yield := yieldNum / yieldDen
	bond.YieldPaidUSD += yield
	bond.LastPayHeight = bond.LastPayHeight + elapsed // advanced up to currentHeight or MatureHeight

	return yield, nil
}

func (bc *BondController) Redeem(bondID int64, currentHeight int64) (int64, error) {
	bond, ok := bc.bonds[uint64(bondID)]
	if !ok {
		return 0, fmt.Errorf("bond %d not found", bondID)
	}
	if bond.Redeemed {
		return 0, ErrAlreadyRedeemed
	}
	if currentHeight < bond.MatureHeight {
		return 0, ErrBondNotMatured
	}

	bond.Redeemed = true
	return bond.PrincipalUSD, nil
}

func (bc *BondController) GetBond(id uint64) (*Bond, bool) {
	b, ok := bc.bonds[id]
	return b, ok
}
