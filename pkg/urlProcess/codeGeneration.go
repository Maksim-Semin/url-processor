package urlProcess

import (
	"crypto/rand"
	"log"
)

const (
	alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits        = "0123456789"
	specialChars  = "_"
	codeLength    = 10
)

// GenerateUniqueCode generates a unique code for creating a shortened code
func GenerateUniqueCode() (string, error) {
	code := make([]byte, codeLength)

	for i := 0; i < codeLength; i++ {
		var chars string

		switch i {
		case 0:
			chars = alphabetLower
		case 1:
			chars = alphabetUpper
		case 2:
			chars = digits
		case 3:
			chars = specialChars
		default:
			chars = alphabetLower + alphabetUpper + digits + specialChars
		}

		code[i] = randomChar(chars)
	}

	return string(code), nil
}

// randomChar returns a random character from the chars string
func randomChar(chars string) byte {
	index := randIndex(len(chars))
	return chars[index]
}

// randIndex generates a random index in the range [0, max)
func randIndex(max int) int {
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Fatalf("error on generate random index: %v", err)
	}
	return int(b[0]) % max
}
