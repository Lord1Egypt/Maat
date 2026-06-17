package token

import "testing"

func TestStandardSumsToTotalSupply(t *testing.T) {
	if s := SumAllocations(Standard()); s != TotalSupply {
		t.Fatalf("allocations sum to %d, want %d", s, TotalSupply)
	}
}

func find(name string) Allocation {
	for _, a := range Standard() {
		if a.Name == name {
			return a
		}
	}
	panic("not found: " + name)
}

func TestCommunityTGEAndLinear(t *testing.T) {
	c := find("Community Sale") // 300M, 10% TGE, 12mo linear
	if u, _ := c.Unlocked(0); u != 30_000_000 {
		t.Fatalf("TGE unlock=%d, want 30M", u)
	}
	if u, _ := c.Unlocked(6); u != 30_000_000+270_000_000*6/12 {
		t.Fatalf("mid unlock=%d, want 165M", u)
	}
	if u, _ := c.Unlocked(12); u != 300_000_000 {
		t.Fatalf("final unlock=%d, want 300M", u)
	}
}

func TestTeamCliffThenLinear(t *testing.T) {
	tm := find("Team") // 150M, 0 TGE, 12mo cliff, 24mo linear
	if u, _ := tm.Unlocked(11); u != 0 {
		t.Fatalf("pre-cliff unlock=%d, want 0", u)
	}
	if u, _ := tm.Unlocked(12); u != 0 {
		t.Fatalf("at-cliff unlock=%d, want 0", u)
	}
	if u, _ := tm.Unlocked(24); u != 75_000_000 {
		t.Fatalf("mid-vest unlock=%d, want 75M", u)
	}
	if u, _ := tm.Unlocked(36); u != 150_000_000 {
		t.Fatalf("full unlock=%d, want 150M", u)
	}
}

func TestCirculatingMonotonicAndBounded(t *testing.T) {
	allocs := Standard()
	prev := int64(-1)
	for m := 0; m <= 72; m++ {
		c, err := CirculatingSupply(allocs, m)
		if err != nil {
			t.Fatalf("month %d: %v", m, err)
		}
		if c < prev {
			t.Fatalf("circulating decreased at month %d: %d < %d", m, c, prev)
		}
		if c > TotalSupply {
			t.Fatalf("circulating %d exceeds total at month %d", c, m)
		}
		prev = c
	}
	if final, _ := CirculatingSupply(allocs, 72); final != TotalSupply {
		t.Fatalf("not fully vested by month 72: %d", final)
	}
}

func TestInflationCapEnforced(t *testing.T) {
	if _, err := BlockReward(TotalSupply, 1001, 1); err != ErrBadInflation {
		t.Fatalf("err=%v, want ErrBadInflation", err)
	}
	if _, err := BlockReward(TotalSupply, 1000, 1); err != nil {
		t.Fatalf("10%% should be allowed: %v", err)
	}
}

func TestBlockReward(t *testing.T) {
	// 1B supply, 10% annual, ~15.7M blocks/yr (2s blocks) -> small positive reward
	r, err := BlockReward(TotalSupply, 1000, 15_768_000)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if r <= 0 {
		t.Fatalf("reward=%d, want positive", r)
	}
}
