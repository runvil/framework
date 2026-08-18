// Package web provides the Runvil web framework: a routing and rendering
// layer built on the Go standard library's net/http and html/template.
//
// Pages are defined as routes, rendered through templates, and can be served
// live over HTTP or exported to a directory as a complete static website.
package web

import (
	"html/template"
	"strconv"
	"strings"
)

// defaultThemeKey is the localStorage key used when Theme.StorageKey is empty.
const defaultThemeKey = "runvil-theme"

// Theme configures the light/dark theming system for a site. It follows the
// system color scheme by default, persists the user's choice, and switches
// between light and dark without server involvement.
type Theme struct {
	// StorageKey is the localStorage key persisting the user's choice.
	// Defaults to "runvil-theme".
	StorageKey string
	// Default is the initial preference when nothing is stored:
	// "auto", "light", or "dark". Defaults to "auto".
	Default string
}

func (t Theme) storageKey() string {
	if t.StorageKey == "" {
		return defaultThemeKey
	}
	return t.StorageKey
}

func (t Theme) defaultTheme() string {
	switch t.Default {
	case "light", "dark":
		return t.Default
	default:
		return "auto"
	}
}

// Script returns the inline <script> block for the document head. It applies
// the stored or system theme before first paint (avoiding a flash of the wrong
// theme), keeps the color-scheme meta in sync, and wires any element carrying
// the data-theme-toggle attribute to cycle light and dark modes.
func (t Theme) Script() template.HTML {
	return template.HTML("<script>" + themeJS(t.storageKey(), t.defaultTheme()) + "</script>")
}

// Button returns the theme-switcher button markup: a sun/moon toggle that the
// theme script wires automatically. Place it anywhere in the page.
func (t Theme) Button() template.HTML {
	return template.HTML(`<button type="button" class="theme-toggle" data-theme-toggle aria-label="Toggle light/dark theme" title="Toggle theme"><svg class="icon-sun" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg><svg class="icon-moon" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg></button>`)
}

// themeJS generates the plain-JavaScript theming runtime.
func themeJS(key, def string) string {
	var b strings.Builder
	b.WriteString(`(function(){var k=` + strconv.Quote(key) + `,d=` + strconv.Quote(def) + `,t=d;try{t=localStorage.getItem(k)||d}catch(e){}var m=window.matchMedia('(prefers-color-scheme: dark)');function apply(){var dark=t==='dark'||(t==='auto'&&m.matches);var r=document.documentElement;r.dataset.theme=dark?'dark':'light';r.style.colorScheme=dark?'dark':'light';var meta=document.querySelector('meta[name="color-scheme"]');if(meta)meta.content=dark?'dark':'light'}apply();if(m.addEventListener){m.addEventListener('change',function(){if(t==='auto')apply()})}window.runvilTheme={get:function(){return t},set:function(x){t=x;try{localStorage.setItem(k,x)}catch(e){}apply()},toggle:function(){var dark=t==='dark'||(t==='auto'&&m.matches);window.runvilTheme.set(dark?'light':'dark')}};document.addEventListener('click',function(e){var el=e.target&&e.target.closest?e.target.closest('[data-theme-toggle]'):null;if(el)window.runvilTheme.toggle()})})();`)
	return b.String()
}
