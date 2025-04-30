// Package bundletype represents the bundle type in the system.
package bundletype

import "fmt"

// The set of types that can be used.
var (
	Personal  = newType("PERSONAL")
	Shareable = newType("SHAREABLE")
)

// =============================================================================

// Set of known housing types.
var bundleTypes = make(map[string]BundleType)

// BundleType represents a type in the system.
type BundleType struct {
	value string
}

func newType(bundleType string) BundleType {
	ht := BundleType{bundleType}
	bundleTypes[bundleType] = ht
	return ht
}

// String returns the name of the type.
func (ht BundleType) String() string {
	return ht.value
}

// Equal provides support for the go-cmp package and testing.
func (ht BundleType) Equal(ht2 BundleType) bool {
	return ht.value == ht2.value
}

// MarshalText provides support for logging and any marshal needs.
func (ht BundleType) MarshalText() ([]byte, error) {
	return []byte(ht.value), nil
}

// =============================================================================

// Parse parses the string value and returns a bundle type if one exists.
func Parse(value string) (BundleType, error) {
	typ, exists := bundleTypes[value]
	if !exists {
		return BundleType{}, fmt.Errorf("invalid bundle type %q", value)
	}

	return typ, nil
}

// MustParse parses the string value and returns a bundle type if one exists. If
// an error occurs the function panics.
func MustParse(value string) BundleType {
	typ, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return typ
}
