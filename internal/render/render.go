package render

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/tnthanh47/Booking/internal/config"
	"github.com/tnthanh47/Booking/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var pathToTemplate = "../../templates"

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplateCache(config *config.AppConfig) {
	app = config
}

func InitData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, tmpData *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//tc, _ = CreateTemplateCache()

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from cache.")
	}

	buf := new(bytes.Buffer)

	tmpData = InitData(tmpData, r)

	_ = t.Execute(buf, tmpData)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error when write template to browser")
	}

}

func Test() {
	var test = map[string]int{}
	fmt.Println(test)
}

func CreateTemplateCache() (map[string]*template.Template, error) {
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
