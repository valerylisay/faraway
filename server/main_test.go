package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	reader bufio.Reader
	writer bufio.Writer
}

func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) Read(b []byte) (n int, err error)   { return m.reader.Read(b) }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return m.writer.Write(b) }

func TestIsValidNonce(t *testing.T) {
	tests := []struct {
		name       string
		challenge  string
		difficulty int
		nonce      string
		want       bool
	}{
		{
			name:       "ok",
			challenge:  "easy",
			difficulty: 1,
			nonce:      "20",
			want:       true,
		},
		{
			name:       "invalid nonce",
			challenge:  "invalid",
			difficulty: 1,
			nonce:      "42",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := isValidNonce(tt.challenge, tt.difficulty, tt.nonce)
			if res != tt.want {
				t.Errorf("Expected %v got %v", tt.want, res)
			}
		})
	}
}

func TestSendChallenge(t *testing.T) {
	var buf bytes.Buffer
	conn := &mockConn{writer: *bufio.NewWriter(&buf)}
	*difficulty = 3

	challenge, err := sendChallenge(conn)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := fmt.Sprintf("CHALLENGE:%s:%d\n", challenge, *difficulty)
	conn.writer.Flush()
	if buf.String() != expected {
		t.Errorf("Expected %s got %s", expected, buf.String())
	}
}
