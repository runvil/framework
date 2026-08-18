package ssg

import (
	"strings"
	"testing"
)

func TestScopeCSS(t *testing.T) {
	in := `
h1 { color: red; }
h1, .x a:hover { margin: 0; }
:root { --a: 1; }
@media (max-width: 48rem) { .m { padding: 0; } }
@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }
@font-face { font-family: X; src: url(x.woff2); }
.x::before { content: "a, b"; }
`
	got := scopeCSS(`[data-rv-component="c"]`, in)
	for _, want := range []string{
		`[data-rv-component="c"] h1{ color: red; }`,
		`[data-rv-component="c"] h1, [data-rv-component="c"] .x a:hover{ margin: 0; }`,
		`:root{ --a: 1; }`,
		`@media (max-width: 48rem) {[data-rv-component="c"] .m{ padding: 0; }}`,
		`@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }`,
		`@font-face { font-family: X; src: url(x.woff2); }`,
		`[data-rv-component="c"] .x::before{ content: "a, b"; }`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("scoped CSS missing %q", want)
		}
	}
	if strings.Contains(got, `[data-rv-component="c"] :root`) {
		t.Error(":root must not be scoped")
	}
}

func TestScopeFragmentInjectsRoot(t *testing.T) {
	got, err := scopeFragment("c", `<p>hi</p><p>bye</p>`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), `<p data-rv-component="c">hi</p>`) {
		t.Errorf("root element not scoped: %s", got)
	}
	if strings.Contains(string(got), `<p data-rv-component="c">bye</p>`) {
		t.Error("non-root elements must not carry the scope attribute")
	}
}

func TestScopeDocumentInjectsHtml(t *testing.T) {
	got, err := scopeDocument("base", "<!DOCTYPE html><html><head></head><body>x</body></html>")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), `<!DOCTYPE html>`) {
		t.Error("doctype lost")
	}
	if !strings.Contains(string(got), `<html data-rv-layout="base">`) {
		t.Errorf("html element not scoped: %s", got)
	}
}
