package dns_toy

import (
	"bytes"
	"testing"
)

func TestParseHeader(t *testing.T) {
	b := HeaderToBytes(DNSHeader{
		ID:             0x1234,
		Flags:          0x5678,
		QuestionsCount: 0x9abc,
	})
	header := ParseHeader(bytes.NewBuffer(b))
	if header.ID != 0x1234 {
		t.Errorf("ID should be 0x1234, but got %x", header.ID)
	}
	if header.Flags != 0x5678 {
		t.Errorf("Flags should be 0x5678, but got %x", header.Flags)
	}
}
