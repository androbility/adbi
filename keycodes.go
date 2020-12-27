package adbi

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type Keyevent uint32

func (k Keyevent) Trigger() []byte {
	return k.TriggerWithRepeat(1)
}

func (k Keyevent) TriggerWithRepeat(n int) []byte {
	if k < end_button_input {
		events := strings.Repeat(fmt.Sprintf(" %d", uint32(k)), n)
		return []byte(fmt.Sprintf("input keyevent %s\n", strings.TrimFunc(events, unicode.IsSpace)))
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
		if k == KEYCODE_MOUSE_SCROLL_UP_READING_SPEED {
			return []byte(fmt.Sprint("input swipe 100 0 100 900 5000\n"))
		} else {
			return []byte(fmt.Sprint("input swipe 100 900 100 0 5000\n"))
		}
	}

	return nil
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
