package adbi

import (
	"fmt"
	"sort"
	"strings"
)

type Keyevent uint32

func (k Keyevent) Trigger() []byte {
	if k < end_button_input {
		return []byte(fmt.Sprintf("input keyevent %d\n", uint32(k)))
	}

	switch {
	case k > begin_text_input && k < end_text_input:
		// idk yet
	case k > begin_instant_mouse_input && k < end_instant_mouse_input:
		if k == KEYCODE_MOUSE_SCROLL_UP {
			return []byte(fmt.Sprint("input swipe 100 0 100 900\n"))
		} else {
			return []byte(fmt.Sprint("input swipe 100 900 100 0\n"))
		}
	case k > begin_slow_mouse_input && k < end_slow_mouse_input:
		if k == KEYCODE_MOUSE_SCROLL_UP {
			return []byte(fmt.Sprint("input swipe 100 0 100 900 5000\n"))
		} else {
			return []byte(fmt.Sprint("input swipe 100 900 100 5000\n"))
		}
	}

	return []byte{}
}

func (k Keyevent) TriggerWithRepeat(n int) []byte {
	events := strings.Repeat(fmt.Sprintf(" %d", uint32(k)), n)
	return []byte(fmt.Sprintf("input keyevent %s\n", events))
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
