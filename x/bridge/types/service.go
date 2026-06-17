package types

import (
	"context"
)

type MsgBridgeInResponse struct{}

type MsgBridgeOutResponse struct {
	WithdrawalId uint64 `json:"withdrawal_id"`
	Status       int32  `json:"status"`
}

type MsgExecuteWithdrawalResponse struct{}

type MsgCancelWithdrawalResponse struct{}

type MsgServer interface {
	BridgeIn(context.Context, *MsgBridgeIn) (*MsgBridgeInResponse, error)
	BridgeOut(context.Context, *MsgBridgeOut) (*MsgBridgeOutResponse, error)
	ExecuteWithdrawal(context.Context, *MsgExecuteWithdrawal) (*MsgExecuteWithdrawalResponse, error)
	CancelWithdrawal(context.Context, *MsgCancelWithdrawal) (*MsgCancelWithdrawalResponse, error)
}
