package util

import (
	"log"
	"testing"

	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// 0x785f1ce5a63a4735449caa1b8fb9fb51aaec8e79e870ebefa43140b79ea1954d00986413f31b8a153d52d70d6507b0273bdfcfbdf8f49f44a7d219afd9875ea800

// VerifySignature xác minh chữ ký dựa trên thông điệp và địa chỉ mong đợi
func verifySignature(message, signatureHex, expectedAddress string) bool {
	// Decode chữ ký từ hex
	signature, err := hexutil.Decode(signatureHex)
	if err != nil {
		log.Println("Lỗi decode chữ ký:", err)
		return false
	}

	// Kiểm tra và điều chỉnh byte v (trừ 27 vì đã cộng 27 khi ký)
	if signature[64] < 27 {
		log.Println("Byte v không hợp lệ")
		return false
	}
	signature[64] -= 27

	// Tạo hash của thông điệp giống như khi ký
	hash := accounts.TextHash([]byte(message))

	// Khôi phục public key từ chữ ký và hash
	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		log.Println("Lỗi khôi phục public key:", err)
		return false
	}

	// Lấy địa chỉ từ public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	// So sánh với địa chỉ mong đợi
	return recoveredAddress == expectedAddress
}

// TestSignature kiểm tra việc tạo và xác minh chữ ký
func TestSignature(t *testing.T) {
	// Tạo khóa riêng tư
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		t.Fatal(err)
	}

	// Thông điệp cần ký
	message := "Sign this message to authenticate:"
	hash := accounts.TextHash([]byte(message))

	// Tạo chữ ký
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	// Điều chỉnh byte v để khớp với Ethereum
	signature[64] += 27

	// Chuyển chữ ký thành hex string
	signatureHex := hexutil.Encode(signature)
	fmt.Println("Chữ ký:", signatureHex)

	// Lấy địa chỉ mong đợi từ khóa riêng tư
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	fmt.Println("Địa chỉ mong đợi:", expectedAddress)

	// Xác minh chữ ký
	isValid := verifySignature(message, signatureHex, expectedAddress)
	if !isValid {
		t.Errorf("Chữ ký không hợp lệ, địa chỉ khôi phục không khớp với địa chỉ mong đợi")
	} else {
		fmt.Println("Chữ ký hợp lệ:", isValid)
	}
}
