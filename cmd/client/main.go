package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/security"
)

const maxAttempts uint64 = 100_000_000
const difficulty = 3 // Difficulty the same as on server

func main() {
	rand.Seed(time.Now().UnixNano())

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading challenge:", err)
		return
	}
	challenge := buffer[:n]

	fmt.Printf("Challenge received: %x\n", challenge)

	startTime := time.Now()
	token := security.GenerateTokenByChallenge(challenge, difficulty, maxAttempts)
	endTime := time.Now()

	if token == nil {
		fmt.Println("PoW solution was not found")
		return
	}

	fmt.Printf("Solution found: %x\n", token)
	fmt.Println("Time taken to solve PoW:", endTime.Sub(startTime))

	_, err = conn.Write(token)
	if err != nil {
		fmt.Printf("Error sending response: %v\n", err)
		return
	}

	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Println("Quote received:", strings.TrimSpace(string(buffer[:n])))
}
