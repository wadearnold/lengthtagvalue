package lengthtagvalue

import (
	"testing"
)

func TestParseString(t *testing.T) {
	testData := "041500360103101V00100130221262615326738433730068005100148151402300008001182V 2049830446200301D01702305140774154494006037VRV00423P 01086642V00VC0178830514077415449400699Visa"

	records, err := ParseString(testData)
	if err != nil {
		t.Fatalf("ParseString failed: %v", err)
	}

	// Verify we parsed some records
	if len(records) == 0 {
		t.Fatal("No records parsed")
	}

	t.Logf("Parsed %d records:", len(records))
	for i, record := range records {
		t.Logf("Record %d: Length=%d, Tag=%s, Value=%s", i+1, record.Length, record.Tag, record.Value)
	}

	// Test specific first record
	if len(records) > 0 {
		first := records[0]
		if first.Length != 41 {
			t.Errorf("Expected length 41, got %d", first.Length)
		}
		if first.Tag != "50" {
			t.Errorf("Expected tag '50', got '%s'", first.Tag)
		}
		expectedValue := "0360103101V0010013022126261532673843373"
		if first.Value != expectedValue {
			t.Errorf("Expected value '%s', got '%s'", expectedValue, first.Value)
		}
	}
}

func TestParserMethods(t *testing.T) {
	testData := "00701Hello00702World"
	parser := NewParserFromString(testData)

	// Test HasMore
	if !parser.HasMore() {
		t.Error("Expected HasMore to be true initially")
	}

	// Parse first record
	record1, err := parser.ParseNext()
	if err != nil {
		t.Fatalf("Failed to parse first record: %v", err)
	}

	if record1.Length != 7 || record1.Tag != "01" || record1.Value != "Hello" {
		t.Errorf("First record incorrect: Length=%d, Tag=%s, Value=%s", record1.Length, record1.Tag, record1.Value)
	}

	// Test position - should be at 5 (header) + 5 (value) = 10
	if parser.Position() != 10 {
		t.Errorf("Expected position 10, got %d", parser.Position())
	}

	// Parse second record
	record2, err := parser.ParseNext()
	if err != nil {
		t.Fatalf("Failed to parse second record: %v", err)
	}

	if record2.Length != 7 || record2.Tag != "02" || record2.Value != "World" {
		t.Errorf("Second record incorrect: Length=%d, Tag=%s, Value=%s", record2.Length, record2.Tag, record2.Value)
	}

	// Should be at end now
	if parser.HasMore() {
		t.Error("Expected HasMore to be false at end")
	}
}

func TestParserReset(t *testing.T) {
	testData := "00701Hello"
	parser := NewParserFromString(testData)

	// Parse once
	_, err := parser.ParseNext()
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Reset and parse again
	parser.Reset()
	if parser.Position() != 0 {
		t.Error("Reset didn't reset position to 0")
	}

	record, err := parser.ParseNext()
	if err != nil {
		t.Fatalf("Failed to parse after reset: %v", err)
	}

	if record.Value != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", record.Value)
	}
}

func TestErrorCases(t *testing.T) {
	// Test insufficient data for header
	parser := NewParserFromString("123")
	_, err := parser.ParseNext()
	if err == nil {
		t.Error("Expected error for insufficient header data")
	}

	// Test insufficient data for value
	parser = NewParserFromString("10001") // Says length 100 but only 1 byte follows (need 98 for value)
	_, err = parser.ParseNext()
	if err == nil {
		t.Error("Expected error for insufficient value data")
	}

	// Test invalid length format
	parser = NewParserFromString("ABC01data")
	_, err = parser.ParseNext()
	if err == nil {
		t.Error("Expected error for invalid length format")
	}
}

func BenchmarkParseString(b *testing.B) {
	testData := "041500360103101V00100130221262615326738433730068005100148151402300008001182V 2049830446200301D01702305140774154494006037VRV00423P 01086642V00VC0178830514077415449400699Visa"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseString(testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
