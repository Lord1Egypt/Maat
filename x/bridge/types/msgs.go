package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgBridgeIn{}
	_ sdk.Msg = &MsgBridgeOut{}
	_ sdk.Msg = &MsgExecuteWithdrawal{}
	_ sdk.Msg = &MsgCancelWithdrawal{}
)

type MsgBridgeIn struct {
	Depositor string `json:"depositor"`
	Denom     string `json:"denom"`
	Amount    int64  `json:"amount"`
	Proof     []byte `json:"proof"`
}

func (m *MsgBridgeIn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Depositor); err != nil {
		return ErrInvalidMsg.Wrapf("invalid depositor address: %s", err)
	}
	if m.Denom == "" {
		return ErrInvalidMsg.Wrap("denom cannot be empty")
	}
	if m.Amount <= 0 {
		return ErrInvalidMsg.Wrapf("amount must be positive: %d", m.Amount)
	}
	if len(m.Proof) == 0 {
		return ErrInvalidMsg.Wrap("proof cannot be empty")
	}
	return nil
}

func (m *MsgBridgeIn) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Depositor)
	return []sdk.AccAddress{addr}
}

func (m *MsgBridgeIn) Reset() {
	*m = MsgBridgeIn{}
}

func (m *MsgBridgeIn) String() string {
	return m.Depositor + ":" + m.Denom
}

func (m *MsgBridgeIn) ProtoMessage() {}

type MsgBridgeOut struct {
	Withdrawer string `json:"withdrawer"`
	Denom      string `json:"denom"`
	ToAddress  string `json:"to_address"`
	Amount     int64  `json:"amount"`
}

func (m *MsgBridgeOut) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Withdrawer); err != nil {
		return ErrInvalidMsg.Wrapf("invalid withdrawer address: %s", err)
	}
	if m.Denom == "" {
		return ErrInvalidMsg.Wrap("denom cannot be empty")
	}
	if m.ToAddress == "" {
		return ErrInvalidMsg.Wrap("to_address cannot be empty")
	}
	if m.Amount <= 0 {
		return ErrInvalidMsg.Wrapf("amount must be positive: %d", m.Amount)
	}
	return nil
}

func (m *MsgBridgeOut) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Withdrawer)
	return []sdk.AccAddress{addr}
}

func (m *MsgBridgeOut) Reset() {
	*m = MsgBridgeOut{}
}

func (m *MsgBridgeOut) String() string {
	return m.Withdrawer + ":" + m.Denom + "->" + m.ToAddress
}

func (m *MsgBridgeOut) ProtoMessage() {}

type MsgExecuteWithdrawal struct {
	Executor     string `json:"executor"`
	WithdrawalID uint64 `json:"withdrawal_id"`
}

func (m *MsgExecuteWithdrawal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Executor); err != nil {
		return ErrInvalidMsg.Wrapf("invalid executor address: %s", err)
	}
	if m.WithdrawalID == 0 {
		return ErrInvalidMsg.Wrap("withdrawal id cannot be zero")
	}
	return nil
}

func (m *MsgExecuteWithdrawal) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Executor)
	return []sdk.AccAddress{addr}
}

func (m *MsgExecuteWithdrawal) Reset() {
	*m = MsgExecuteWithdrawal{}
}

func (m *MsgExecuteWithdrawal) String() string {
	return m.Executor + ":" + fmt.Sprintf("%d", m.WithdrawalID)
}

func (m *MsgExecuteWithdrawal) ProtoMessage() {}

type MsgCancelWithdrawal struct {
	Authority    string `json:"authority"`
	WithdrawalID uint64 `json:"withdrawal_id"`
}

func (m *MsgCancelWithdrawal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return ErrInvalidMsg.Wrapf("invalid authority address: %s", err)
	}
	if m.WithdrawalID == 0 {
		return ErrInvalidMsg.Wrap("withdrawal id cannot be zero")
	}
	return nil
}

func (m *MsgCancelWithdrawal) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgCancelWithdrawal) Reset() {
	*m = MsgCancelWithdrawal{}
}

func (m *MsgCancelWithdrawal) String() string {
	return m.Authority + ":" + fmt.Sprintf("%d", m.WithdrawalID)
}

func (m *MsgCancelWithdrawal) ProtoMessage() {}
