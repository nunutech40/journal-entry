package account

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for accounts.
type Handler struct {
	svc       Service
	templates map[string]*template.Template
}

// NewHandler creates a new account Handler.
func NewHandler(svc Service, templates map[string]*template.Template) *Handler {
	return &Handler{svc: svc, templates: templates}
}

// RegisterRoutes registers account routes on the given chi router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", h.HandleList)
		r.Get("/new", h.HandleCreateForm)
		r.Post("/", h.HandleCreate)
		r.Get("/{id}/edit", h.HandleEditForm)
		r.Put("/{id}", h.HandleUpdate)
		r.Delete("/{id}", h.HandleDelete)
	})
}

// HandleList renders the full account list page.
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.svc.ListAccounts(r.Context())
	if err != nil {
		log.Printf("[ERROR] list accounts: %v", err)
		http.Error(w, "Gagal memuat daftar akun", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Accounts":   accounts,
		"TypeLabels": buildTypeLabels(),
	}

	tmpl, ok := h.templates["account/list.html"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("[RENDER ERROR] account list: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// HandleCreateForm renders the create account form.
func (h *Handler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Account":    nil,
		"Types":      ValidTypes(),
		"TypeLabels": buildTypeLabels(),
		"IsEdit":     false,
	}

	tmpl, ok := h.templates["account/form.html"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("[RENDER ERROR] account form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// HandleCreate processes the create account form submission.
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal membaca form", http.StatusBadRequest)
		return
	}

	req := CreateRequest{
		Code:        r.FormValue("code"),
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Description: r.FormValue("description"),
	}

	_, err := h.svc.CreateAccount(r.Context(), req)
	if err != nil {
		// Re-render form with error
		data := map[string]any{
			"Account":    req,
			"Types":      ValidTypes(),
			"TypeLabels": buildTypeLabels(),
			"IsEdit":     false,
			"Error":      friendlyError(err),
		}
		tmpl := h.templates["account/form.html"]
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		tmpl.ExecuteTemplate(w, "base", data)
		return
	}

	// Success — redirect to list with toast
	w.Header().Set("HX-Redirect", "/accounts")
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Akun berhasil ditambahkan!", "type": "success"}}`)
	w.WriteHeader(http.StatusCreated)
}

// HandleEditForm renders the edit account form.
func (h *Handler) HandleEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	acc, err := h.svc.GetAccount(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			http.Error(w, "Akun tidak ditemukan", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Account":    acc,
		"Types":      ValidTypes(),
		"TypeLabels": buildTypeLabels(),
		"IsEdit":     true,
	}

	tmpl, ok := h.templates["account/form.html"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("[RENDER ERROR] account edit form: %v", err)
	}
}

// HandleUpdate processes the edit account form submission.
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

	req := UpdateRequest{
		Code:        r.FormValue("code"),
		Name:        r.FormValue("name"),
		Type:        r.FormValue("type"),
		Description: r.FormValue("description"),
		IsActive:    r.FormValue("is_active") == "on",
	}

	_, err = h.svc.UpdateAccount(r.Context(), id, req)
	if err != nil {
		acc, _ := h.svc.GetAccount(r.Context(), id)
		data := map[string]any{
			"Account":    acc,
			"Types":      ValidTypes(),
			"TypeLabels": buildTypeLabels(),
			"IsEdit":     true,
			"Error":      friendlyError(err),
		}
		tmpl := h.templates["account/form.html"]
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		tmpl.ExecuteTemplate(w, "base", data)
		return
	}

	w.Header().Set("HX-Redirect", "/accounts")
	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Akun berhasil diperbarui!", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}

// HandleDelete deletes an account.
func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteAccount(r.Context(), id); err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			http.Error(w, "Akun tidak ditemukan", http.StatusNotFound)
			return
		}
		http.Error(w, "Gagal menghapus akun", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"showToast": {"message": "Akun berhasil dihapus!", "type": "success"}}`)
	w.WriteHeader(http.StatusOK)
}

// --- Helpers ---

func buildTypeLabels() map[string]string {
	labels := make(map[string]string)
	for _, t := range ValidTypes() {
		labels[t] = TypeLabel(t)
	}
	return labels
}

func friendlyError(err error) string {
	switch {
	case errors.Is(err, ErrCodeRequired):
		return "Kode akun wajib diisi"
	case errors.Is(err, ErrNameRequired):
		return "Nama akun wajib diisi"
	case errors.Is(err, ErrInvalidType):
		return "Tipe akun tidak valid"
	case errors.Is(err, ErrCodeExists):
		return "Kode akun sudah digunakan"
	case errors.Is(err, ErrAccountNotFound):
		return "Akun tidak ditemukan"
	default:
		return "Terjadi kesalahan, silakan coba lagi"
	}
}
