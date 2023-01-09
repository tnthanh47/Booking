package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/tnthanh47/Booking/internal/config"
	"github.com/tnthanh47/Booking/internal/handlers"
	"github.com/tnthanh47/Booking/internal/models"
	"github.com/tnthanh47/Booking/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {

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

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandler(repo)
	render.NewTemplateCache(&appConfig)

	fmt.Printf(fmt.Sprintf("Start Application listening to Port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
