package adbi

import (
	"fmt"
	"sort"
	"strings"
)

type Keyevent uint32

func (k Keyevent) Trigger() []byte {
	return []byte(fmt.Sprintf("input keyevent %d\n", uint32(k)))
}

func (k Keyevent) Rune() rune {
	return rune(k)
}

// Key returns a Keyevent representing the provided name.
//
// Returns KEYCODE_UNKNOWN for invalid key names.
func Key(name string) Keyevent {
	if code, ok := keycodeLookupTable[strings.ToUpper(name)]; ok {
		return code
	}

	return KEYCODE_UNKNOWN
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
