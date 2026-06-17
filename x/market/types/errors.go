package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInsufficientBacking = errors.Register(ModuleName, 2, "trade would break 100% backing")
	ErrMarketHalted        = errors.Register(ModuleName, 3, "market halted due to insufficient backing")
	ErrBadAmount           = errors.Register(ModuleName, 4, "amount must be positive")
	ErrOracleHalted        = errors.Register(ModuleName, 5, "oracle is halted for denom")
	ErrNoOraclePrice       = errors.Register(ModuleName, 6, "no oracle price available")
	ErrUnknownDenom        = errors.Register(ModuleName, 7, "unknown reserve denom")
	ErrInvalidMsg          = errors.Register(ModuleName, 8, "invalid message")
	ErrSameDenom           = errors.Register(ModuleName, 9, "from and to denom must differ")
)
