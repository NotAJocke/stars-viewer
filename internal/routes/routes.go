package routes

import (
	"fmt"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", MainRoute)

	return mux
}

func MainRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
