package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

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
