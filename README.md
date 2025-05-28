# lengthtagvalue
Go library for parsing your TLV (Length-Tag-Value) format.
* Three bytes Length
* two bytes Tag
* length of Value 

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
