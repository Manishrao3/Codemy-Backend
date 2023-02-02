package main

import (
	"RestfulApi/router"
	"fmt"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	r := router.Router()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	fmt.Println("Starting server on the port 8080...")
	handler := cors.Handler(r)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
