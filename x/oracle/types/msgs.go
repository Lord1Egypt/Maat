package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgSubmitPrice{}
	_ sdk.Msg = &MsgResumeOracle{}
)

type MsgSubmitPrice struct {
	Feeder string `json:"feeder"`
	Denom  string `json:"denom"`
	Price  int64  `json:"price"`
}

func (m *MsgSubmitPrice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Feeder); err != nil {
		return ErrInvalidMsg.Wrapf("invalid feeder address: %s", err)
	}
	if m.Denom == "" {
		return ErrInvalidMsg.Wrap("denom cannot be empty")
	}
	if m.Price <= 0 {
		return ErrInvalidMsg.Wrapf("price must be positive: %d", m.Price)
	}
	return nil
}

func (m *MsgSubmitPrice) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Feeder)
	return []sdk.AccAddress{addr}
}

func (m *MsgSubmitPrice) Reset() {
	*m = MsgSubmitPrice{}
}

func (m *MsgSubmitPrice) String() string {
	return m.Feeder + ":" + m.Denom
}

func (m *MsgSubmitPrice) ProtoMessage() {}

type MsgResumeOracle struct {
	Authority string `json:"authority"`
	Denom     string `json:"denom"`
}

func (m *MsgResumeOracle) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidMsg.Wrapf("invalid authority address: %s", err)
	}
	if m.Denom == "" {
		return ErrInvalidMsg.Wrap("denom cannot be empty")
	}
	return nil
}

func (m *MsgResumeOracle) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgResumeOracle) Reset() {
	*m = MsgResumeOracle{}
}

func (m *MsgResumeOracle) String() string {
	return m.Authority + ":" + m.Denom
}

func (m *MsgResumeOracle) ProtoMessage() {}
