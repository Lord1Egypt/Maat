package types

import (
	"encoding/binary"
)

const (
	ModuleName = "bridge"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	ParamsKey        = []byte{0x01}
	LimiterPrefix    = []byte{0x02}
	WithdrawalPrefix = []byte{0x03}
)

func LimiterKey(denom string) []byte {
	return append(LimiterPrefix, []byte(denom)...)
}

func WithdrawalKey(denom string, id uint64) []byte {
	denomBytes := []byte(denom)
	key := make([]byte, len(WithdrawalPrefix)+len(denomBytes)+8)
	copy(key[0:], WithdrawalPrefix)
	copy(key[len(WithdrawalPrefix):], denomBytes)
	binary.BigEndian.PutUint64(key[len(WithdrawalPrefix)+len(denomBytes):], id)
	return key
}
