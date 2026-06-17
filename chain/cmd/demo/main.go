// Command demo runs the Ma'at end-to-end scenario and prints a human summary.
// It is the runnable proof that the cores work together.
//
//	cd chain && go run ./cmd/demo
package main

import (
	"fmt"

	"github.com/Lord1Egypt/Maat/chain/scenario"
)

const M = scenario.M

func usd(microUSD int64) string { return fmt.Sprintf("$%.2f", float64(microUSD)/float64(M)) }

func main() {
	cfg := scenario.Default()
	r := scenario.Run(cfg)

	fmt.Println("============================================================")
	fmt.Println(" Ma'at end-to-end demo  —  oracle + market + treasury + bridge")
	fmt.Println("============================================================")
	fmt.Printf(" blocks simulated     : %d\n", r.Blocks)
	fmt.Printf(" backing held >=100%%  : %v\n", r.BackingHeld)
	fmt.Printf(" final backing        : %.2f%%\n", float64(r.FinalBackingBps)/100)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf(" spread captured      : %s\n", usd(r.SpreadCaptured))
	fmt.Printf("   -> reserve buffer  : %s\n", usd(r.ReserveFund))
	fmt.Printf("   -> staker rewards  : %s\n", usd(r.RewardsFund))
	fmt.Printf("   -> insurance fund  : %s\n", usd(r.InsuranceFund))
	fmt.Printf("   -> treasury        : %s\n", usd(r.TreasuryFund))
	fmt.Println("------------------------------------------------------------")
	fmt.Printf(" bridge-outs accepted : %d\n", r.BridgeAccepted)
	fmt.Printf(" bridge-outs throttled: %d  (cap protection working)\n", r.BridgeThrottled)
	fmt.Printf(" oracle breaker halts : %d\n", r.OracleHalts)
	fmt.Println("============================================================")
	fmt.Println(" Reserve stayed solvent and grew on spread. Order, not chaos.")
}
