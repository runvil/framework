// Package ssg provides an Astro/Svelte-inspired static site generator built
// on html/template: sites are composed from named components with scoped
// styles, wrapped by layouts with a content slot, and rendered to static HTML
// plus a single collected stylesheet in one build step.
package ssg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/runvil/framework/ui"
)

// Component is a reusable piece of UI. Body is the html/template source;
// Style is CSS scoped to the component's rendered root element.
type Component struct {
	Name  string
	Body  string
	Style string
}

// Layout wraps page content in a shared shell. Body receives a LayoutData
// value (fields Title, Content, Data). Style is scoped to the layout's root
// element.
type Layout struct {
	Name  string
	Body  string
	Style string
}

// LayoutData is the value passed to a layout template.
type LayoutData struct {
	// Title is the page title.
	Title string
	// Content is the rendered page component HTML.
	Content template.HTML
	// Data is the page's data value.
	Data any
}

// Page is one output page of the site.
type Page struct {
	// Path is the output URL: "/" for the root, "/about" for a clean URL, or
	// "/about.html" for an explicit filename.
	Path string
	// Title is passed to the layout.
	Title string
	// Layout names the Layout wrapping this page.
	Layout string
	// Root names the root component of the page.
	Root string
	// Data is passed to the root component.
	Data any
	// Props, when the Site emits props, is serialized to a data-props
	// attribute on the page's root element for client-side mounting.
	Props any
	// Scripts are <script src> URLs injected before the closing body tag,
	// empty by default.
	Scripts []string
}

// Site builds a static website from components, layouts, pages, and assets.
type Site struct {
	funcs     template.FuncMap
	layouts   map[string]*Layout
	comps     map[string]*Component
	pages     []*Page
	assets    map[string]string
	reg       *ui.Registry
	emitProps bool

	set *template.Template
	mu  sync.Mutex
	// used records every component and layout whose styles were referenced
	// during a build.
	used map[string]bool
}

// New returns an empty Site.
func New() *Site {
	return &Site{
		layouts: map[string]*Layout{},
		comps:   map[string]*Component{},
		assets:  map[string]string{},
		used:    map[string]bool{},
	}
}

// Funcs registers template functions available to components and layouts.
func (s *Site) Funcs(f template.FuncMap) *Site {
	if s.funcs == nil {
		s.funcs = template.FuncMap{}
	}
	for name, fn := range f {
		s.funcs[name] = fn
	}
	return s
}

// Component registers a component by name.
func (s *Site) Component(c Component) *Site {
	s.comps[c.Name] = &c
	return s
}

// Layout registers a layout by name.
func (s *Site) Layout(l Layout) *Site {
	s.layouts[l.Name] = &l
	return s
}

// Page registers an output page.
func (s *Site) Page(p Page) *Site {
	s.pages = append(s.pages, &p)
	return s
}

// Asset registers a static file by its output path, e.g. "favicon.ico" or
// "assets/theme.css".
func (s *Site) Asset(name, body string) *Site {
	s.assets[name] = body
	return s
}

// Assets returns the names of all registered assets in sorted order.
func (s *Site) Assets() []string {
	names := make([]string, 0, len(s.assets))
	for name := range s.assets {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// AssetBody returns the content of a registered asset.
func (s *Site) AssetBody(name string) (string, bool) {
	body, ok := s.assets[name]
	return body, ok
}

// Open implements fs.FS over the site's asset map, letting an App serve
// assets from a live Site through a single static mount.
func (s *Site) Open(name string) (fs.File, error) {
	name = strings.TrimPrefix(name, "/")
	body, ok := s.assets[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return fsFile{name: name, r: strings.NewReader(body)}, nil
}

// fsFile is an fs.File backed by a string asset body.
type fsFile struct {
	name string
	r    *strings.Reader
}

func (f fsFile) Stat() (fs.FileInfo, error) {
	return fs.FileInfo(fi{name: f.name, size: int64(f.r.Len())}), nil
}
func (f fsFile) Read(p []byte) (int, error) { return f.r.Read(p) }
func (f fsFile) Close() error               { return nil }

type fi struct {
	name string
	size int64
}

func (i fi) Name() string       { return path.Base(i.name) }
func (i fi) Size() int64        { return i.size }
func (i fi) Mode() fs.FileMode  { return 0o444 }
func (i fi) ModTime() time.Time { return time.Time{} }
func (i fi) IsDir() bool        { return false }
func (i fi) Sys() any           { return nil }

// Build renders every page and writes the site into outDir. It returns the
// paths of the files created, relative to outDir.
func (s *Site) Build(outDir string) ([]string, error) {
	if err := s.prepare(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.used = map[string]bool{}
	s.mu.Unlock()

	var created []string
	for _, p := range s.pages {
		body, err := s.renderPage(p)
		if err != nil {
			return nil, err
		}
		rel := pageRelPath(p.Path)
		out := filepath.Join(outDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(out, []byte(body), 0o644); err != nil {
			return nil, err
		}
		created = append(created, rel)
	}

	if css := s.collectedCSS(); css != "" {
		rel := "assets/style.css"
		out := filepath.Join(outDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(out, []byte(css), 0o644); err != nil {
			return nil, err
		}
		created = append(created, rel)
	}

	names := make([]string, 0, len(s.assets))
	for name := range s.assets {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		rel := cleanAssetPath(name)
		out := filepath.Join(outDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(out, []byte(s.assets[name]), 0o644); err != nil {
			return nil, err
		}
		created = append(created, rel)
	}
	return created, nil
}

// Handler returns an HTTP handler serving the site's pages and assets,
// rendered on demand for development.
func (s *Site) Handler() http.Handler {
	if err := s.prepare(); err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		})
	}
	mux := http.NewServeMux()
	for _, p := range s.pages {
		pp := p
		seen := map[string]bool{}
		paths := []string{p.Path}
		if !strings.HasSuffix(p.Path, ".html") {
			clean := strings.TrimSuffix(p.Path, "/")
			if slash := clean + "/"; !seen[slash] {
				paths = append(paths, slash)
			}
		}
		for _, rp := range paths {
			if seen[rp] {
				continue
			}
			seen[rp] = true
			mux.HandleFunc(rp, func(w http.ResponseWriter, r *http.Request) {
				body, err := s.renderPage(pp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = io.WriteString(w, body)
			})
		}
	}
	for name, body := range s.assets {
		assetName, assetBody := name, body
		mux.HandleFunc("/"+assetName, func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, assetName, time.Time{}, strings.NewReader(assetBody))
		})
	}
	return mux
}

// Registry sets the component registry used to resolve component names that
// don't correspond to template bodies. Resolution order: registered Go-native
// components first, then template components defined on the site.
func (s *Site) Registry(reg *ui.Registry) *Site {
	s.reg = reg
	return s
}

// EmitProps enables serializing a page's Props to a data-props attribute on
// the page's root element, providing a stable mount point for client-side
// frameworks. Off by default; pages with a nil Props are left untouched.
func (s *Site) EmitProps() *Site {
	s.emitProps = true
	return s
}

// RenderPage renders a page's root component inside its layout as a complete
// HTML document. It is the shared render path used by both live SSR and
// static export, guaranteeing identical output for identical input.
func (s *Site) RenderPage(p *Page) (string, error) {
	if err := s.prepare(); err != nil {
		return "", err
	}
	return s.renderPage(p)
}

// renderComponent renders a component with data, marks it used, and injects
// its scope attribute onto the rendered root element.
func (s *Site) renderComponent(name string, data any) (template.HTML, error) {
	if _, ok := s.comps[name]; ok {
		s.markUsed(name)
		var buf bytes.Buffer
		if err := s.set.ExecuteTemplate(&buf, name, data); err != nil {
			return "", err
		}
		return scopeFragment(name, buf.String())
	}
	if s.reg != nil {
		if out, err := s.reg.Render(name, data); err == nil {
			return out, nil
		} else if !strings.HasPrefix(err.Error(), "ui: undefined component") {
			return "", err
		}
	}
	return "", fmt.Errorf("ssg: undefined component %q", name)
}

// renderPage renders a page's root component inside its layout.
func (s *Site) renderPage(p *Page) (string, error) {
	root, err := s.renderComponent(p.Root, p.Data)
	if err != nil {
		return "", err
	}
	if s.emitProps && p.Props != nil {
		root, err = injectDataProps(root, p.Props)
		if err != nil {
			return "", err
		}
	}
	_, ok := s.layouts[p.Layout]
	if !ok {
		return "", fmt.Errorf("ssg: undefined layout %q", p.Layout)
	}
	s.markUsed(p.Layout)
	var buf bytes.Buffer
	if err := s.set.ExecuteTemplate(&buf, p.Layout, LayoutData{
		Title:   p.Title,
		Content: root,
		Data:    p.Data,
	}); err != nil {
		return "", err
	}
	// Inject the layout scope attribute onto the document root so layout
	// styles can be scoped under [data-rv-layout="name"].
	scoped, err := scopeDocument(p.Layout, buf.String())
	if err != nil {
		return "", err
	}
	body := string(scoped)
	if len(p.Scripts) > 0 {
		body = injectScripts(body, p.Scripts)
	}
	return body, nil
}

// markUsed records that a component or layout was referenced.
func (s *Site) markUsed(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.used[name] = true
}

// collectedCSS concatenates the scoped styles of every component and layout
// referenced during the last build, in deterministic order.
func (s *Site) collectedCSS() string {
	var out strings.Builder
	names := make([]string, 0, len(s.used))
	for name := range s.used {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if c, ok := s.comps[name]; ok && strings.TrimSpace(c.Style) != "" {
			out.WriteString("/* " + name + " */\n")
			out.WriteString(scopeCSS(`[data-rv-component="`+name+`"]`, c.Style, true))
			out.WriteString("\n")
			continue
		}
		if l, ok := s.layouts[name]; ok && strings.TrimSpace(l.Style) != "" {
			out.WriteString("/* " + name + " */\n")
			out.WriteString(scopeCSS(`[data-rv-layout="`+name+`"]`, l.Style, false))
			out.WriteString("\n")
		}
	}
	return out.String()
}

// prepare parses component and layout templates into one template set.
func (s *Site) prepare() error {
	if s.set != nil {
		return nil
	}
	funcs := template.FuncMap{}
	for name, fn := range s.funcs {
		funcs[name] = fn
	}
	funcs["component"] = s.renderComponent
	set := template.New("").Funcs(funcs)
	compNames := make([]string, 0, len(s.comps))
	for name := range s.comps {
		compNames = append(compNames, name)
	}
	sort.Strings(compNames)
	for _, name := range compNames {
		if _, err := set.New(name).Parse(s.comps[name].Body); err != nil {
			return fmt.Errorf("ssg: component %q: %w", name, err)
		}
	}
	layoutNames := make([]string, 0, len(s.layouts))
	for name := range s.layouts {
		layoutNames = append(layoutNames, name)
	}
	sort.Strings(layoutNames)
	for _, name := range layoutNames {
		if _, err := set.New(name).Parse(s.layouts[name].Body); err != nil {
			return fmt.Errorf("ssg: layout %q: %w", name, err)
		}
	}
	s.set = set
	return nil
}

// pageRelPath maps a page path to a filesystem path: "/" -> "index.html",
// "/about" -> "about/index.html", "/about.html" -> "about.html".
func pageRelPath(p string) string {
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return "index.html"
	}
	if strings.HasSuffix(p, ".html") {
		return p
	}
	return filepath.Join(filepath.FromSlash(p), "index.html")
}

// cleanAssetPath normalizes an asset path, refusing directory traversal.
func cleanAssetPath(name string) string {
	name = strings.TrimPrefix(name, "/")
	cleaned := path.Clean(name)
	if cleaned == "." || strings.HasPrefix(cleaned, "..") {
		return "assets/_"
	}
	return cleaned
}

// scopeFragment injects the component scope attribute onto the first element
// of an HTML fragment and re-renders the fragment.
func scopeFragment(name, body string) (template.HTML, error) {
	frag, err := html.ParseFragment(strings.NewReader(body), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return "", err
	}
	target := firstElement(frag)
	if target == nil {
		return template.HTML(body), nil
	}
	for _, a := range target.Attr {
		if a.Key == "data-rv-component" {
			// Root already carries a scope attribute from a nested
			// component; this wrapper has no markup of its own.
			return template.HTML(body), nil
		}
	}
	target.Attr = append(target.Attr, html.Attribute{Key: "data-rv-component", Val: name})
	var buf bytes.Buffer
	for _, n := range frag {
		if err := html.Render(&buf, n); err != nil {
			return "", err
		}
	}
	return template.HTML(buf.String()), nil
}

// scopeDocument injects the layout scope attribute onto the <html> element of
// a full document and re-renders it, preserving the doctype.
func scopeDocument(name, doc string) (template.HTML, error) {
	root, err := html.Parse(strings.NewReader(doc))
	if err != nil {
		return "", err
	}
	for n := root.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode && n.Data == "html" {
			n.Attr = append(n.Attr, html.Attribute{Key: "data-rv-layout", Val: name})
			break
		}
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, root); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// firstElement returns the first element node of a parsed fragment.
func firstElement(frag []*html.Node) *html.Node {
	for _, n := range frag {
		if n.Type == html.ElementNode {
			return n
		}
	}
	return nil
}

// injectDataProps serializes props to a data-props attribute on the first
// element of a rendered fragment, giving client-side frameworks a stable,
// opt-in hydration payload.
func injectDataProps(fragHTML template.HTML, props any) (template.HTML, error) {
	data, err := json.Marshal(props)
	if err != nil {
		return "", fmt.Errorf("ssg: encode page props: %w", err)
	}
	frag, err := html.ParseFragment(strings.NewReader(string(fragHTML)), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return "", err
	}
	target := firstElement(frag)
	if target == nil {
		return fragHTML, nil
	}
	target.Attr = append(target.Attr, html.Attribute{Key: "data-props", Val: string(data)})
	var buf bytes.Buffer
	for _, n := range frag {
		if err := html.Render(&buf, n); err != nil {
			return "", err
		}
	}
	return template.HTML(buf.String()), nil
}

// injectScripts inserts <script src> tags before the closing </body> tag,
// the documented mount script slot, which is empty by default.
func injectScripts(body string, scripts []string) string {
	var sb strings.Builder
	for _, src := range scripts {
		sb.WriteString(`<script src="`)
		sb.WriteString(src)
		sb.WriteString(`"></script>`)
	}
	sb.WriteString("</body>")
	idx := strings.LastIndex(body, "</body>")
	if idx < 0 {
		return body + sb.String()
	}
	return body[:idx] + sb.String()
}
