package main

import (
	"net/http"
)

func (app *application) logError(request *http.Request, err error) {
	app.logger.Error("system error", "method", request.Method, "URI", request.URL.Path, "err", err.Error())
}

func (app *application) serverErrorResponse(responseWriter http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)
	app.renderError(responseWriter, request, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) notFoundResponse(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (app *application) rateLimitExceededResponse(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusTooManyRequests, "You're making requests too fast. Wait a few minutes and try again")
}

func (app *application) unprocessableEntityResponse(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
}

func (app *application) serviceUnavailable(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusServiceUnavailable, "Service is temporarily unavailable. Wait a few minutes and try again")
}

func (app *application) badRequest(responseWriter http.ResponseWriter, request *http.Request) {
	app.renderError(responseWriter, request, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
}
