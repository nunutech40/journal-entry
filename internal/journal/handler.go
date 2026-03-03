package journal

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"journal-entry/internal/account"
)

// Handler handles HTTP requests for journal entries.
type Handler struct {
	svc        Service
	accountSvc account.Service
	templates  map[string]*template.Template
}

// NewHandler creates a new journal Handler.
func NewHandler(svc Service, accountSvc account.Service, templates map[string]*template.Template) *Handler {
	return &Handler{svc: svc, accountSvc: accountSvc, templates: templates}
}

// RegisterRoutes registers journal routes on the given chi router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/journals", func(r chi.Router) {
		r.Get("/", h.HandleList)
		r.Get("/new", h.HandleCreateForm)
		r.Post("/", h.HandleCreate)
		r.Get("/{id}/edit", h.HandleEditForm)
		r.Put("/{id}", h.HandleUpdate)
		r.Delete("/{id}", h.HandleDelete)
	})
}

// HandleList renders the journal entry list page.
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.ListEntries(r.Context())
	if err != nil {
		log.Printf("[ERROR] list entries: %v", err)
		http.Error(w, "Gagal memuat daftar jurnal", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Entries": entries,
	}

	h.render(w, "journal/list.html", data)
}

// HandleCreateForm renders the journal entry create form.
func (h *Handler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountSvc.ListAccounts(r.Context())
	if err != nil {
		log.Printf("[ERROR] list accounts for form: %v", err)
		http.Error(w, "Gagal memuat daftar akun", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Accounts": accounts,
		"IsEdit":   false,
		"Entry":    nil,
	}

	h.render(w, "journal/form.html", data)
}

// HandleCreate processes the journal entry creation.
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	req := h.parseFormToRequest(r)

	_, err := h.svc.CreateEntry(r.Context(), req)
	if err != nil {
		h.renderFormWithError(w, r, nil, false, friendlyError(err))
		return
	}

	w.Header().Set("HX-Redirect", "/journals")
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Jurnal berhasil dibuat!", "type": "success"}}`)
	w.WriteHeader(http.StatusCreated)
}

// HandleEditForm renders the journal entry edit form.
func (h *Handler) HandleEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	entry, err := h.svc.GetEntry(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrEntryNotFound) {
			http.Error(w, "Jurnal tidak ditemukan", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	accounts, err := h.accountSvc.ListAccounts(r.Context())
	if err != nil {
		http.Error(w, "Gagal memuat daftar akun", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Accounts": accounts,
		"IsEdit":   true,
		"Entry":    entry,
	}

	h.render(w, "journal/form.html", data)
}

// HandleUpdate processes the journal entry update.
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	createReq := h.parseFormToRequest(r)
	updateReq := UpdateRequest{
		EntryDate:       createReq.EntryDate,
		Description:     createReq.Description,
		ReferenceNumber: createReq.ReferenceNumber,
		Lines:           createReq.Lines,
	}

	entry, _ := h.svc.GetEntry(r.Context(), id)
	_, err = h.svc.UpdateEntry(r.Context(), id, updateReq)
	if err != nil {
		h.renderFormWithError(w, r, entry, true, friendlyError(err))
		return
	}

	w.Header().Set("HX-Redirect", "/journals")
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Jurnal berhasil diperbarui!", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}

// HandleDelete deletes a journal entry.
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteEntry(r.Context(), id); err != nil {
		if errors.Is(err, ErrEntryNotFound) {
			http.Error(w, "Jurnal tidak ditemukan", http.StatusNotFound)
			return
		}
		http.Error(w, "Gagal menghapus jurnal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Jurnal berhasil dihapus!", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}

// --- Helpers ---

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

func (h *Handler) renderFormWithError(w http.ResponseWriter, r *http.Request, entry *JournalEntry, isEdit bool, errMsg string) {
	accounts, _ := h.accountSvc.ListAccounts(r.Context())
	data := map[string]any{
		"Accounts": accounts,
		"IsEdit":   isEdit,
		"Entry":    entry,
		"Error":    errMsg,
		"FormData": r.Form,
	}
	tmpl := h.templates["journal/form.html"]
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusUnprocessableEntity)
	tmpl.ExecuteTemplate(w, "base", data)
}

// parseFormToRequest parses HTML form data into a CreateRequest.
// Lines are submitted as line_account_id[], line_debit[], line_credit[], line_description[].
func (h *Handler) parseFormToRequest(r *http.Request) CreateRequest {
	accountIDs := r.Form["line_account_id[]"]
	debits := r.Form["line_debit[]"]
	credits := r.Form["line_credit[]"]
	descs := r.Form["line_description[]"]

	var lines []CreateLineRequest
	for i := 0; i < len(accountIDs); i++ {
		accID, _ := strconv.Atoi(accountIDs[i])
		debit := parseFloat(safeIndex(debits, i))
		credit := parseFloat(safeIndex(credits, i))

		// Skip completely empty lines
		if accID == 0 && debit == 0 && credit == 0 {
			continue
		}

		lines = append(lines, CreateLineRequest{
			AccountID:   accID,
			Description: safeIndex(descs, i),
			Debit:       debit,
			Credit:      credit,
		})
	}

	return CreateRequest{
		EntryDate:       r.FormValue("entry_date"),
		Description:     r.FormValue("description"),
		ReferenceNumber: r.FormValue("reference_number"),
		Lines:           lines,
	}
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func safeIndex(slice []string, i int) string {
	if i < len(slice) {
		return slice[i]
	}
	return ""
}

func friendlyError(err error) string {
	switch {
	case errors.Is(err, ErrDateRequired):
		return "Tanggal wajib diisi"
	case errors.Is(err, ErrDateInvalid):
		return "Format tanggal tidak valid (gunakan YYYY-MM-DD)"
	case errors.Is(err, ErrDescriptionRequired):
		return "Keterangan wajib diisi"
	case errors.Is(err, ErrMinTwoLines):
		return "Minimal 2 baris entri diperlukan"
	case errors.Is(err, ErrLineNoAccount):
		return err.Error()
	case errors.Is(err, ErrLineNoAmount):
		return err.Error()
	case errors.Is(err, ErrLineBothAmounts):
		return err.Error()
	case errors.Is(err, ErrNotBalanced):
		return "Total debit harus sama dengan total kredit"
	case errors.Is(err, ErrAccountNotFound):
		return err.Error()
	case errors.Is(err, ErrEntryNotFound):
		return "Jurnal tidak ditemukan"
	default:
		return "Terjadi kesalahan, silakan coba lagi"
	}
}
