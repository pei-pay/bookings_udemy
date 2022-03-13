package handlers

import (
	"net/http"

	"github.com/pei-pay/bookings_udemy/pkg/config"
	"github.com/pei-pay/bookings_udemy/pkg/models"
	"github.com/pei-pay/bookings_udemy/pkg/render"
)

// Repo the repositrory used by the handerers
var Repo *Repository

// Repositroy is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHanders sets the repository for the handers
func NewHanders(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIp

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
