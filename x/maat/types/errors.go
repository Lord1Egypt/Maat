package types

import "cosmossdk.io/errors"

var (
	ErrInvalidInflation  = errors.Register(ModuleName, 2, "inflation exceeds maximum allowed")
	ErrInvalidParams     = errors.Register(ModuleName, 3, "invalid module parameters")
	ErrGenesisValidation = errors.Register(ModuleName, 4, "genesis validation failed")
	ErrMintFailed        = errors.Register(ModuleName, 5, "block reward minting failed")
	ErrSupplyMismatch    = errors.Register(ModuleName, 6, "allocation sum does not match total supply")
)
