package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/tnthanh47/Booking/internal/config"
	"github.com/tnthanh47/Booking/internal/models"
	"github.com/tnthanh47/Booking/internal/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var appConfig config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	// Register to store kind of data in session
	gob.Register(models.Reservation{})

	appConfig.IsProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = false
	session.Cookie.Secure = appConfig.IsProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := NewRepo(&appConfig)
	NewHandler(repo)
	render.NewTemplateCache(&appConfig)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/SearchAvailability", Repo.SearchAvailabilityJSON)
	mux.Post("/PostedAvailabilityJSON", Repo.PostedAvailabilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func NoSurf(next http.Handler) http.Handler {

	csrf := nosurf.New(next)

	csrf.SetBaseCookie(
		http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   appConfig.IsProduction,
			SameSite: http.SameSiteLaxMode,
		},
	)

	return csrf
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplate))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
