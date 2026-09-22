package main

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func renderTemplate(responseWriter http.ResponseWriter, template templ.Component) error {
	return template.Render(context.Background(), responseWriter)
}
