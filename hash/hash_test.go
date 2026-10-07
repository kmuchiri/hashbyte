package hash

import (
	"bytes"
	"testing"
)

func TestNewBuffers(t *testing.T) {
	buf1, buf2, buf3 := NewBuffers()
	if len(buf1) != 20 {
		t.Errorf("expected buf1 length 20, got %d", len(buf1))
	}
	if len(buf2) != 40 {
		t.Errorf("expected buf2 length 40, got %d", len(buf2))
	}
	if len(buf3) != 32 {
		t.Errorf("expected buf3 length 32, got %d", len(buf3))
	}
	if !bytes.Equal(buf1[:16], FixedKey) {
		t.Errorf("expected buf1 to start with FixedKey")
	}
	if !bytes.Equal(buf2[:16], FixedKey) {
		t.Errorf("expected buf2 to start with FixedKey")
	}
}

func TestAudible(t *testing.T) {
	buf1, buf2, buf3 := NewBuffers()
	candidate := uint32(0x12345678)
	res := Audible(candidate, buf1, buf2, buf3)
	if len(res) != 20 {
		t.Errorf("expected result length 20")
	}
}
