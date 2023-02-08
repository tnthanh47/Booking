package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/tnthanh47/Booking/internal/config"
	"github.com/tnthanh47/Booking/internal/driver"
	"github.com/tnthanh47/Booking/internal/handlers"
	"github.com/tnthanh47/Booking/internal/helper"
	"github.com/tnthanh47/Booking/internal/models"
	"github.com/tnthanh47/Booking/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Start Application listening to Port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {

	// Register to store kind of data in session
	gob.Register(models.Reservation{})

	appConfig.IsProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = false
	session.Cookie.Secure = appConfig.IsProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	appConfig.Session = session

	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandler(repo)
	render.NewTemplateCache(&appConfig)
	helper.NewHelper(&appConfig)

	return nil
}
