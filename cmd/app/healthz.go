package main

import "net/http"

func (app *application) healthzHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte("ok"))
}
