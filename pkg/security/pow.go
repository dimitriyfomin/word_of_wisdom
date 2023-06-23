package security

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

func GenerateChallenge(difficulty int) ([]byte, error) {
	buffer := make([]byte, 8)
	n, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha256.Sum256(append(buffer[:n], []byte(timestamp)...))
	return h[:n], nil
}

func VerifyPoW(challenge, token []byte, difficulty int) bool {
	h := sha256.Sum256(append(challenge, token...))
	for _, v := range h[:difficulty] {
		if v != 0 {
			return false
		}
	}
	return true
}

func GenerateTokenByChallenge(challenge []byte, difficulty int, maxAttempts uint64) []byte {
	var nonce uint64
	timestampBytes := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	for nonce < maxAttempts {
		token := append(timestampBytes, []byte(fmt.Sprint(nonce))...)
		if VerifyPoW(challenge, token[:], difficulty) {
			return token[:]
		}
		nonce++
	}
	return nil
}
