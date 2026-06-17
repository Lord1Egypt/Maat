package types

import "context"

type MsgSwapResponse struct {
	Executed int64 `json:"executed"`
	SpreadUSD int64 `json:"spread_usd"`
}

type MsgServer interface {
	Swap(context.Context, *MsgSwap) (*MsgSwapResponse, error)
}
