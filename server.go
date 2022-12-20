package main

import (
	"net/http"

	"PasswordServer2/api"
	"PasswordServer2/frontend"
)

func main() {
	http.Handle("/", frontend.SvelteKitHandler("/"))
	http.Handle("/api/v1/", api.APIHandler("/api/v1"))

	http.ListenAndServe("0.0.0.0:8000", nil)
}
