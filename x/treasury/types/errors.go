package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidParams = errors.Register(ModuleName, 2, "invalid params")
	ErrInvalidSplit  = errors.Register(ModuleName, 3, "invalid split configuration")
	ErrBadAmount     = errors.Register(ModuleName, 4, "amount must be non-negative")
)
