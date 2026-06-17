package keeper

import (
	"context"
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Lord1Egypt/Maat/chain/bridge"
	"github.com/Lord1Egypt/Maat/x/bridge/types"
)

type msgServer struct {
	k *Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{k: k}
}

func (s *msgServer) BridgeIn(ctx context.Context, msg *types.MsgBridgeIn) (*types.MsgBridgeInResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	vaaObj, err := bridge.UnmarshalVAA(msg.Proof)
	if err != nil {
		return nil, types.ErrInvalidProof.Wrapf("failed to unmarshal VAA: %s", err)
	}
	transfer, err := bridge.ParseTokenBridgeTransfer(vaaObj.Payload)
	if err != nil {
		return nil, types.ErrInvalidProof.Wrapf("failed to parse token bridge transfer: %s", err)
	}

	parsedAmt := int64(binary.BigEndian.Uint64(transfer.Amount[24:32]))
	if parsedAmt != msg.Amount {
		return nil, types.ErrInvalidProof.Wrapf("VAA amount %d does not match message amount %d", parsedAmt, msg.Amount)
	}

	r := s.k.marketKeeper.GetReserve(msg.Denom)
	if r == nil {
		return nil, types.ErrUnknownDenom.Wrapf("no reserve for %s", msg.Denom)
	}

	// 1. Core logic
	if err := r.BridgeIn(msg.Amount); err != nil {
		return nil, fmt.Errorf("core bridge in failed: %w", err)
	}
	s.k.marketKeeper.SetReserve(msg.Denom, r)

	// 2. Token minting & transfer
	depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid depositor: %s", err)
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(msg.Denom, msg.Amount))
	if err := s.k.bankKeeper.MintCoins(sdkCtx, types.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("failed to mint coins: %w", err)
	}

	if err := s.k.bankKeeper.SendCoinsFromModuleToAccount(sdkCtx, types.ModuleName, depositor, coins); err != nil {
		return nil, fmt.Errorf("failed to send coins: %w", err)
	}

	return &types.MsgBridgeInResponse{}, nil
}

func (s *msgServer) BridgeOut(ctx context.Context, msg *types.MsgBridgeOut) (*types.MsgBridgeOutResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	withdrawer, err := sdk.AccAddressFromBech32(msg.Withdrawer)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid withdrawer: %s", err)
	}

	lim := s.k.GetLimiter(msg.Denom)
	w, err := lim.RequestOut(msg.Amount, sdkCtx.BlockHeight())
	if err != nil {
		return nil, wrapCoreError(err)
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(msg.Denom, msg.Amount))
	if err := s.k.bankKeeper.SendCoinsFromAccountToModule(sdkCtx, withdrawer, types.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("failed to lock coins: %w", err)
	}

	s.k.SetWithdrawer(w.ID, msg.Withdrawer)

	if w.Status == bridge.StatusImmediate {
		r := s.k.marketKeeper.GetReserve(msg.Denom)
		if r == nil {
			return nil, types.ErrUnknownDenom.Wrapf("no reserve for %s", msg.Denom)
		}
		if err := r.BridgeOut(msg.Amount); err != nil {
			return nil, fmt.Errorf("core bridge out failed: %w", err)
		}
		s.k.marketKeeper.SetReserve(msg.Denom, r)

		if err := s.k.bankKeeper.BurnCoins(sdkCtx, types.ModuleName, coins); err != nil {
			return nil, fmt.Errorf("failed to burn coins: %w", err)
		}
		w.Status = bridge.StatusExecuted
	}

	return &types.MsgBridgeOutResponse{
		WithdrawalId: w.ID,
		Status:       int32(w.Status),
	}, nil
}

func (s *msgServer) ExecuteWithdrawal(ctx context.Context, msg *types.MsgExecuteWithdrawal) (*types.MsgExecuteWithdrawalResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var foundLimiter *bridge.Limiter
	var foundDenom string
	var foundW *bridge.Withdrawal
	for _, denom := range s.k.AllDenoms() {
		lim := s.k.GetLimiter(denom)
		if w, ok := lim.Get(int64(msg.WithdrawalID)); ok {
			foundLimiter = lim
			foundDenom = denom
			foundW = w
			break
		}
	}
	if foundLimiter == nil {
		return nil, types.ErrNotFound.Wrapf("withdrawal %d", msg.WithdrawalID)
	}

	if err := foundLimiter.Execute(int64(msg.WithdrawalID), sdkCtx.BlockHeight()); err != nil {
		return nil, wrapCoreError(err)
	}

	r := s.k.marketKeeper.GetReserve(foundDenom)
	if r == nil {
		return nil, types.ErrUnknownDenom.Wrapf("no reserve for %s", foundDenom)
	}
	if err := r.BridgeOut(foundW.Amount); err != nil {
		return nil, fmt.Errorf("core bridge out failed: %w", err)
	}
	s.k.marketKeeper.SetReserve(foundDenom, r)

	coins := sdk.NewCoins(sdk.NewInt64Coin(foundDenom, foundW.Amount))
	if err := s.k.bankKeeper.BurnCoins(sdkCtx, types.ModuleName, coins); err != nil {
		return nil, fmt.Errorf("failed to burn coins: %w", err)
	}

	return &types.MsgExecuteWithdrawalResponse{}, nil
}

func (s *msgServer) CancelWithdrawal(ctx context.Context, msg *types.MsgCancelWithdrawal) (*types.MsgCancelWithdrawalResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if msg.Authority != s.k.authority {
		return nil, sdkerrors.ErrUnauthorized.Wrapf("expected authority %s, got %s", s.k.authority, msg.Authority)
	}

	var foundLimiter *bridge.Limiter
	var foundDenom string
	var foundW *bridge.Withdrawal
	for _, denom := range s.k.AllDenoms() {
		lim := s.k.GetLimiter(denom)
		if w, ok := lim.Get(int64(msg.WithdrawalID)); ok {
			foundLimiter = lim
			foundDenom = denom
			foundW = w
			break
		}
	}
	if foundLimiter == nil {
		return nil, types.ErrNotFound.Wrapf("withdrawal %d", msg.WithdrawalID)
	}

	if err := foundLimiter.Cancel(int64(msg.WithdrawalID)); err != nil {
		return nil, wrapCoreError(err)
	}

	withdrawerStr, ok := s.k.GetWithdrawer(msg.WithdrawalID)
	if !ok {
		return nil, types.ErrNotFound.Wrapf("withdrawer for withdrawal %d not found", msg.WithdrawalID)
	}
	withdrawer, err := sdk.AccAddressFromBech32(withdrawerStr)
	if err != nil {
		return nil, fmt.Errorf("invalid withdrawer address registered: %w", err)
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(foundDenom, foundW.Amount))
	if err := s.k.bankKeeper.SendCoinsFromModuleToAccount(sdkCtx, types.ModuleName, withdrawer, coins); err != nil {
		return nil, fmt.Errorf("failed to refund coins: %w", err)
	}

	return &types.MsgCancelWithdrawalResponse{}, nil
}

func wrapCoreError(err error) error {
	switch err {
	case bridge.ErrBadAmount:
		return types.ErrBadAmount.Wrap(err.Error())
	case bridge.ErrCapExceeded:
		return types.ErrCapExceeded.Wrap(err.Error())
	case bridge.ErrNotFound:
		return types.ErrNotFound.Wrap(err.Error())
	case bridge.ErrNotMatured:
		return types.ErrNotMatured.Wrap(err.Error())
	case bridge.ErrNotPending:
		return types.ErrNotPending.Wrap(err.Error())
	default:
		return err
	}
}
