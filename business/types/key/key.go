// Package key represents a key in the system.
package key

import (
	"fmt"
	"regexp"
)

// Key represents a key in the system.
type Key struct {
	value string
}

// String returns the value of the key.
func (n Key) String() string {
	return n.value
}

// Equal provides support for the go-cmp package and testing.
func (n Key) Equal(n2 Key) bool {
	return n.value == n2.value
}

// MarshalText provides support for logging and any marshal needs.
func (n Key) MarshalText() ([]byte, error) {
	return []byte(n.value), nil
}

// =============================================================================

var keyRegEx = regexp.MustCompile("^[a-zA-Z0-9]{1,600}$")

// Parse parses the string value and returns a key if the value complies
// with the rules for a key.
func Parse(value string) (Key, error) {
	if !keyRegEx.MatchString(value) {
		return Key{}, fmt.Errorf("invalid key %q", value)
	}

	return Key{value}, nil
}

// MustParse parses the string value and returns a key if the value
// complies with the rules for a key. If an error occurs the function panics.
func MustParse(value string) Key {
	key, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return key
}

// =============================================================================

// Null represents a key in the system that can be empty.
type Null struct {
	value string
	valid bool
}

// String returns the value of the key.
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

// ParseNull parses the string value and returns a key if the value complies
// with the rules for a key.
func ParseNull(value string) (Null, error) {
	if value == "" {
		return Null{}, nil
	}

	if !keyRegEx.MatchString(value) {
		return Null{}, fmt.Errorf("invalid key %q", value)
	}

	return Null{value, true}, nil
}

// MustParseNull parses the string value and returns a key if the value
// complies with the rules for a key. If an error occurs the function panics.
func MustParseNull(value string) Null {
	key, err := ParseNull(value)
	if err != nil {
		panic(err)
	}

	return key
}
