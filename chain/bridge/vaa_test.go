package bridge

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestParseVAA(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteByte(1) // version
	binary.Write(&buf, binary.BigEndian, uint32(2)) // guardian set index
	buf.WriteByte(1) // signatures length
	buf.WriteByte(0) // guardian index
	var dummySig [65]byte
	dummySig[0] = 0xAA
	dummySig[64] = 0xBB
	buf.Write(dummySig[:])

	binary.Write(&buf, binary.BigEndian, uint32(100000)) // timestamp
	binary.Write(&buf, binary.BigEndian, uint32(42))     // nonce
	binary.Write(&buf, binary.BigEndian, uint16(2))      // emitter chain
	var dummyEmitter [32]byte
	dummyEmitter[0] = 0x11
	buf.Write(dummyEmitter[:])
	binary.Write(&buf, binary.BigEndian, uint64(999))    // sequence
	buf.WriteByte(15)                                    // consistency level

	// Payload (TokenBridgeTransfer layout)
	buf.WriteByte(1) // payload type
	var dummyAmount [32]byte
	dummyAmount[31] = 0xFF
	buf.Write(dummyAmount[:])
	var dummyToken [32]byte
	dummyToken[31] = 0xEE
	buf.Write(dummyToken[:])
	binary.Write(&buf, binary.BigEndian, uint16(2)) // token chain
	var dummyTo [32]byte
	dummyTo[31] = 0xDD
	buf.Write(dummyTo[:])
	binary.Write(&buf, binary.BigEndian, uint16(1)) // to chain
	var dummyFee [32]byte
	buf.Write(dummyFee[:])

	v, err := UnmarshalVAA(buf.Bytes())
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v.Version != 1 || v.GuardianSetIndex != 2 || len(v.Signatures) != 1 {
		t.Fatalf("unexpected VAA fields: %+v", v)
	}

	if v.Timestamp != 100000 || v.Sequence != 999 {
		t.Fatalf("unexpected VAA timestamp/sequence: %+v", v)
	}

	transfer, err := ParseTokenBridgeTransfer(v.Payload)
	if err != nil {
		t.Fatalf("parse transfer error: %v", err)
	}

	if transfer.PayloadType != 1 || transfer.TokenChain != 2 || transfer.ToChain != 1 {
		t.Fatalf("unexpected transfer fields: %+v", transfer)
	}

	if transfer.Amount[31] != 0xFF || transfer.TokenAddress[31] != 0xEE || transfer.ToAddress[31] != 0xDD {
		t.Fatalf("unexpected transfer addresses/amount: %+v", transfer)
	}
}

func TestParseVAAErrors(t *testing.T) {
	_, err := UnmarshalVAA([]byte{1, 2})
	if err == nil {
		t.Fatal("expected error for short VAA")
	}

	_, err = ParseTokenBridgeTransfer([]byte{1, 2})
	if err == nil {
		t.Fatal("expected error for short transfer payload")
	}
}
