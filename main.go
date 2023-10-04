package main

import (
	"net/http" // Relative import for custom packages
	// Relative import for custom packages

	"github.com/gorilla/mux"
)

// Main function
func main() {
	r := mux.NewRouter()
	// Define routes
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/verify", verifyRegistrationHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/verifylogin", verifyLoginHandler).Methods("GET")
	r.HandleFunc("/refreshtoken", refreshTokenHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("POST")

	r.HandleFunc("/update-tag", UpdateTagHandler).Methods("POST")
	r.HandleFunc("/avatar", UpdateAvatarHandler).Methods("POST")

	http.Handle("/", r)

	// Start your HTTP server
	http.ListenAndServe(":8080", nil)
}
