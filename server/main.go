package main

import (
	"bufio"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const clientTimeout = 10 * time.Second

var (
	addr       = flag.String("addr", ":8080", "Listen address")
	difficulty = flag.Int("difficulty", 4, "Challenge difficulty")
)

var quotes = []string{
	"To get somewhere new, you must first decide that you are tired of being where you are.",
	"And suddenly you realize it’s time for a change. You are ready for beginnings.",
	"One small crack does not mean you are broken, it means you were put to the test and didn’t fall apart.",
	"Surround yourself with people… Who illuminate your path. Who pushes you to dig deeper. Who makes you happy? Who makes you laugh?",
	"It might take a day, it might take a year. Just hold onto faith and let go of fear.",
	"Every second is a chance to turn your life around.",
	"Never give up on something you really want, it is difficult to wait, but more difficult to regret.",
	"Always pray to have eyes that see the best, a heart that forgives the worst, a mind that forgets the bad, and a soul that never loses faith.",
	"Sometimes life gets worse before it gets better. But, it always gets better. Just remember who put you there and who helped you up.",
	"Even though I have no idea what tomorrow will bring, I’m going to think positive and hope for the best.",
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s", *addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept conn: %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(clientTimeout))

	challenge, err := sendChallenge(conn)
	if err != nil {
		log.Printf("Failed to send challenge: %v", err)
		return
	}

	if err = validateNonce(conn, challenge); err != nil {
		log.Printf("Failed to validate nonce: %v", err)
		return
	}

	conn.Write([]byte(getRandomQuote() + "\n"))
}

func sendChallenge(conn net.Conn) (string, error) {
	challenge := generateChallenge()
	msg := fmt.Sprintf("CHALLENGE:%s:%d\n", challenge, *difficulty)
	_, err := conn.Write([]byte(msg))
	return challenge, err
}

func generateChallenge() string {
	return fmt.Sprintf("%x", rand.Int63())
}

func validateNonce(conn net.Conn, challenge string) error {
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read nonce: %v", err)
	}

	nonce := scanner.Text()
	if !isValidNonce(challenge, *difficulty, nonce) {
		return errors.New("invalid nonce")
	}

	return nil
}

func isValidNonce(challenge string, difficulty int, nonce string) bool {
	hash := sha256.Sum256([]byte(challenge + nonce))
	target := strings.Repeat("0", difficulty)
	return strings.HasPrefix(fmt.Sprintf("%x", hash), target)
}

func getRandomQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
