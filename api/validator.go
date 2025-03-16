package api

import (
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
)

func IsValidEthereumAddress(address string) bool {
	if len(address) != 42 || address[:2] != "0x" {
		return false
	}
	matched, _ := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	if !matched {
		return false
	}
	return common.IsHexAddress(address)
}

var EthAddressValidator validator.Func = func(fl validator.FieldLevel) bool {
	address := fl.Field().String()
	return IsValidEthereumAddress(address)
}
