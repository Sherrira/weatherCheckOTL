package usecase

import (
	"strconv"
)


func IsValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	_, err := strconv.Atoi(cep)
	return err == nil
}
