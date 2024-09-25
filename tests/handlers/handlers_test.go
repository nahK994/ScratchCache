package handlers_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nahK994/TinyCache/pkg/handlers"
)

func TestHandleGET(t *testing.T) {
	// Key does not exist
	_, err := handlers.HandleCommand("*2\r\n$3\r\nGET\r\n$17\r\nnon_existing_key\r\n")
	if err == nil {
		t.Errorf("Expected error for non-existing key, got none")
	}

	// Key exists, type int
	handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$6\r\nnumber\r\n$2\r\n10\r\n")
	resp, err := handlers.HandleCommand("*2\r\n$3\r\nGET\r\n$6\r\nnumber\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != "$2\r\n10\r\n" {
		t.Errorf("Expected '$2\\r\\n10\\r\\n', got %s", resp)
	}

	// Key exists, type string
	handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n")
	resp, err = handlers.HandleCommand("*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != "$5\r\nhello\r\n" {
		t.Errorf("Expected '$5\\r\\nhello\\r\\n', got %s", resp)
	}
}

func TestHandleSET(t *testing.T) {
	// Set key and get its value
	resp, err := handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != "+OK\r\n" {
		t.Errorf("Expected '+OK\\r\\n', got %s", resp)
	}

	resp, err = handlers.HandleCommand("*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != "$5\r\nvalue\r\n" {
		t.Errorf("Expected '$5\\r\\nvalue\\r\\n', got %s", resp)
	}
}

func TestHandleEXISTS(t *testing.T) {
	// Key does not exist
	resp, err := handlers.HandleCommand("*2\r\n$6\r\nEXISTS\r\n$6\r\nno_key\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":0\r\n" {
		t.Errorf("Expected ':0\\r\\n', got %s", resp)
	}

	// Key exists
	handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n")
	resp, err = handlers.HandleCommand("*2\r\n$6\r\nEXISTS\r\n$5\r\nmykey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":1\r\n" {
		t.Errorf("Expected ':1\\r\\n', got %s", resp)
	}
}

func TestHandleDEL(t *testing.T) {
	// Key does not exist
	resp, err := handlers.HandleCommand("*2\r\n$3\r\nDEL\r\n$6\r\nno_key\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":0\r\n" {
		t.Errorf("Expected ':0\\r\\n', got %s", resp)
	}

	// Key exists, and is deleted
	handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nvalue\r\n")
	resp, err = handlers.HandleCommand("*2\r\n$3\r\nDEL\r\n$5\r\nmykey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":1\r\n" {
		t.Errorf("Expected ':1\\r\\n', got %s", resp)
	}

	// Key should no longer exist
	resp, err = handlers.HandleCommand("*2\r\n$6\r\nEXISTS\r\n$5\r\nmykey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":0\r\n" {
		t.Errorf("Expected ':0\\r\\n', got %s", resp)
	}
}

func TestHandleINCR(t *testing.T) {
	// Key does not exist, should initialize to 0 and increment
	resp, err := handlers.HandleCommand("*2\r\n$4\r\nINCR\r\n$6\r\nnewkey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":1\r\n" {
		t.Errorf("Expected ':1\\r\\n', got %s", resp)
	}

	// Key exists, increment further
	resp, err = handlers.HandleCommand("*2\r\n$4\r\nINCR\r\n$6\r\nnewkey\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":2\r\n" {
		t.Errorf("Expected ':2\\r\\n', got %s", resp)
	}

	// Error on wrong type (set string)
	handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n")
	_, err = handlers.HandleCommand("*2\r\n$4\r\nINCR\r\n$5\r\nmykey\r\n")
	if err == nil {
		t.Errorf("Expected error for INCR on non-integer key, got none")
	}
}

func TestHandleLPUSH(t *testing.T) {
	// Push elements to a list
	resp, err := handlers.HandleCommand("*5\r\n$5\r\nLPUSH\r\n$6\r\nmylist\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":3\r\n" {
		t.Errorf("Expected ':3\\r\\n', got %s", resp)
	}
}

func TestHandleLRANGE(t *testing.T) {
	_, err := handlers.HandleCommand("*4\r\n$6\r\nLRANGE\r\n$6\r\nmylist1\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	if err == nil {
		t.Errorf("Expected empty list error, got %v", err)
	}

	handlers.HandleCommand("*5\r\n$5\r\nLPUSH\r\n$6\r\nmylist1\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n")
	resp, err1 := handlers.HandleCommand("*4\r\n$6\r\nLRANGE\r\n$6\r\nmylist1\r\n$1\r\n0\r\n$2\r\n-1\r\n")
	if err1 != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !contains(resp, "$1\r\na\r\n") {
		t.Errorf("Expected response to contain '$1\\r\\na\\r\\n', got %s", resp)
	}

}

func TestHandleEXPIRE(t *testing.T) {
	// Set key with an expiration
	resp, err := handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$7\r\nexp_key\r\n$5\r\nvalue\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	resp, err = handlers.HandleCommand("*3\r\n$6\r\nEXPIRE\r\n$7\r\nexp_key\r\n$1\r\n5\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != "+OK\r\n" { // 1 if the key was set to expire
		t.Errorf("Expected +OK\r\n, got %s", resp)
	}

	// Check the TTL
	resp, err = handlers.HandleCommand("*2\r\n$3\r\nTTL\r\n$7\r\nexp_key\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resp != ":4\r\n" { // Should reflect the remaining time
		t.Errorf("Expected TTL response to be ':5\\r\\n', got %s", resp)
	}

	// Wait for expiration
	time.Sleep(6 * time.Second)
	_, err = handlers.HandleCommand("*2\r\n$3\r\nTTL\r\n$7\r\nexp_key\r\n")
	if err == nil {
		t.Errorf("Expected type error error, got %v", err)
	}
}

// Helper function to check if a string contains another string
func contains(resp, substr string) bool {
	return strings.Contains(resp, substr)
}

func TestHandleTTL(t *testing.T) {
	// Set key with expiration
	handlers.HandleCommand("*3\r\n$3\r\nSET\r\n$7\r\nttl_key\r\n$5\r\nvalue\r\n")
	handlers.HandleCommand("*3\r\n$6\r\nEXPIRE\r\n$7\r\nttl_key\r\n$2\r\n10\r\n")

	// Check TTL value
	resp, err := handlers.HandleCommand("*2\r\n$3\r\nTTL\r\n$7\r\nttl_key\r\n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Expected response should be a positive integer representing remaining seconds
	if !isPositiveInteger(resp) {
		t.Errorf("Expected positive integer, got %s", resp)
	}

	// Wait for the key to expire
	time.Sleep(15 * time.Second)
	_, err = handlers.HandleCommand("*2\r\n$3\r\nTTL\r\n$7\r\nttl_key\r\n")
	if err == nil {
		t.Errorf("Expected type error, got %v", err)
	}
}

// Helper function to check if the response is a positive integer
func isPositiveInteger(resp string) bool {
	if len(resp) < 3 {
		return false
	}
	if resp[0] != ':' {
		return false
	}
	num, err := strconv.Atoi(resp[1 : len(resp)-2]) // Remove ':', and CRLF
	return err == nil && num > 0
}
