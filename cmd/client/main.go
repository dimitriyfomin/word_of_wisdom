package main

import (
	"errors"
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

	resp, err := processCommand(conn, "v1.CHG", nil)
	if err != nil {
		fmt.Println("Error sending challenge request command:", err)
		return
	}
	fmt.Printf("Challenge: %x\n", resp)
	challenge := resp

	startTime := time.Now()
	token := security.GenerateTokenByChallenge(challenge, difficulty, maxAttempts)
	endTime := time.Now()

	if token == nil {
		fmt.Println("PoW solution was not found")
		return
	}

	fmt.Printf("Solution found: %x\n", token)
	fmt.Println("Time taken to solve PoW:", endTime.Sub(startTime))

	resp, err = processCommand(conn, "v1.CHGT", token)
	if err != nil {
		fmt.Println("Error sending challenge token command:", err)
		return
	}

	resp, err = processCommand(conn, "v1.QTR", nil)
	if err != nil {
		fmt.Println("Error sending random quote:", err)
		return
	}
	fmt.Println("Quote received:", string(resp))
}

func processCommand(conn net.Conn, command string, payload []byte) ([]byte, error) {
	var message []byte
	if payload != nil {
		message = append([]byte(command+" "), payload...)
	} else {
		message = []byte(command)
	}
	_, err := conn.Write(message)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	resp := string(buffer[:n])
	if strings.HasPrefix(resp, "ERR ") {
		return nil, errors.New(resp[4:])
	}

	if strings.HasPrefix(resp, "OK ") {
		return buffer[3:n], nil
	}
	if resp == "OK" {
		return nil, nil
	}
	return nil, errors.New("Unexpected format of response: " + resp)
}
