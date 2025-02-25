package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const solveTimeout = 10 * time.Second

var serverAddr = flag.String("server-addr", ":8080", "Server address")

type Challenge struct {
	value      string
	difficulty int
}

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	challenge, err := readChallenge(conn)
	if err != nil {
		log.Fatalf("Failed to read challenge: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), solveTimeout)
	defer cancel()
	nonce, err := solveChallenge(ctx, challenge)
	if err != nil {
		log.Fatalf("Failed to solve challenge: %v", err)
	}

	quote, err := exchangeNonceToQuote(conn, nonce)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	fmt.Printf("Received quote: %s", quote)
}

func readChallenge(conn net.Conn) (*Challenge, error) {
	scanner := bufio.NewScanner(conn)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	parts := strings.Split(scanner.Text(), ":")
	if len(parts) != 3 && parts[0] != "CHALLENGE" {
		return nil, errors.New("invalid challenge format")
	}

	difficulty, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errors.New("invalid difficulty")
	}

	return &Challenge{value: parts[1], difficulty: difficulty}, nil
}

func solveChallenge(ctx context.Context, challenge *Challenge) (string, error) {
	var nonce int
	target := strings.Repeat("0", challenge.difficulty)

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("solve timeout")
		default:
			hash := sha256.Sum256([]byte(challenge.value + strconv.Itoa(nonce)))
			if strings.HasPrefix(fmt.Sprintf("%x", hash), target) {
				return strconv.Itoa(nonce), nil
			}
			nonce++
		}
	}
}

func exchangeNonceToQuote(conn net.Conn, nonce string) (string, error) {
	conn.Write([]byte(nonce + "\n"))

	scanner := bufio.NewScanner(conn)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return scanner.Text(), nil
}
