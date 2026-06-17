package treasury

import (
	"testing"
)

func TestPharaohBonds(t *testing.T) {
	bc := NewBondController()

	// Issue bond
	// Principal = $10,000, 8% coupon, duration = 10,000 blocks
	bond, err := bc.Issue("owner1", 10000, 800, 100, 10000)
	if err != nil {
		t.Fatalf("failed to issue bond: %v", err)
	}

	if bond.ID != 1 || bond.Owner != "owner1" || bond.PrincipalUSD != 10000 || bond.MatureHeight != 10100 {
		t.Fatalf("unexpected bond values: %+v", bond)
	}

	// Yield accrual
	// Principal = 10,000, coupon = 800 bps, elapsed = 1000 blocks, blocks/year = 100,000
	// yield = 10,000 * 800 * 1000 / (10000 * 100,000) = 8 USD
	yield, err := bc.AccrueYield(int64(bond.ID), 1100, 100000)
	if err != nil {
		t.Fatalf("failed to accrue yield: %v", err)
	}
	if yield != 8 {
		t.Fatalf("expected yield 8, got %d", yield)
	}

	// Redemption before maturity must fail
	_, err = bc.Redeem(int64(bond.ID), 5000)
	if err != ErrBondNotMatured {
		t.Fatalf("expected ErrBondNotMatured, got %v", err)
	}

	// Redemption at maturity
	principal, err := bc.Redeem(int64(bond.ID), 10100)
	if err != nil {
		t.Fatalf("failed to redeem bond: %v", err)
	}
	if principal != 10000 {
		t.Fatalf("expected principal 10000, got %d", principal)
	}
	if !bond.Redeemed {
		t.Fatal("expected bond to be redeemed")
	}

	// Double redemption must fail
	_, err = bc.Redeem(int64(bond.ID), 10100)
	if err != ErrAlreadyRedeemed {
		t.Fatalf("expected ErrAlreadyRedeemed, got %v", err)
	}
}
