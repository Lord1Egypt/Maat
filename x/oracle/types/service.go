package types

import "context"

type MsgSubmitPriceResponse struct{}

type MsgResumeOracleResponse struct{}

type MsgServer interface {
	SubmitPrice(context.Context, *MsgSubmitPrice) (*MsgSubmitPriceResponse, error)
	ResumeOracle(context.Context, *MsgResumeOracle) (*MsgResumeOracleResponse, error)
}
