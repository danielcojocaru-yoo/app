package handlers

import (
	"encoding/json"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Error:   "Bad Request",
			Message: "Invalid JSON data",
		})
		return
	}
	// Invalidate any existing code for the user's email
	invalidateCode(user.Email)

	// Check if email exists
	exists, err := emailExists(user.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Message: "Error checking email existence",
		})
		return
	}

	if exists {
		// Generate a new code
		codeInfo := generateCode()
		loginCodes[user.Email] = codeInfo
		// Pass the "login" templateType when calling sendVerificationEmail
		if err := sendVerificationEmail(user.Email, codeInfo.Code, "login"); err != nil {
			// Handle error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Message: "I've just sent you a code via email. Please check your email!",
		})
	} else {
		// Email not found
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Error:   "Unauthorized Access",
			Message: "The email address does not exist in the database. Would you like to create an account?",
		})
	}
}
