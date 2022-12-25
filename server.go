package main

import (
	"fmt"
	"net/http"
	"os"

	"PasswordServer2/api"
	"PasswordServer2/frontend"

	psDatabase "PasswordServer2/lib/database"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	psDatabase.DatabaseConnect()

	http.Handle("/", frontend.SvelteKitHandler("/"))
	http.Handle("/api/v1/", api.APIHandler("/api/v1"))

	listenLocation := "0.0.0.0:8000"
	envListenLocation := os.Getenv("LISTEN_LOCATION")
	if envListenLocation != "" {
		listenLocation = envListenLocation
	}

	fmt.Println("Listening on: " + listenLocation)
	fmt.Println("Starting server")

	http.ListenAndServe(listenLocation, nil)
}
