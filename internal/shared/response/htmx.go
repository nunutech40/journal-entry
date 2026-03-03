package response

import (
	"html/template"
	"log"
	"net/http"
)

// RenderPage renders a full page template within the base layout.
// The template must define a "content" block.
func RenderPage(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("[RENDER ERROR] %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderPartial renders a partial template (for HTMX swaps).
// Does not include the base layout.
func RenderPartial(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("[RENDER ERROR] partial %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderError renders an error message as HTML partial.
// Typically swapped into an error container via HTMX.
func RenderError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	// Simple error HTML that can be swapped into any container
	w.Write([]byte(`<div class="c-alert c-alert--error">` + template.HTMLEscapeString(message) + `</div>`))
}

// SetTrigger sets the HX-Trigger header to fire a client-side event.
// Commonly used to trigger Alpine.js toast notifications.
//
// Example: SetTrigger(w, "showToast") → Alpine listens @show-toast.window
func SetTrigger(w http.ResponseWriter, event string) {
	w.Header().Set("HX-Trigger", event)
}

// SetTriggerWithData sets HX-Trigger with a JSON payload.
// Example: SetTriggerWithData(w, `{"showToast": {"message": "Berhasil!", "type": "success"}}`)
func SetTriggerWithData(w http.ResponseWriter, jsonPayload string) {
	w.Header().Set("HX-Trigger", jsonPayload)
}

// Redirect tells HTMX to do a client-side redirect.
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
}

// IsHTMX checks if the request was made by HTMX.
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
