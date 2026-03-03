package dashboard

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"journal-entry/internal/journal"
	"journal-entry/internal/report"
)

// Summary holds aggregated financial data for the dashboard.
type Summary struct {
	TotalAsset     float64
	TotalLiability float64
	TotalRevenue   float64
	TotalExpense   float64
	TotalCOGS      float64
	TotalEquity    float64
	NetIncome      float64 // Revenue - Expense - COGS
}

// Handler handles dashboard HTTP requests.
type Handler struct {
	reportSvc  report.Service
	journalSvc journal.Service
	templates  map[string]*template.Template
}

// NewHandler creates a new dashboard Handler.
func NewHandler(reportSvc report.Service, journalSvc journal.Service, templates map[string]*template.Template) *Handler {
	return &Handler{reportSvc: reportSvc, journalSvc: journalSvc, templates: templates}
}

// HandleDashboard renders the dashboard page.
func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get financial summary from trial balance
	summary := h.getSummary(ctx)

	// Get recent journal entries
	entries, err := h.journalSvc.ListEntries(ctx)
	if err != nil {
		log.Printf("[ERROR] list entries for dashboard: %v", err)
	}

	// Limit to 10 latest entries
	recentEntries := entries
	if len(recentEntries) > 10 {
		recentEntries = recentEntries[:10]
	}

	data := map[string]any{
		"Summary":       summary,
		"RecentEntries": recentEntries,
	}

	tmpl, ok := h.templates["dashboard/index.html"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("[RENDER ERROR] dashboard: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// getSummary aggregates financial data from the trial balance.
func (h *Handler) getSummary(ctx context.Context) Summary {
	tb, err := h.reportSvc.GetTrialBalance(ctx)
	if err != nil {
		log.Printf("[ERROR] get trial balance for dashboard: %v", err)
		return Summary{}
	}

	var s Summary
	for _, row := range tb.Rows {
		net := row.DebitTotal - row.CreditTotal
		switch row.AccountType {
		case "asset":
			s.TotalAsset += net
		case "liability":
			s.TotalLiability += -net // credit-normal
		case "equity":
			s.TotalEquity += -net // credit-normal
		case "revenue":
			s.TotalRevenue += -net // credit-normal
		case "expense":
			s.TotalExpense += net
		case "cogs":
			s.TotalCOGS += net
		}
	}

	s.NetIncome = s.TotalRevenue - s.TotalExpense - s.TotalCOGS

	return s
}
