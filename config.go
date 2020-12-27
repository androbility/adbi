package adbi

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
)

// LoadConfigFile attempts to read a config-file from configDir.
// Failing that, LoadConfigFile parses defaultBindings for the
// configuration.
func LoadConfigFile(configDir, defaultBindings string) map[rune]Keyevent {
	configFile := os.ExpandEnv(configDir) + "/config.json"

	// 1. Open configFile; fallback to defaultBindings on error
	// 2. get io.ReadCloser from either.
	// 3. Unmarshal the json.
	// 4. Cast values as strings.
	keymapAsBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Debug("Error reading config file.  Using default keybindings.")

		keymapAsBytes = []byte(defaultBindings)
	}
	keymapReader := strings.NewReader(string(keymapAsBytes))

	keybindingsConfig, err := ioutil.ReadAll(keymapReader)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error reading keybindings.  Aborting startup.")
	}

	// Now json.Unmarshal.
	var keybindingsContainer map[string]interface{}
	err = json.Unmarshal(keybindingsConfig, &keybindingsContainer)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error unmarshaling the configuration.  Aborting startup.")
	}

	keybindings, ok := keybindingsContainer["keybindings"].(map[string]interface{})
	if !ok {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Malformed configuration.  Aborting startup.")
	}

	keymap := map[rune]Keyevent{}
	// e.g., "h": "KEYCODE_HOME"
	for keycodeNewMapping, keycodeName := range keybindings {
		// Skip on empty value
		if len(keycodeNewMapping) == 0 {
			continue
		}
		keycodeName = strings.ToUpper(keycodeName.(string))

		// Is the specified keycodeName valid?  If not, KEYCODE_UNKNOWN.
		newKeyCode := Key(keycodeName.(string))

		// Extract the first rune from keycodeNewMapping.
		// That's the key we want to set.
		key, _ := utf8.DecodeRuneInString(keycodeNewMapping)

		keymap[key] = newKeyCode
	}

	// If no keys are defined, fail.
	if len(keymap) == 0 {
		log.Fatal("No keybindings are defined; aborting.")
	}

	return keymap
}
