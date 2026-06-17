package types

const (
	ModuleName = "market"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	ParamsKey      = []byte{0x01}
	ReservePrefix  = []byte{0x02}
	MarketHaltKey  = []byte{0x03}
)

func ReserveKey(denom string) []byte {
	return append(ReservePrefix, []byte(denom)...)
}
