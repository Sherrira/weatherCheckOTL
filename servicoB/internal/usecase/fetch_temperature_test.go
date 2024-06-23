package usecase

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestIsValidCEP(t *testing.T) {
	assert.True(t, IsValidCEP("12345678"))
	assert.False(t, IsValidCEP("1234567"))
	assert.False(t, IsValidCEP("123456789"))
	assert.False(t, IsValidCEP("1234567a"))
}

func TestIsValidCEPFuzzy(t *testing.T) {
	assert := assert.New(t)
	SEED := time.Now().UnixNano()
	r := rand.New(rand.NewSource(SEED))
	fmt.Println(r.Uint64())
	fmt.Println(r.Uint64())

	for i := 0; i < 100; i++ {
		cep := fmt.Sprintf("%08d", rand.Intn(1000000000))
		isValid := IsValidCEP(cep)
		if isValid {
			assert.Len(cep, 8, "The CEP should be valid")
		} else {
			assert.NotEqual(8, len(cep), "The CEP should be invalid")
		}
	}
}
