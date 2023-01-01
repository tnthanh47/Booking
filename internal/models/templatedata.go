package models

import "github.com/tnthanh47/Booking/internal/forms"

type TemplateData struct {
	MapString map[string]string
	CSRFToken string
	Data      map[string]interface{}
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
