package zeroex

import "github.com/ethereum/go-ethereum/signer/core"

const (
	// ZeroExProtocolName is the EIP-712 domain name of the 0x protocol
	ZeroExProtocolName = "0x Protocol"

	// ZeroExProtocolVersion is the EIP-712 domain version of the 0x protocol
	ZeroExProtocolVersion = "3.0.0"

	// TypeEIP712Domain is the name of the EIP-712 domain type
	TypeEIP712Domain = "EIP712Domain"

	// TypeZeroExTransaction is the name of the 0x transaction type
	TypeZeroExTransaction = "ZeroExTransaction"

	// TypeOrder is the name of the 0x order type
	TypeOrder = "Order"
)

// EIP712Types are the EIP-712 type definitions for the relevant 0x types and domain
var EIP712Types = core.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"ZeroExTransaction": {
		{
			Name: "salt",
			Type: "uint256",
		},
		{
			Name: "expirationTimeSeconds",
			Type: "uint256",
		},
		{
			Name: "gasPrice",
			Type: "uint256",
		},
		{
			Name: "signerAddress",
			Type: "address",
		},
		{
			Name: "data",
			Type: "bytes",
		},
	},
	"Order": {
		{
			Name: "makerAddress",
			Type: "address",
		},
		{
			Name: "takerAddress",
			Type: "address",
		},
		{
			Name: "feeRecipientAddress",
			Type: "address",
		},
		{
			Name: "senderAddress",
			Type: "address",
		},
		{
			Name: "makerAssetAmount",
			Type: "uint256",
		},
		{
			Name: "takerAssetAmount",
			Type: "uint256",
		},
		{
			Name: "makerFee",
			Type: "uint256",
		},
		{
			Name: "takerFee",
			Type: "uint256",
		},
		{
			Name: "expirationTimeSeconds",
			Type: "uint256",
		},
		{
			Name: "salt",
			Type: "uint256",
		},
		{
			Name: "makerAssetData",
			Type: "bytes",
		},
		{
			Name: "takerAssetData",
			Type: "bytes",
		},
		{
			Name: "makerFeeAssetData",
			Type: "bytes",
		},
		{
			Name: "takerFeeAssetData",
			Type: "bytes",
		},
	},
}
