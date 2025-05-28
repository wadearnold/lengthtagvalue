# lengthtagvalue
Go library for parsing your TLV (Tag-Length-Value) format.

```go
// Simple parsing
records, err := tlvparser.ParseString(yourData)

// Stateful parsing
parser := tlvparser.NewParserFromString(yourData)
for parser.HasMore() {
    record, err := parser.ParseNext()
    // Process record...
}
```