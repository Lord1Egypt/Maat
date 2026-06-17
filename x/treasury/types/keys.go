package types

const (
	ModuleName = "treasury"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	ParamsKey   = []byte{0x01}
	BalancesKey = []byte{0x02}
)
