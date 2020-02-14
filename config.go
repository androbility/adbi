package adbi

import (
	"strings"
	"unicode/utf8"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfigFile attempts to read a config-file from configDir.
// Failing that, LoadConfigFile parses defaultBindings for the
// configuration.
func LoadConfigFile(configDir, defaultBindings string) map[rune]Keyevent {
	keymap := map[rune]Keyevent{}

	viper.SetConfigName("config")
	viper.AddConfigPath(configDir)
	viper.SetDefault("keybindings", keymap)

	if err := viper.ReadInConfig(); err != nil {
		log.Debug("Configuration file not found; using defaults.")

		viper.Reset()
		viper.SetConfigType("json")
		r := strings.NewReader(defaultBindings)
		if err = viper.ReadConfig(r); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Debug("Loading default configuration failed.  Please contact the developer.")
		}
	}

	// e.g., "h": "KEYCODE_HOME"
	for keycodeNewMapping, keycodeName := range viper.GetStringMapString("keybindings") {
		// Skip on empty value
		if len(keycodeNewMapping) == 0 {
			continue
		}
		keycodeName = strings.ToUpper(keycodeName)

		// Is the specified keycodeName valid?  If not, KEYCODE_UNKNOWN.
		newKeyCode := Key(keycodeName)

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
