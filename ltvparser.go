package lengthtagvalue

// ltvparser.go - A simple LTV (Length-Tag-Value) parser in Go
// This package provides functionality to parse LTV encoded data streams.

import (
	"errors"
	"fmt"
	"strconv"
)

// Record represents a parsed TLV record
type Record struct {
	Length int
	Tag    string
	Value  string
}

// Parser handles the parsing of TLV encoded byte streams
type Parser struct {
	data []byte
	pos  int
}

// NewParser creates a new parser for the given data
func NewParser(data []byte) *Parser {
	return &Parser{
		data: data,
		pos:  0,
	}
}

// NewParserFromString creates a new parser from a string
func NewParserFromString(data string) *Parser {
	return &Parser{
		data: []byte(data),
		pos:  0,
	}
}

// HasMore returns true if there's more data to parse
func (p *Parser) HasMore() bool {
	return p.pos < len(p.data)
}

// ParseNext parses the next record from the stream
func (p *Parser) ParseNext() (*Record, error) {
	if !p.HasMore() {
		return nil, errors.New("no more data to parse")
	}

	// Need at least 5 bytes for length (3) + tag (2)
	if p.pos+5 > len(p.data) {
		return nil, errors.New("insufficient data for header")
	}

	// Parse 3-byte length
	lengthStr := string(p.data[p.pos : p.pos+3])
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid length format: %s", lengthStr)
	}
	p.pos += 3

	// Parse 2-byte tag
	tag := string(p.data[p.pos : p.pos+2])
	p.pos += 2

	// Check if we have enough data for the value
	if p.pos+length > len(p.data) {
		return nil, fmt.Errorf("insufficient data for value: need %d bytes, have %d", length, len(p.data)-p.pos)
	}

	// The length includes the tag (2 bytes), so value length is length - 2
	valueLength := length - 2
	if valueLength < 0 {
		return nil, fmt.Errorf("invalid length: %d is too small (must be at least 2 to include tag)", length)
	}

	// Check if we have enough data for the value
	if p.pos+valueLength > len(p.data) {
		return nil, fmt.Errorf("insufficient data for value: need %d bytes, have %d", valueLength, len(p.data)-p.pos)
	}

	// Parse value
	value := string(p.data[p.pos : p.pos+valueLength])
	p.pos += valueLength

	return &Record{
		Length: length,
		Tag:    tag,
		Value:  value,
	}, nil
}

// ParseAll parses all records from the stream
func (p *Parser) ParseAll() ([]*Record, error) {
	var records []*Record

	for p.HasMore() {
		record, err := p.ParseNext()
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}

	return records, nil
}

// Reset resets the parser to the beginning of the data
func (p *Parser) Reset() {
	p.pos = 0
}

// Position returns the current position in the data
func (p *Parser) Position() int {
	return p.pos
}

// Remaining returns the number of bytes remaining to be parsed
func (p *Parser) Remaining() int {
	return len(p.data) - p.pos
}

// ParseString is a convenience function to parse a string directly
func ParseString(data string) ([]*Record, error) {
	parser := NewParserFromString(data)
	return parser.ParseAll()
}

// ParseBytes is a convenience function to parse bytes directly
func ParseBytes(data []byte) ([]*Record, error) {
	parser := NewParser(data)
	return parser.ParseAll()
}
