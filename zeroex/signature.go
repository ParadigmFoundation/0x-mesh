package zeroex

import (
	"github.com/0xProject/0x-mesh/ethereum/signer"
)

// ECSignatureLength is the length, in bytes, of a ECSignature
const ECSignatureLength = 66

// SignatureType represents the type of 0x signature encountered
type SignatureType uint8

// SignatureType values
const (
	IllegalSignature SignatureType = iota
	InvalidSignature
	EIP712Signature
	EthSignSignature
	WalletSignature
	ValidatorSignature
	PreSignedSignature
	EIP1271WalletSignature
	NSignatureTypesSignature
)

// ECSignatureToBytes converts a 0x ECSignature to it's bytes representation
// Ideally this would be a method on *signer.ECSignature
func ECSignatureToBytes(ecSignature *signer.ECSignature, sigType SignatureType) []byte {
	signature := make([]byte, ECSignatureLength)
	signature[0] = ecSignature.V
	copy(signature[1:33], ecSignature.R[:])
	copy(signature[33:65], ecSignature.S[:])

	// append signature type byte
	signature[65] = byte(sigType)
	return signature
}
