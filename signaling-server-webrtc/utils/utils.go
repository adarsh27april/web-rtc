package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateShortID(lengths ...int) string {
	length := 6
	if len(lengths) > 0 {
		length = lengths[0]
	}

	id := make([]byte, length)
	for i := range id {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62))))
		if err != nil {
			// Log the fallback
			fmt.Printf("[WARN] crypto/rand failed: %v — falling back to math/rand (less secure)\n", err)
			id[i] = base62[mrand.Intn(len(base62))]
		} else {
			id[i] = base62[n.Int64()]
		}
	}

	return string(id)
}

func Ptr[T any](t T) *T {
	return &t
}
