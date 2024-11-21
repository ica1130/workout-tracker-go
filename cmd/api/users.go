package main

import (
	"fmt"
	"net/http"
)

func (app *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "here are the users")
}
