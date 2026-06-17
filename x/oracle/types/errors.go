package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrNoSources    = errors.Register(ModuleName, 2, "no price sources")
	ErrAllFiltered  = errors.Register(ModuleName, 3, "all sources filtered as outliers")
	ErrOracleHalted = errors.Register(ModuleName, 4, "oracle circuit breaker halted")
	ErrUnknownDenom = errors.Register(ModuleName, 5, "unknown denom")
	ErrUnauthorized = errors.Register(ModuleName, 6, "unauthorized feeder")
	ErrInvalidMsg   = errors.Register(ModuleName, 7, "invalid message")
)
