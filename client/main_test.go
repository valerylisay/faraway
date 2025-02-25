package main

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"strings"
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

func TestReadChallenge(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValue   string
		wantDiff    int
		expectError bool
	}{
		{
			name:      "ok",
			input:     "CHALLENGE:abc123:4\n",
			wantValue: "abc123",
			wantDiff:  4,
		},
		{
			name:        "invalid format",
			input:       "INVALIDFORMAT\n",
			expectError: true,
		},
		{
			name:        "invalid difficulty",
			input:       "CHALLENGE:abc:diff\n",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConn{
				reader: *bufio.NewReader(bytes.NewBufferString(tt.input)),
			}

			challenge, err := readChallenge(conn)
			if tt.expectError {
				if err == nil {
					t.Fatal("Expected err got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if challenge.value != tt.wantValue {
				t.Errorf("Expected value %s got %s", tt.wantValue, challenge.value)
			}

			if challenge.difficulty != tt.wantDiff {
				t.Errorf("Expected difficulty %d got %d", tt.wantDiff, challenge.difficulty)
			}
		})
	}
}

func TestSolveChallenge(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		challenge   *Challenge
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "ok",
			challenge: &Challenge{
				value:      "easy",
				difficulty: 1,
			},
			timeout: time.Second,
		},
		{
			name: "solve timeout",
			challenge: &Challenge{
				value:      "difficult",
				difficulty: 10,
			},
			timeout:     time.Millisecond,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(ctx, tt.timeout)
			defer cancel()

			_, err := solveChallenge(ctx, tt.challenge)
			if tt.expectError {
				if err == nil {
					t.Fatal("Expected err got nil")
				}
				if err.Error() != "solve timeout" {
					t.Fatalf("Unexpected error: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}

func TestExchangeNonceToQuote(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		nonce  string
		expect string
	}{
		{
			name:   "ok",
			input:  "test quote\n",
			nonce:  "123",
			expect: "test quote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			conn := &mockConn{
				reader: *bufio.NewReader(bytes.NewBufferString(tt.input)),
				writer: *bufio.NewWriter(&buf),
			}

			quote, err := exchangeNonceToQuote(conn, tt.nonce)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if quote != tt.expect {
				t.Errorf("Expected quote %s got %s", tt.expect, quote)
			}

			conn.writer.Flush()
			sent := buf.String()
			if !strings.Contains(sent, tt.nonce+"\n") {
				t.Errorf("Expected nonce %s got %s", tt.nonce, sent)
			}
		})
	}
}
