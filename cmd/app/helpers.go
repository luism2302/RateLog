package main

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
)

func (app *application) renderTemplate(responseWriter http.ResponseWriter, template templ.Component) error {
	return template.Render(context.Background(), responseWriter)
}

func (app *application) readTitleParam(request *http.Request) string {
	params := httprouter.ParamsFromContext(request.Context())
	return params.ByName("title")
}
