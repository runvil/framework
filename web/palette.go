// Palette is a named color schema for a single mode (light or dark).
package web

import (
	"reflect"
	"strings"
)

// Color is a CSS color value, e.g. "#0b58a1", "rgb(11 88 161)", or
// "oklch(0.5 0.2 250)".
type Color string

// Palette is a named color schema for a single mode (light or dark). Every
// base color has a -content companion: the color of text and icons rendered
// on top of it. Empty fields fall back to the mode's default palette, so a
// site can override just a few colors.
type Palette struct {
	Base1            Color `css:"base-1"`
	Base1Content     Color `css:"base-1-content"`
	Base2            Color `css:"base-2"`
	Base2Content     Color `css:"base-2-content"`
	Base3            Color `css:"base-3"`
	Base3Content     Color `css:"base-3-content"`
	Primary          Color `css:"primary"`
	PrimaryContent   Color `css:"primary-content"`
	Secondary        Color `css:"secondary"`
	SecondaryContent Color `css:"secondary-content"`
	Accent           Color `css:"accent"`
	AccentContent    Color `css:"accent-content"`
	Ghost            Color `css:"ghost"`
	GhostContent     Color `css:"ghost-content"`
	Neutral          Color `css:"neutral"`
	NeutralContent   Color `css:"neutral-content"`
	Success          Color `css:"success"`
	SuccessContent   Color `css:"success-content"`
	Info             Color `css:"info"`
	InfoContent      Color `css:"info-content"`
	Warning          Color `css:"warning"`
	WarningContent   Color `css:"warning-content"`
	Error            Color `css:"error"`
	ErrorContent     Color `css:"error-content"`
}

// DefaultLightPalette is the palette applied when no light override is set.
var DefaultLightPalette = Palette{
	Base1:            "#fdfdfb",
	Base1Content:     "#1a1a1a",
	Base2:            "#f5f2eb",
	Base2Content:     "#666",
	Base3:            "#efece4",
	Base3Content:     "#4a4a4a",
	Primary:          "#0b58a1",
	PrimaryContent:   "#ffffff",
	Secondary:        "#6b7280",
	SecondaryContent: "#ffffff",
	Accent:           "#b45309",
	AccentContent:    "#ffffff",
	Ghost:            "#e9e6dc",
	GhostContent:     "#666",
	Neutral:          "#e3e0d8",
	NeutralContent:   "#666",
	Success:          "#15803d",
	SuccessContent:   "#ffffff",
	Info:             "#0369a1",
	InfoContent:      "#ffffff",
	Warning:          "#b45309",
	WarningContent:   "#ffffff",
	Error:            "#b91c1c",
	ErrorContent:     "#ffffff",
}

// DefaultDarkPalette is the palette applied when no dark override is set.
var DefaultDarkPalette = Palette{
	Base1:            "#17171a",
	Base1Content:     "#e6e4dd",
	Base2:            "#202024",
	Base2Content:     "#9a968b",
	Base3:            "#2a2a30",
	Base3Content:     "#c6c2b8",
	Primary:          "#7ab4ef",
	PrimaryContent:   "#0b0e14",
	Secondary:        "#8b8f9c",
	SecondaryContent: "#0e0e11",
	Accent:           "#f0b429",
	AccentContent:    "#241a00",
	Ghost:            "#26262b",
	GhostContent:     "#9a968b",
	Neutral:          "#2e2e34",
	NeutralContent:   "#9a968b",
	Success:          "#4ade80",
	SuccessContent:   "#0b1f10",
	Info:             "#60a5fa",
	InfoContent:      "#0a1628",
	Warning:          "#fbbf24",
	WarningContent:   "#241a00",
	Error:            "#f87171",
	ErrorContent:     "#260a0a",
}

// cssVars emits "name: value;" declarations for every set color, filling in
// non-empty fields from overrides and empty fields from the defaults palette.
func (p Palette) cssVars(defaults Palette) string {
	var b strings.Builder
	rv, rdd := reflect.ValueOf(p), reflect.ValueOf(defaults)
	for i := 0; i < rv.NumField(); i++ {
		v := rv.Field(i).String()
		if v == "" {
			v = rdd.Field(i).String()
		}
		if v == "" {
			continue
		}
		b.WriteString("--")
		b.WriteString(rv.Type().Field(i).Tag.Get("css"))
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString(";")
	}
	return b.String()
}
