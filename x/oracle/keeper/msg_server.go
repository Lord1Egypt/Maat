package keeper

import (
	"context"
	"fmt"

	"github.com/Lord1Egypt/Maat/x/oracle/types"
)

type msgServer struct {
	k *Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{k: k}
}

func (s *msgServer) SubmitPrice(ctx context.Context, msg *types.MsgSubmitPrice) (*types.MsgSubmitPriceResponse, error) {
	if !s.k.IsFeeder(msg.Feeder) {
		return nil, types.ErrUnauthorized.Wrapf("address %s is not a whitelisted feeder", msg.Feeder)
	}
	if s.k.IsHalted(msg.Denom) {
		return nil, types.ErrOracleHalted.Wrapf("denom %s is halted", msg.Denom)
	}
	s.k.getOrCreateOracle(msg.Denom)
	s.k.AppendVote(msg.Denom, msg.Price)
	return &types.MsgSubmitPriceResponse{}, nil
}

func (s *msgServer) ResumeOracle(ctx context.Context, msg *types.MsgResumeOracle) (*types.MsgResumeOracleResponse, error) {
	o := s.k.GetOracle(msg.Denom)
	if o == nil {
		return nil, types.ErrUnknownDenom.Wrapf("no oracle for denom %s", msg.Denom)
	}
	if !o.Halted {
		return nil, fmt.Errorf("oracle for denom %s is not halted", msg.Denom)
	}
	o.Resume()
	return &types.MsgResumeOracleResponse{}, nil
}
