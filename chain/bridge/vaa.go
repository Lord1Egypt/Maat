package bridge

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrInvalidVAAFormat = errors.New("bridge: invalid VAA binary format")
	ErrUnsupportedVAA   = errors.New("bridge: unsupported VAA version or payload")
)

type Signature struct {
	GuardianIndex uint8
	Signature     [65]byte
}

type VAA struct {
	Version           uint8
	GuardianSetIndex  uint32
	Signatures        []Signature
	Timestamp         uint32
	Nonce             uint32
	EmitterChain      uint16
	EmitterAddress    [32]byte
	Sequence          uint64
	ConsistencyLevel  uint8
	Payload           []byte
}

type TokenBridgeTransfer struct {
	PayloadType  uint8
	Amount       [32]byte
	TokenAddress [32]byte
	TokenChain   uint16
	ToAddress    [32]byte
	ToChain      uint16
	Fee          [32]byte
}

func UnmarshalVAA(data []byte) (*VAA, error) {
	if len(data) < 6 {
		return nil, ErrInvalidVAAFormat
	}

	buf := bytes.NewReader(data)

	var version uint8
	if err := binary.Read(buf, binary.BigEndian, &version); err != nil {
		return nil, err
	}
	if version != 1 {
		return nil, fmt.Errorf("%w: version %d", ErrUnsupportedVAA, version)
	}

	var guardianSetIndex uint32
	if err := binary.Read(buf, binary.BigEndian, &guardianSetIndex); err != nil {
		return nil, err
	}

	var lenSignatures uint8
	if err := binary.Read(buf, binary.BigEndian, &lenSignatures); err != nil {
		return nil, err
	}

	sigSize := int(lenSignatures) * 66
	if buf.Len() < sigSize {
		return nil, ErrInvalidVAAFormat
	}

	signatures := make([]Signature, lenSignatures)
	for i := 0; i < int(lenSignatures); i++ {
		var idx uint8
		if err := binary.Read(buf, binary.BigEndian, &idx); err != nil {
			return nil, err
		}
		var sigBytes [65]byte
		if err := binary.Read(buf, binary.BigEndian, &sigBytes); err != nil {
			return nil, err
		}
		signatures[i] = Signature{
			GuardianIndex: idx,
			Signature:     sigBytes,
		}
	}

	if buf.Len() < 51 {
		return nil, ErrInvalidVAAFormat
	}

	var timestamp uint32
	if err := binary.Read(buf, binary.BigEndian, &timestamp); err != nil {
		return nil, err
	}

	var nonce uint32
	if err := binary.Read(buf, binary.BigEndian, &nonce); err != nil {
		return nil, err
	}

	var emitterChain uint16
	if err := binary.Read(buf, binary.BigEndian, &emitterChain); err != nil {
		return nil, err
	}

	var emitterAddress [32]byte
	if err := binary.Read(buf, binary.BigEndian, &emitterAddress); err != nil {
		return nil, err
	}

	var sequence uint64
	if err := binary.Read(buf, binary.BigEndian, &sequence); err != nil {
		return nil, err
	}

	var consistencyLevel uint8
	if err := binary.Read(buf, binary.BigEndian, &consistencyLevel); err != nil {
		return nil, err
	}

	payload := make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &payload); err != nil {
		return nil, err
	}

	return &VAA{
		Version:          version,
		GuardianSetIndex: guardianSetIndex,
		Signatures:       signatures,
		Timestamp:        timestamp,
		Nonce:            nonce,
		EmitterChain:     emitterChain,
		EmitterAddress:   emitterAddress,
		Sequence:         sequence,
		ConsistencyLevel: consistencyLevel,
		Payload:          payload,
	}, nil
}

func ParseTokenBridgeTransfer(payload []byte) (*TokenBridgeTransfer, error) {
	if len(payload) < 133 {
		return nil, ErrInvalidVAAFormat
	}

	buf := bytes.NewReader(payload)

	var payloadType uint8
	if err := binary.Read(buf, binary.BigEndian, &payloadType); err != nil {
		return nil, err
	}
	if payloadType != 1 && payloadType != 3 {
		return nil, fmt.Errorf("%w: payload type %d", ErrUnsupportedVAA, payloadType)
	}

	var amount [32]byte
	if err := binary.Read(buf, binary.BigEndian, &amount); err != nil {
		return nil, err
	}

	var tokenAddress [32]byte
	if err := binary.Read(buf, binary.BigEndian, &tokenAddress); err != nil {
		return nil, err
	}

	var tokenChain uint16
	if err := binary.Read(buf, binary.BigEndian, &tokenChain); err != nil {
		return nil, err
	}

	var toAddress [32]byte
	if err := binary.Read(buf, binary.BigEndian, &toAddress); err != nil {
		return nil, err
	}

	var toChain uint16
	if err := binary.Read(buf, binary.BigEndian, &toChain); err != nil {
		return nil, err
	}

	var fee [32]byte
	if err := binary.Read(buf, binary.BigEndian, &fee); err != nil {
		return nil, err
	}

	return &TokenBridgeTransfer{
		PayloadType:  payloadType,
		Amount:       amount,
		TokenAddress: tokenAddress,
		TokenChain:   tokenChain,
		ToAddress:    toAddress,
		ToChain:      toChain,
		Fee:          fee,
	}, nil
}
