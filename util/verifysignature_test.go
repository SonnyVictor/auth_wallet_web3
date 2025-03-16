package util

import (
	"log"
	"testing"

	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func verifySignature(message, signatureHex, expectedAddress string) bool {
	signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		log.Println("Err decode:", err)
		return false
	}

	if signature[64] < 27 {
		log.Println("Byte v invalid")
		return false
	}
	signature[64] -= 27

	hash := accounts.TextHash([]byte(message))

	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		log.Println(" public key:", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	return recoveredAddress == expectedAddress
}

func TestSignature(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		t.Fatal(err)
	}

	message := "Sign this message to authenticate:"
	hash := accounts.TextHash([]byte(message))

	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	signature[64] += 27
	signatureHex := hexutil.Encode(signature)
	fmt.Println("signatureHex", signatureHex)

	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	fmt.Println("Address expectedAddress", expectedAddress)

	isValid := verifySignature(message, signatureHex, expectedAddress)
	if !isValid {
		t.Errorf("Signature is valid")
	} else {
		fmt.Println("Correctly:", isValid)
	}
}
