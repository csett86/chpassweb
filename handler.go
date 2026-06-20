package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates map[string]*template.Template

func init() {
	// Parse base template first
	base := template.New("")
	base = template.Must(base.ParseFS(templateFS, "templates/base.html"))
	
	// Clone base and parse each page template to override blocks
	promptTmpl := template.Must(base.Clone())
	promptTmpl = template.Must(promptTmpl.ParseFS(templateFS, "templates/prompt.html"))
	
	errorTmpl := template.Must(base.Clone())
	errorTmpl = template.Must(errorTmpl.ParseFS(templateFS, "templates/error.html"))
	
	successTmpl := template.Must(base.Clone())
	successTmpl = template.Must(successTmpl.ParseFS(templateFS, "templates/success.html"))
	
	templates = map[string]*template.Template{
		"prompt.html": promptTmpl,
		"error.html":   errorTmpl,
		"success.html": successTmpl,
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	// Check if REMOTE_USER header exists
	if values := r.Header.Values("REMOTE_USER"); len(values) == 0 {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	
	username := r.Header.Get("REMOTE_USER")
	if username == "" {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	session := store.Create(username)
	go runPam(session)

	prompt, ok := <-session.PromptCh
	if !ok {
		http.Error(w, "PAM failed to start", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	renderPrompt(w, prompt)
}

func handleRespond(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Session expired", http.StatusBadRequest)
		return
	}

	session := store.Get(cookie.Value)
	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	response := r.FormValue("response")
	if response == "" {
		http.Error(w, "Response required", http.StatusBadRequest)
		return
	}

	session.ResponseCh <- response

	select {
	case prompt, ok := <-session.PromptCh:
		if !ok {
			if err := <-session.ResultCh; err != nil {
				renderError(w, err.Error())
			} else {
				renderSuccess(w)
			}
			store.Delete(session.ID)
			return
		}
		renderPrompt(w, prompt)
	case err := <-session.ResultCh:
		if err != nil {
			renderError(w, err.Error())
		} else {
			renderSuccess(w)
		}
		store.Delete(session.ID)
	}
}

func renderPrompt(w http.ResponseWriter, prompt string) {
	templates["prompt.html"].ExecuteTemplate(w, "base.html", struct {
		Prompt string
	}{Prompt: prompt})
}

func renderError(w http.ResponseWriter, errorMsg string) {
	templates["error.html"].ExecuteTemplate(w, "base.html", struct {
		Error string
	}{Error: errorMsg})
}

func renderSuccess(w http.ResponseWriter) {
	templates["success.html"].ExecuteTemplate(w, "base.html", nil)
}
