package types

const (
	ModuleName = "oracle"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	ParamsKey  = []byte{0x01}
	MidPrefix  = []byte{0x02}
	VotePrefix = []byte{0x03}
	HaltPrefix = []byte{0x04}
)

func MidKey(denom string) []byte {
	return append(MidPrefix, []byte(denom)...)
}

func VoteKey(denom string) []byte {
	return append(VotePrefix, []byte(denom)...)
}

func HaltKey(denom string) []byte {
	return append(HaltPrefix, []byte(denom)...)
}
