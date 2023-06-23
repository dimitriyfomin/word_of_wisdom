package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/security"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/wisdom"
)

const difficulty = 3 // We set the difficulty of the PoW

func main() {
	rand.Seed(time.Now().UnixNano())

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Server is running!")
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	challenge, err := security.GenerateChallenge(difficulty)
	if err != nil {
		fmt.Printf("Error generating challenge: %v\n", err)
		return
	}

	_, err = conn.Write(challenge)
	if err != nil {
		fmt.Printf("Error sending challenge: %v\n", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if !security.VerifyPoW(challenge, buffer[:n], difficulty) {
		fmt.Println("Proof of work verification failed")
		return
	}
	fmt.Println("Proof of work verification successful")

	quoteToSend := wisdom.GetRandomQuote()

	fmt.Println("Sending the quote:", quoteToSend)

	_, err = conn.Write([]byte(quoteToSend + "\n"))
	if err != nil {
		fmt.Printf("Error sending quote: %v\n", err)
		return
	}
}
