package main

import (
	"net/http"

	"PasswordServer2/api"
	"PasswordServer2/frontend"

	psDatabase "PasswordServer2/lib/database"
	psErrors "PasswordServer2/lib/errors"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	database := psDatabase.DatabaseConnect()
	if database == nil {
		panic(psErrors.ErrorLoadingDatabase)
	}

	http.Handle("/", frontend.SvelteKitHandler("/"))
	http.Handle("/api/v1/", api.APIHandler("/api/v1"))

	http.ListenAndServe("0.0.0.0:8000", nil)
}
