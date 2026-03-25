package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Метод для генерации sellerId
func GenerateSellerId() int {
	minValue := 111111
	maxValue := 999999
	num, err := rand.Int(rand.Reader, big.NewInt(int64(maxValue-minValue+1)))
	if err != nil {
		panic(err)
	}
	return int(num.Int64()) + minValue
}

// Метод для генерации рандомной строки заданной длины
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		randCharsetIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset)-1)))
		b[i] = charset[randCharsetIndex.Int64()]
	}
	return string(b)
}
