package handlers

import (
	"encoding/json"
	"github.com/tnthanh47/Booking/pkg/config"
	"github.com/tnthanh47/Booking/pkg/models"
	"github.com/tnthanh47/Booking/pkg/render"
	"log"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, request *http.Request) {

	remoteIp := request.RemoteAddr

	m.App.Session.Put(request.Context(), "remote_ip", remoteIp)
	//Perform some logic
	strMap := map[string]string{}
	strMap["test"] = "hello"
	render.RenderTemplate(
		w, request, "home.page.html", &models.TemplateData{
			MapString: strMap,
		},
	)
}

func (m *Repository) About(w http.ResponseWriter, req *http.Request) {

	remoteIp := m.App.Session.GetString(req.Context(), "remote_ip")
	sessionLifeTime := m.App.Session.Lifetime

	strMap := map[string]string{}
	strMap["test"] = "hello"
	strMap["remote_ip"] = remoteIp
	strMap["session_life_time"] = sessionLifeTime.String()
	render.RenderTemplate(
		w, req, "about.page.html", &models.TemplateData{
			MapString: strMap,
		},
	)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) SearchAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		OK:      true,
		Message: "success",
	}
	out, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	write, err := w.Write(out)
	if err != nil {
		return
	}

	log.Println(write)
}

func (m *Repository) PostedAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Posted"))
}
