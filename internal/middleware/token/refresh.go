package token

import (
	"crypto/rand"

	"github.com/mpu-cad/gw-backend-go/internal/entity"
)

func GenerateRefreshToken() string {
	b := make([]byte, entity.LenRefreshToken)

	_, _ = rand.Read(b)

	token := make([]rune, entity.LenRefreshToken)
	for i := range entity.LenRefreshToken {
		index := int(b[i]) % len(entity.AllSymbol)
		token[i] = entity.AllSymbol[index]
	}

	return string(token)
}
