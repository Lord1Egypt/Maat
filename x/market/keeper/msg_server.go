package keeper

import (
	"context"
	"fmt"

	"github.com/Lord1Egypt/Maat/chain/pricing"
	"github.com/Lord1Egypt/Maat/chain/treasury"
	"github.com/Lord1Egypt/Maat/x/market/types"
)

const nativeDenom = "umaat"

type msgServer struct {
	k *Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{k: k}
}

func (s *msgServer) Swap(ctx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	if s.k.IsMarketHalted() {
		return nil, types.ErrMarketHalted.Wrap("global market halt active")
	}

	if s.k.oracleKeeper.IsHalted(msg.FromDenom) {
		return nil, types.ErrOracleHalted.Wrapf("denom %s", msg.FromDenom)
	}
	if s.k.oracleKeeper.IsHalted(msg.ToDenom) {
		return nil, types.ErrOracleHalted.Wrapf("denom %s", msg.ToDenom)
	}

	var spreadCaptured int64
	// Sell side: user sells FromDenom to protocol
	if msg.FromDenom != nativeDenom {
		r := s.k.GetReserve(msg.FromDenom)
		if r == nil {
			return nil, types.ErrUnknownDenom.Wrapf("no reserve for %s", msg.FromDenom)
		}
		q, err := s.k.GetQuote(msg.FromDenom, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get quote for %s: %w", msg.FromDenom, err)
		}
		if err := r.SellWrapped(msg.Amount, q); err != nil {
			return nil, wrapCoreError(err)
		}
		spreadCaptured += (q.Mid - q.Buy) * msg.Amount / pricing.MicroUSD
	}

	// Buy side: user buys ToDenom from protocol
	if msg.ToDenom != nativeDenom {
		r := s.k.GetReserve(msg.ToDenom)
		if r == nil {
			return nil, types.ErrUnknownDenom.Wrapf("no reserve for %s", msg.ToDenom)
		}
		q, err := s.k.GetQuote(msg.ToDenom, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to get quote for %s: %w", msg.ToDenom, err)
		}
		if err := r.BuyWrapped(msg.Amount, q); err != nil {
			return nil, wrapCoreError(err)
		}
		spreadCaptured += (q.Sell - q.Mid) * msg.Amount / pricing.MicroUSD
	}

	if spreadCaptured > 0 && s.k.treasuryKeeper != nil {
		if _, err := s.k.treasuryKeeper.Collect(spreadCaptured, treasury.DefaultSpreadSplit); err != nil {
			return nil, fmt.Errorf("treasury collect failed: %w", err)
		}
	}

	return &types.MsgSwapResponse{
		Executed:  msg.Amount,
		SpreadUSD: spreadCaptured,
	}, nil
}

func wrapCoreError(err error) error {
	switch {
	case err == pricing.ErrInsufficientBacking:
		return types.ErrInsufficientBacking.Wrap(err.Error())
	case err == pricing.ErrMarketHalted:
		return types.ErrMarketHalted.Wrap(err.Error())
	case err == pricing.ErrBadAmount:
		return types.ErrBadAmount.Wrap(err.Error())
	default:
		return err
	}
}
