package order

import (
	"encoding/binary"
	"math/big"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"

	"github.com/darmayasa221/polymarket-go/internal/commons/polyid"
)

// UnsignedOrder contains the full 12-field struct that gets EIP-712 signed.
// Field names and types match the Polymarket CTF Exchange contract ABI exactly.
type UnsignedOrder struct {
	Salt          *big.Int
	Maker         string // EOA address
	Signer        string // same as Maker for EOA
	Taker         string // zero address for CLOB orders
	TokenID       polyid.TokenID
	MakerAmount   decimal.Decimal
	TakerAmount   decimal.Decimal
	Expiration    int64
	Nonce         int64
	FeeRateBps    uint64
	Side          Side
	SignatureType uint8
}

// EIP-712 domain and type constants for Polymarket CTF Exchange.
const (
	eip712DomainName    = "Polymarket CTF Exchange"
	eip712DomainVersion = "1"
	// exchangeContract is the main (non-neg-risk) exchange on Polygon mainnet.
	exchangeContract = "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E"
)

// keccak256 computes the Ethereum Keccak256 hash (not standard SHA3).
func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

// SigningHash computes the EIP-712 signing hash for an UnsignedOrder.
// The returned 32 bytes are what the EOA private key signs via ECDSA.
// Port of: https://github.com/Polymarket/clob-client/blob/main/src/signing/eip712.ts
func SigningHash(o UnsignedOrder, chainID int64) ([]byte, error) {
	domainSeparator := computeDomainSeparator(chainID)
	structHash := computeStructHash(o)

	// EIP-712 final hash: keccak256("\x19\x01" + domainSeparator + structHash)
	msg := make([]byte, 0, 66)
	msg = append(msg, 0x19, 0x01)
	msg = append(msg, domainSeparator...)
	msg = append(msg, structHash...)
	return keccak256(msg), nil
}

// domainTypeHash is the EIP-712 type hash for the domain separator.
var domainTypeHash = keccak256([]byte(
	"EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)",
))

// orderTypeHash is the EIP-712 type hash for the Order struct.
var orderTypeHash = keccak256([]byte(
	"Order(uint256 salt,address maker,address signer,address taker,uint256 tokenId," +
		"uint256 makerAmount,uint256 takerAmount,uint256 expiration,uint256 nonce," +
		"uint256 feeRateBps,uint8 side,uint8 signatureType)",
))

func computeDomainSeparator(chainID int64) []byte {
	nameHash := keccak256([]byte(eip712DomainName))
	versionHash := keccak256([]byte(eip712DomainVersion))
	contractHash := padAddress(exchangeContract)
	chainBig := new(big.Int).SetInt64(chainID)

	data := make([]byte, 0, 5*32)
	data = append(data, pad32(domainTypeHash)...)
	data = append(data, pad32(nameHash)...)
	data = append(data, pad32(versionHash)...)
	data = append(data, padBigInt(chainBig)...)
	data = append(data, contractHash...)
	return keccak256(data)
}

func computeStructHash(o UnsignedOrder) []byte {
	makerHash := padAddress(o.Maker)
	signerHash := padAddress(o.Signer)
	takerHash := padAddress(o.Taker)
	tokenIDBig, _ := new(big.Int).SetString(o.TokenID.String(), 10)
	makerAmountBig := decimalToUint256(o.MakerAmount)
	takerAmountBig := decimalToUint256(o.TakerAmount)

	data := make([]byte, 0, 12*32)
	data = append(data, pad32(orderTypeHash)...)
	data = append(data, padBigInt(o.Salt)...)
	data = append(data, makerHash...)
	data = append(data, signerHash...)
	data = append(data, takerHash...)
	data = append(data, padBigInt(tokenIDBig)...)
	data = append(data, padBigInt(makerAmountBig)...)
	data = append(data, padBigInt(takerAmountBig)...)
	data = append(data, padInt64(o.Expiration)...)
	data = append(data, padInt64(o.Nonce)...)
	data = append(data, padUint64(o.FeeRateBps)...)
	data = append(data, padUint8(uint8(o.Side))...)
	data = append(data, padUint8(o.SignatureType)...)
	return keccak256(data)
}

// pad32 zero-pads a byte slice to 32 bytes (left-padded).
func pad32(b []byte) []byte {
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}

func padBigInt(n *big.Int) []byte {
	if n == nil {
		return make([]byte, 32)
	}
	return pad32(n.Bytes())
}

func padAddress(addr string) []byte {
	if len(addr) >= 2 && addr[:2] == "0x" {
		addr = addr[2:]
	}
	b := make([]byte, 20)
	for i := 0; i+1 < len(addr) && i < 40; i += 2 {
		hi := hexByte(addr[i])
		lo := hexByte(addr[i+1])
		b[i/2] = (hi << 4) | lo
	}
	return pad32(b)
}

func padInt64(n int64) []byte {
	out := make([]byte, 32)
	binary.BigEndian.PutUint64(out[24:], uint64(n)) //nolint:gosec // EIP-712 fields are always non-negative
	return out
}

func padUint64(n uint64) []byte {
	out := make([]byte, 32)
	binary.BigEndian.PutUint64(out[24:], n)
	return out
}

func padUint8(n uint8) []byte {
	out := make([]byte, 32)
	out[31] = n
	return out
}

func hexByte(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

// decimalToUint256 converts a decimal amount to its uint256 big.Int representation.
// CLOB amounts are in USDC.e with 6 decimals — multiply by 10^6.
func decimalToUint256(d decimal.Decimal) *big.Int {
	scaled := d.Mul(decimal.NewFromInt(1_000_000))
	bi, _ := new(big.Int).SetString(scaled.StringFixed(0), 10)
	return bi
}
