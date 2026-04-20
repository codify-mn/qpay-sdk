package qpay

import "testing"

func TestMCCCodes_noDuplicates(t *testing.T) {
	seen := make(map[string]int)
	for i, c := range MCCCodes {
		if prev, ok := seen[c.Code]; ok {
			t.Errorf("duplicate code %s at indices %d and %d", c.Code, prev, i)
		}
		seen[c.Code] = i
		if c.NameMongolian == "" {
			t.Errorf("empty name for code %s", c.Code)
		}
	}
}

func TestMCCByCode(t *testing.T) {
	m, ok := MCCByCode("5411")
	if !ok || m.NameMongolian == "" {
		t.Fatalf("expected to find 5411, got %+v ok=%v", m, ok)
	}
	if _, ok := MCCByCode("00000"); ok {
		t.Fatal("expected unknown code to return false")
	}
}

func TestIsValidMCC(t *testing.T) {
	if !IsValidMCC("4814") {
		t.Fatal("4814 must be valid")
	}
	if IsValidMCC("nope") {
		t.Fatal("nope must be invalid")
	}
}
