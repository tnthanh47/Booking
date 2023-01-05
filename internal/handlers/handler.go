package handlers

import (
	"encoding/json"
	"github.com/tnthanh47/Booking/internal/config"
	"github.com/tnthanh47/Booking/internal/forms"
	"github.com/tnthanh47/Booking/internal/models"
	"github.com/tnthanh47/Booking/internal/render"
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

	remoteIp := Repo.App.Session.GetString(req.Context(), "remote_ip")
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

// Đặt trước
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(
		w, r, "make-reservation.page.html", &models.TemplateData{
			Form: forms.New(nil),
		},
	)
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)
	form.Has("first_name", r)
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservarion"] = reservation

		render.RenderTemplate(
			w, r, "make-reservation.page.html", &models.TemplateData{
				Form: form,
				Data: data,
			},
		)
	}
}
