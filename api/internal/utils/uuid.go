package utils

import (
	"crypto/rand"
	"fmt"
	mathrand "math/rand"
	"time"
)

// GenerateUUID generates a v4 UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to a pseudo-random approach
		rng := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
		for i := range b {
			b[i] = byte(rng.Intn(256))
		}
	}
	
	// Version 4 UUID
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant bits
	
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}