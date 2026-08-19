// Package theme defines the application palette and toggle styling
// (FRK-STR-006). Palette default: neon green light/dark, as the ecosystem
// default.
package theme

import "github.com/runvil/framework/ui"

// New returns the application theme with the default Runvil palette.
func New() *ui.Theme {
	return &ui.Theme{
		Light: ui.Palette{Primary: "#00c853"},
		Dark:  ui.Palette{Primary: "#5cff8e"},
	}
}
