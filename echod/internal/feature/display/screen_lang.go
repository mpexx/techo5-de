//go:build !dot

package display

import "github.com/HuskerMinion/techo5/echod/internal/config"

// screenLang keeps the German fork's date and weather in German when the screen listens for all
// languages. A specific screen language still controls those words as in upstream TECHO5.
func screenLang() string {
	if lang := config.Get().Screen.Language; lang != "" {
		return lang
	}
	return "de"
}
