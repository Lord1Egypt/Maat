package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrBadAmount    = errors.Register(ModuleName, 2, "amount must be positive")
	ErrCapExceeded  = errors.Register(ModuleName, 3, "withdrawal cap exceeded")
	ErrNotFound     = errors.Register(ModuleName, 4, "pending withdrawal not found")
	ErrNotMatured   = errors.Register(ModuleName, 5, "delay window not elapsed")
	ErrNotPending   = errors.Register(ModuleName, 6, "withdrawal not in pending state")
	ErrInvalidMsg    = errors.Register(ModuleName, 7, "invalid message")
	ErrInvalidProof  = errors.Register(ModuleName, 8, "invalid proof")
	ErrUnknownDenom  = errors.Register(ModuleName, 9, "unknown reserve denom")
)
