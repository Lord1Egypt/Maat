package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSwap{}

type MsgSwap struct {
	Trader    string `json:"trader"`
	FromDenom string `json:"from_denom"`
	ToDenom   string `json:"to_denom"`
	Amount    int64  `json:"amount"`
}

func (m *MsgSwap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Trader); err != nil {
		return ErrInvalidMsg.Wrapf("invalid trader address: %s", err)
	}
	if m.FromDenom == "" {
		return ErrInvalidMsg.Wrap("from_denom cannot be empty")
	}
	if m.ToDenom == "" {
		return ErrInvalidMsg.Wrap("to_denom cannot be empty")
	}
	if m.FromDenom == m.ToDenom {
		return ErrSameDenom.Wrapf("%s", m.FromDenom)
	}
	if m.Amount <= 0 {
		return ErrInvalidMsg.Wrapf("amount must be positive: %d", m.Amount)
	}
	return nil
}

func (m *MsgSwap) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Trader)
	return []sdk.AccAddress{addr}
}

func (m *MsgSwap) Reset() {
	*m = MsgSwap{}
}

func (m *MsgSwap) String() string {
	return m.Trader + ":" + m.FromDenom + "->" + m.ToDenom
}

func (m *MsgSwap) ProtoMessage() {}
