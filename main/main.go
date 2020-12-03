package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// --- Configuration --- //

var config Configuration

func getConfig(ENV string) Configuration {
	file, err := os.Open(fmt.Sprintf("config.%s.json", ENV))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Configuration
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// --- Main --- //

func main() {
	// Get configuration
	ENV := os.Getenv("ENV")
	if ENV == "" {
		ENV = "dev"
	}
	fmt.Println(fmt.Sprintf("Running in ENV: %s", ENV))
	config = getConfig(ENV)

	db = connectDb(config.Db)
	defer db.Close()
	pingDb(db)

	// Init router
	r := mux.NewRouter()

	// Route handlers
	r.HandleFunc("/auth/", home).Methods("GET")
	r.HandleFunc("/auth/register", registerPage).Methods("GET")
	r.HandleFunc("/auth/register", register).Methods("POST")
	r.HandleFunc("/auth/login", loginPage).Methods("GET")
	r.HandleFunc("/auth/login", login).Methods("POST")
	r.HandleFunc("/auth/password", passwordPage).Methods("GET")
	r.HandleFunc("/auth/password", password).Methods("POST")
	r.HandleFunc("/auth/logout", logout).Methods("GET")

	// CORS
	// handler := cors.Default().Handler(r)

	// Run server
	port := 8080
	fmt.Println(fmt.Sprintf("Serving on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}