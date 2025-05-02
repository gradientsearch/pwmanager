// Package entry represents a entry in the system.
package entry

import (
	"fmt"
	"regexp"
)

// Entry represents a entry in the system.
type Entry struct {
	value string
}

// String returns the value of the entry.
func (n Entry) String() string {
	return n.value
}

// Equal provides support for the go-cmp package and testing.
func (n Entry) Equal(n2 Entry) bool {
	return n.value == n2.value
}

// MarshalText provides support for logging and any marshal needs.
func (n Entry) MarshalText() ([]byte, error) {
	return []byte(n.value), nil
}

// =============================================================================

var entryRegEx = regexp.MustCompile("^[a-zA-Z0-9]{1,600}$")

// Parse parses the string value and returns a entry if the value complies
// with the rules for a entry.
func Parse(value string) (Entry, error) {
	if !entryRegEx.MatchString(value) {
		return Entry{}, fmt.Errorf("invalid entry %q", value)
	}

	return Entry{value}, nil
}

// MustParse parses the string value and returns a entry if the value
// complies with the rules for a entry. If an error occurs the function panics.
func MustParse(value string) Entry {
	entry, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return entry
}

// =============================================================================

// Null represents a entry in the system that can be empty.
type Null struct {
	value string
	valid bool
}

// String returns the value of the entry.
func (n Null) String() string {
	if !n.valid {
		return "NULL"
	}

	return n.value
}

// Valid tests if the value is null.
func (n Null) Valid() bool {
	return n.valid
}

// Equal provides support for the go-cmp package and testing.
func (n Null) Equal(n2 Null) bool {
	return n.value == n2.value
}

// =============================================================================

// ParseNull parses the string value and returns a entry if the value complies
// with the rules for a entry.
func ParseNull(value string) (Null, error) {
	if value == "" {
		return Null{}, nil
	}

	if !entryRegEx.MatchString(value) {
		return Null{}, fmt.Errorf("invalid entry %q", value)
	}

	return Null{value, true}, nil
}

// MustParseNull parses the string value and returns a entry if the value
// complies with the rules for a entry. If an error occurs the function panics.
func MustParseNull(value string) Null {
	entry, err := ParseNull(value)
	if err != nil {
		panic(err)
	}

	return entry
}
