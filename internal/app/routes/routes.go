package routes

import (
	"log"
	"net/http"
	"url-short/internal/app/api"
	"url-short/internal/app/middleware"
	"url-short/internal/app/pages"
)

func Register() {
	api.InitializeShorten()

	http.HandleFunc("/", pages.Main)
	middleware.API("/api/v1/shorten", "POST", api.Shorten)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
