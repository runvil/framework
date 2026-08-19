package web

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/runvil/framework/ui"
	"github.com/runvil/framework/web/ssg"
)

// App assembles a monolith: it composes routes, middleware, pages, static
// assets, theming, and the shared component registry into one HTTP listener
// that serves HTML pages (SSR), JSON APIs (RVF-H3QD8), and embedded assets.
//
// Pages render through the same ssg render path used for static export, so
// live output and exported output are byte-identical for identical input.
type App struct {
	router *Router

	funcs     template.FuncMap
	comps     []ssg.Component
	layouts   []ssg.Layout
	assets    map[string]string
	registry  *ui.Registry
	theme     *ui.Theme
	emitProps bool

	pages   []*PageSpec
	statics []staticFS

	renderSite *ssg.Site
	dirty      bool
}

type staticFS struct {
	prefix string
	fsys   fs.FS
}

// PageSpec describes one page of an App: a root component rendered inside a
// layout. A page is static when Data is fixed, or dynamic when DataFunc
// computes per-request data.
type PageSpec struct {
	// Path is the served URL and export path ("/" for the root).
	Path string
	// Title is passed to the layout.
	Title string
	// Layout names the layout wrapping this page.
	Layout string
	// Root names the root component of the page, resolved through the registry.
	Root string
	// Data is the fixed page data for static pages.
	Data any
	// DataFunc, when set, computes the page data per request. Pages with a
	// DataFunc are dynamic and skipped by Export.
	DataFunc func(*http.Request) (any, error)
	// EmitProps serializes the page data to a data-props attribute on the
	// root element for client-side mounting. Off by default.
	EmitProps bool
	// Scripts are <script src> URLs injected before </body>, empty by default.
	Scripts []string
}

// NewApp returns an App pre-loaded with the default component registry and a
// trusted-HTML template helper.
func NewApp() *App {
	a := &App{
		router:   NewRouter(),
		assets:   map[string]string{},
		registry: ui.Default(),
		funcs:    template.FuncMap{},
	}
	a.funcs = template.FuncMap{
		"html": func(v any) template.HTML {
			return template.HTML(fmt.Sprint(v))
		},
		"themeHead": func() template.HTML {
			if a.theme == nil {
				return template.HTML("")
			}
			return a.theme.Script()
		},
		"themeButton": func() template.HTML {
			if a.theme == nil {
				return template.HTML("")
			}
			return a.theme.Button()
		},
	}
	a.dirty = true
	return a
}

// Router exposes the underlying router for API routes and middleware.
func (a *App) Router() *Router { return a.router }

// Use registers middleware applied to every route.
func (a *App) Use(mw ...Middleware) { a.router.Use(mw...) }

// Method registers an API or HTML handler on the router.
func (a *App) Method(method, pattern string, h HandlerFunc) {
	a.router.Handle(method, pattern, h)
}

// Funcs registers template functions available to components and layouts.
func (a *App) Funcs(f template.FuncMap) *App {
	for name, fn := range f {
		if a.funcs == nil {
			a.funcs = template.FuncMap{}
		}
		a.funcs[name] = fn
	}
	a.dirty = true
	return a
}

// Component registers a template component by name.
func (a *App) Component(c ssg.Component) *App {
	a.comps = append(a.comps, c)
	a.dirty = true
	return a
}

// Layout registers a layout by name.
func (a *App) Layout(l ssg.Layout) *App {
	a.layouts = append(a.layouts, l)
	a.dirty = true
	return a
}

// Asset registers a static file by its output path (e.g. "favicon.ico").
func (a *App) Asset(name, body string) *App {
	a.assets[name] = body
	a.dirty = true
	return a
}

// Theme attaches the theming system; its CSS is collected into the site's
// assets.
func (a *App) Theme(t *ui.Theme) *App {
	a.theme = t
	a.dirty = true
	return a
}

// Registry overrides the component registry used to resolve page roots.
func (a *App) Registry(r *ui.Registry) *App {
	a.registry = r
	a.dirty = true
	return a
}

// EmitProps globally enables data-props emission for pages that request it.
func (a *App) EmitProps() *App {
	a.emitProps = true
	a.dirty = true
	return a
}

// Page registers a page.
func (a *App) Page(p PageSpec) *App {
	a.pages = append(a.pages, &p)
	a.dirty = true
	return a
}

// Static mounts an embed-compatible fs.FS (e.g. an //go:embed value) under
// prefix, serving its files untouched.
func (a *App) Static(prefix string, fsys fs.FS) *App {
	a.statics = append(a.statics, staticFS{prefix: prefix, fsys: fsys})
	return a
}

// Handler returns the assembled http.Handler: routes, middleware, pages, and
// static mounts. It is safe for tests and embedding in larger servers.
func (a *App) Handler() http.Handler {
	site := a.site()
	for _, p := range a.pages {
		spec := p
		a.router.Get(spec.Path, func(w http.ResponseWriter, req *http.Request, _ Params) {
			data, err := pageData(spec, req)
			if err != nil {
				Error(w, err)
				return
			}
			body, err := site.RenderPage(&ssg.Page{
				Path:    spec.Path,
				Title:   spec.Title,
				Layout:  spec.Layout,
				Root:    spec.Root,
				Data:    data,
				Props:   propsFor(spec, data),
				Scripts: spec.Scripts,
			})
			if err != nil {
				Error(w, err)
				return
			}
			HTML(w, http.StatusOK, body)
		})
	}
	for _, name := range site.Assets() {
		assetName := name
		assetBody, _ := site.AssetBody(name)
		a.router.Get("/"+assetName, func(w http.ResponseWriter, _ *http.Request, _ Params) {
			serveAsset(w, assetName, assetBody)
		})
	}
	for _, st := range a.statics {
		a.router.StaticFS(st.prefix, st.fsys)
	}
	return a.router
}

// serveAsset writes an embedded asset body with a content type derived from
// its extension.
func serveAsset(w http.ResponseWriter, name, body string) {
	switch filepath.Ext(name) {
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	_, _ = io.WriteString(w, body)
}

// Export writes the App's static pages as a complete website into outDir,
// byte-identical to the ssg build for the same input. Dynamic pages (those
// with a DataFunc) are skipped.
func (a *App) Export(outDir string) ([]string, error) {
	s := a.exportSite()
	return s.Build(outDir)
}

// Run serves the App over HTTP on addr, shutting down gracefully on
// SIGINT/SIGTERM.
func (a *App) Run(addr string) error {
	srv := &http.Server{Addr: addr, Handler: a.Handler()}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("runvil server started", "addr", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		_ = srv.Close()
		return err
	}
	slog.Info("runvil server stopped")
	return nil
}

// site builds the shared render site, cached until the App is mutated.
func (a *App) site() *ssg.Site {
	if a.renderSite != nil && !a.dirty {
		return a.renderSite
	}
	a.renderSite = a.buildSite(nil)
	a.dirty = false
	return a.renderSite
}

// exportSite builds a site for static export: dynamic pages are skipped.
func (a *App) exportSite() *ssg.Site {
	var static []PageSpec
	for _, p := range a.pages {
		if p.DataFunc == nil {
			static = append(static, *p)
		}
	}
	return a.buildSite(static)
}

func (a *App) buildSite(pages []PageSpec) *ssg.Site {
	s := ssg.New().Funcs(a.funcs).Registry(a.registry)
	if emit := a.emitEnabled(); emit {
		s.EmitProps()
	}
	for _, c := range a.comps {
		s.Component(c)
	}
	for _, l := range a.layouts {
		s.Layout(l)
	}
	for name, body := range a.assets {
		s.Asset(name, body)
	}
	if a.theme != nil {
		s.Asset("assets/theme.css", ui.ThemeModeVarsCSS+"\n"+ui.ThemeToggleCSS)
	}
	if a.registry != nil {
		s.Asset("assets/ui.css", ui.ComponentsCSS())
	}
	for _, p := range pages {
		s.Page(ssg.Page{
			Path:    p.Path,
			Title:   p.Title,
			Layout:  p.Layout,
			Root:    p.Root,
			Data:    p.Data,
			Props:   propsFor(&p, p.Data),
			Scripts: p.Scripts,
		})
	}
	return s
}

// emitEnabled reports whether any page requests data-props emission, either
// globally at the App level or per page.
func (a *App) emitEnabled() bool {
	if a.emitProps {
		return true
	}
	for _, p := range a.pages {
		if p.EmitProps {
			return true
		}
	}
	return false
}

func pageData(spec *PageSpec, req *http.Request) (any, error) {
	if spec.DataFunc == nil {
		return spec.Data, nil
	}
	return spec.DataFunc(req)
}

func propsFor(spec *PageSpec, data any) any {
	if !spec.EmitProps {
		return nil
	}
	return data
}
