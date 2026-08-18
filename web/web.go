// Package web provides the Runvil web framework: a routing and rendering
// layer built on the Go standard library's net/http and html/template.
//
// Pages are defined as routes, rendered through templates, and can be served
// live over HTTP or exported to a directory as a complete static website.
package web

import (
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Params holds the path variables extracted from a route pattern.
type Params map[string]string

// HandlerFunc handles an HTTP request with the route's parameters.
type HandlerFunc func(http.ResponseWriter, *http.Request, Params)

// Route pairs an HTTP method and path pattern with a handler.
type Route struct {
	method   string
	segments []segment
	handle   HandlerFunc
}

type segment struct {
	text     string
	param    string
	hasParam bool
}

// Router registers routes and dispatches requests.
type Router struct {
	routes  []*Route
	statics []staticRoute
	// NotFound, when set, handles unmatched routes.
	NotFound HandlerFunc
}

type staticRoute struct {
	prefix string
	dir    string
}

// NewRouter returns an empty Router.
func NewRouter() *Router {
	return &Router{}
}

// Get registers a handler for GET requests on a path pattern.
func (r *Router) Get(pattern string, h HandlerFunc) {
	r.Handle(http.MethodGet, pattern, h)
}

// Handle registers a handler for the given HTTP method and path pattern.
// Path segments of the form "{name}" become parameters.
func (r *Router) Handle(method, pattern string, h HandlerFunc) {
	if h == nil {
		h = func(http.ResponseWriter, *http.Request, Params) {}
	}
	r.routes = append(r.routes, &Route{method: method, segments: parsePattern(pattern), handle: h})
}

// Static serves files from dir under the URL prefix.
func (r *Router) Static(prefix, dir string) {
	r.statics = append(r.statics, staticRoute{prefix: prefix, dir: dir})
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, rt := range r.routes {
		if rt.method != req.Method {
			continue
		}
		if params, ok := rt.match(req.URL.Path); ok {
			rt.handle(w, req, params)
			return
		}
	}
	for _, st := range r.statics {
		if !strings.HasPrefix(req.URL.Path, st.prefix) {
			continue
		}
		rel := strings.TrimPrefix(req.URL.Path, st.prefix)
		file := filepath.Join(st.dir, filepath.Clean("/"+rel))
		http.ServeFile(w, req, file)
		return
	}
	if r.NotFound != nil {
		r.NotFound(w, req, nil)
		return
	}
	http.NotFound(w, req)
}

// match reports whether path matches the route pattern and returns extracted
// parameters.
func (rt *Route) match(path string) (Params, bool) {
	parts := splitPath(path)
	if len(parts) != len(rt.segments) {
		return nil, false
	}
	params := make(Params, len(rt.segments))
	for i, seg := range rt.segments {
		if seg.hasParam {
			if parts[i] == "" {
				return nil, false
			}
			params[seg.param] = parts[i]
			continue
		}
		if seg.text != parts[i] {
			return nil, false
		}
	}
	return params, true
}

func parsePattern(pattern string) []segment {
	segs := make([]segment, 0, 8)
	for _, part := range splitPath(pattern) {
		if len(part) > 2 && strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			segs = append(segs, segment{param: part[1 : len(part)-1], hasParam: true})
			continue
		}
		segs = append(segs, segment{text: part})
	}
	return segs
}

// splitPath splits a URL path into non-empty segments.
func splitPath(p string) []string {
	return strings.FieldsFunc(p, func(r rune) bool { return r == '/' })
}

// Templates executes named html/templates parsed from an fs.FS.
type Templates struct {
	set *template.Template
}

// ParseTemplates loads the templates matching pattern from fsys.
func ParseTemplates(fsys fs.FS, pattern string) (*Templates, error) {
	set, err := template.ParseFS(fsys, pattern)
	if err != nil {
		return nil, err
	}
	return &Templates{set: set}, nil
}

// Execute writes the named template with data to w.
func (t *Templates) Execute(w io.Writer, name string, data any) error {
	return t.set.ExecuteTemplate(w, name, data)
}

// Page is one concrete page of a static site.
type Page struct {
	// Path is the page URL, e.g. "/" or "/chapters/intro".
	Path string
	// Render writes the full page body.
	Render func(io.Writer) error
}

// Export writes pages and assets into outDir as a complete static website:
// the root page becomes outDir/index.html and every other page becomes
// outDir/<path>/index.html. Assets are written to their declared paths.
func Export(outDir string, pages []Page, assets map[string]string) error {
	sort.Slice(pages, func(i, j int) bool { return pages[i].Path < pages[j].Path })

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	for _, p := range pages {
		dir := outDir
		if rel := strings.Trim(p.Path, "/"); rel != "" {
			dir = filepath.Join(outDir, filepath.FromSlash(rel))
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(dir, "index.html"))
		if err != nil {
			return err
		}
		if err := p.Render(f); err != nil {
			f.Close()
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}

	for name, body := range assets {
		out := filepath.Join(outDir, filepath.FromSlash(cleanAssetPath(name)))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(out, []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// cleanAssetPath normalizes an asset URL into a safe relative path, refusing
// directory traversal.
func cleanAssetPath(name string) string {
	name = strings.TrimPrefix(name, "/")
	cleaned := path.Clean(name)
	if cleaned == "." || strings.HasPrefix(cleaned, "..") {
		return "assets/_"
	}
	return cleaned
}

// HTML writes an HTML response with the given status code.
func HTML(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	io.WriteString(w, body)
}

// Redirect responds with a redirect to the target URL.
func Redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url, code)
}
