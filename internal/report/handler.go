package report

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"journal-entry/internal/account"
)

// Handler handles HTTP requests for reports.
type Handler struct {
	svc        Service
	accountSvc account.Service
	templates  map[string]*template.Template
}

// NewHandler creates a new report Handler.
func NewHandler(svc Service, accountSvc account.Service, templates map[string]*template.Template) *Handler {
	return &Handler{svc: svc, accountSvc: accountSvc, templates: templates}
}

// RegisterRoutes registers report routes.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		r.Get("/ledger", h.HandleLedger)
		r.Get("/trial-balance", h.HandleTrialBalance)
	})
}

// HandleLedger renders the general ledger report.
func (h *Handler) HandleLedger(w http.ResponseWriter, r *http.Request) {
	// Get all accounts for dropdown
	accounts, err := h.accountSvc.ListAccounts(r.Context())
	if err != nil {
		log.Printf("[ERROR] list accounts for ledger: %v", err)
		http.Error(w, "Gagal memuat daftar akun", http.StatusInternalServerError)
		return
	}

	// Parse filter params
	accountIDStr := r.URL.Query().Get("account_id")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	data := map[string]any{
		"Accounts":  accounts,
		"AccountID": accountIDStr,
		"DateFrom":  dateFrom,
		"DateTo":    dateTo,
		"Report":    nil,
	}

	// If account selected, get ledger
	if accountIDStr != "" {
		accountID, err := strconv.Atoi(accountIDStr)
		if err == nil {
			report, err := h.svc.GetLedgerReport(r.Context(), accountID, dateFrom, dateTo)
			if err != nil {
				log.Printf("[ERROR] get ledger: %v", err)
			} else {
				data["Report"] = report
				data["DateFrom"] = report.DateFrom
				data["DateTo"] = report.DateTo
			}
		}
	}

	h.render(w, "report/ledger.html", data)
}

// HandleTrialBalance renders the trial balance report.
func (h *Handler) HandleTrialBalance(w http.ResponseWriter, r *http.Request) {
	tb, err := h.svc.GetTrialBalance(r.Context())
	if err != nil {
		log.Printf("[ERROR] get trial balance: %v", err)
		http.Error(w, "Gagal memuat neraca saldo", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"TrialBalance": tb,
		"TypeLabels":   buildTypeLabels(),
	}

	h.render(w, "report/trial_balance.html", data)
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("[RENDER ERROR] %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func buildTypeLabels() map[string]string {
	return map[string]string{
		"asset":     "Aset",
		"liability": "Kewajiban",
		"equity":    "Ekuitas",
		"revenue":   "Pendapatan",
		"cogs":      "Harga Pokok",
		"expense":   "Beban",
	}
}
