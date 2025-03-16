package util

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(publicAddress, signatureHex string, message []byte) (string, error) {
	signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode signature: %w", err)
	}

	if signature[64] < 27 {
		return "", fmt.Errorf("invalid signature: v byte is too small")
	}
	signature[64] -= 27

	hash := accounts.TextHash(message)

	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return "", fmt.Errorf("failed to recover public key: %w", err)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	if strings.ToLower(recoveredAddress) != strings.ToLower(publicAddress) {
		return "", fmt.Errorf("signature does not match the provided address")
	}

	return recoveredAddress, nil
}
