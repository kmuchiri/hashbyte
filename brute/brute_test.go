package brute

import (
	"encoding/hex"
	"testing"
)

func TestFindHash(t *testing.T) {
	targetHashHex := "05768689284d85b1abe78176134439220f87c209"
	targetHash, _ := hex.DecodeString(targetHashHex)
	match, _, found := findHash(targetHash)
	
	if !found {
		t.Errorf("expected to find match")
	}
	expectedMatch := [4]byte{0, 0, 0, 0}
	if match != expectedMatch {
		t.Errorf("expected match %v, got %v", expectedMatch, match)
	}
}

func TestPrintResult(t *testing.T) {
	printResult("05768689284d85b1abe78176134439220f87c209", [4]byte{0, 0, 0, 0}, true, 0)
	printResult("05768689284d85b1abe78176134439220f87c209", [4]byte{0, 0, 0, 0}, false, 0)
}
