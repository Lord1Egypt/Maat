package types

import "errors"

var (
	ErrBackingBelowMinimum = errors.New("reserve: backing below minimum threshold")
	ErrDenomNotFound       = errors.New("reserve: denom not registered")
)
