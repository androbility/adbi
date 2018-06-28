package adbi

import (
	"fmt"
	"sort"
	"strings"
)

type Keycode uint32

func (k Keycode) Event() ([]byte, bool) {
	if code, ok := keymap[rune(k)]; ok {
		return []byte(fmt.Sprintf("input keyevent %d\n", code)), true
	}

	return nil, false
}

func (k Keycode) Rune() rune {
	return rune(k)
}

// Key returns a Keycode representing the provided name.
//
// Returns KEYCODE_UNKNOWN and an error for invalid key names.
func Key(name string) (Keycode, error) {
	if code, ok := keycodeLookupTable[strings.ToUpper(name)]; ok {
		return code, nil
	}

	return KEYCODE_UNKNOWN, fmt.Errorf("invalid Keycode: %s", name)
}

// KeyNames returns a sorted slice of valid key names.
func KeyNames() []string {
	names := []string{}
	for name, _ := range keycodeLookupTable {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}
