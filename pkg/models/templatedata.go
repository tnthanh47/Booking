package models

type TemplateData struct {
	MapString map[string]string
	CSRFToken string
	Data      map[string]interface{}
	Warning   string
	Error     string
}
