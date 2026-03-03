package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"journal-entry/internal/account"
	"journal-entry/internal/journal"
	"journal-entry/internal/report"
	"journal-entry/internal/shared/middleware"
)

// NewRouter creates and configures the chi router with all routes.
func NewRouter(templates map[string]*template.Template, accountHandler *account.Handler, journalHandler *journal.Handler, reportHandler *report.Handler) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)

	// Serve static files
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "static")
	fileServer := http.FileServer(http.Dir(staticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// === Routes ===

	// Dashboard (home)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, ok := templates["dashboard/index.html"]
		if !ok {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "base", nil); err != nil {
			log.Printf("[RENDER ERROR] dashboard: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Account routes
	accountHandler.RegisterRoutes(r)

	// Journal routes
	journalHandler.RegisterRoutes(r)

	// Report routes
	reportHandler.RegisterRoutes(r)

	return r
}

// ParseTemplates parses all page templates, each composed with the base layout.
// Returns a map of template name → *template.Template.
//
// Each page template is parsed together with the base layout + components,
// so {{ template "content" }} in base.html gets filled by the page template.
func ParseTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)

	// Base layout + components (shared across all pages)
	layoutFiles := []string{
		"templates/layout/base.html",
	}

	// Find all component files
	componentFiles, _ := filepath.Glob("templates/components/*.html")
	sharedFiles := append(layoutFiles, componentFiles...)

	// Page templates — each one is a full page that defines "content" block
	pageTemplates := []string{
		"templates/dashboard/index.html",
		"templates/account/list.html",
		"templates/account/form.html",
		"templates/journal/list.html",
		"templates/journal/form.html",
		"templates/report/ledger.html",
		"templates/report/trial_balance.html",
	}

	funcMap := template.FuncMap{
		"activeNav": func(current, page string) string {
			if current == page {
				return "l-sidebar__link--active"
			}
			return ""
		},
		// fmtNum formats a float64 as an integer with dot thousand separators (Indonesian format)
		"fmtNum": func(n float64) string {
			isNeg := n < 0
			if isNeg {
				n = -n
			}
			intPart := int64(n + 0.5)
			s := fmt.Sprintf("%d", intPart)
			// Add dots every 3 digits from right
			out := make([]byte, 0, len(s)+len(s)/3)
			for i, c := range s {
				if i > 0 && (len(s)-i)%3 == 0 {
					out = append(out, '.')
				}
				out = append(out, byte(c))
			}
			if isNeg {
				return "-" + string(out)
			}
			return string(out)
		},
		// gtf compares float64 > 0
		"gtf": func(n float64) bool {
			return n > 0
		},
	}

	for _, page := range pageTemplates {
		// Parse: base layout + components + this page template
		files := append([]string{}, sharedFiles...)
		files = append(files, page)

		name := page[len("templates/"):]

		tmpl, err := template.New("").Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}

		templates[name] = tmpl
	}

	log.Printf("📄 Parsed %d page templates", len(templates))
	return templates
}
