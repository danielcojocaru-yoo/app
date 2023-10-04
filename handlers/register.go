package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	// Import any other packages you need for user-related operations
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Message: "Method not allowed",
		})
		return
	}

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "", "Eror encode user")
		return
	}

	email := user.Email

	exists, err := emailExists(email)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, "", "Error checking email")
		return
	}

	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{
			Error:   "Unauthorized Access",
			Message: "The email address already exists in our database. Please go to the Login page to connect.",
		})
		return
	}

	codeInfo := generateCode()
	emailCodes[email] = codeExpiry{
		Code: strconv.Itoa(codeInfo.Code),
		Exp:  codeInfo.Exp,
	}
	if err := sendVerificationEmail(email, codeInfo.Code, "/register"); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Message: "Oops, a technical error occurred while sending the email. Please try again later.",
		})
		return
	}

	// Store the expected verification code for comparison later
	expectedCode = codeInfo.Code

	// Email sent successfully
	sendJSONResponse(w, http.StatusOK, "Verification email sent successfully. Please check your email for the verification code.", "")

}

func verifyRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	email := r.FormValue("email")

	storedCode, exists := emailCodes[email]
	if !exists || code != storedCode.Code {
		// Code doesn't match

		sendJSONResponse(w, http.StatusUnauthorized, "Unauthorized", "Invalid verification code")
		return
	} else {

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Response{
				Message: "Error decoding user",
			})
			return

		}

		// Generate JWT token for the registered user
		jwtToken, err := generateJWT(user.Email)
		if err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, "Error", "Error generating JWT token")
		}

		// Generate refresh token for the registered user
		refreshToken := generateRefreshToken(user.Email)

		// Call createUser with tokens
		if err := createUser(user, jwtToken, refreshToken); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, "Error", "Error creating user")
			return
		}

		if err := saveEmail(user.Email); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, "Error", "Error saving email")
			return
		}

		// Return success message along with the JWT token
		sendJSONResponse(w, http.StatusOK, "Registration successful. Please go to the Login page and discover Yoosend!", "")

	}
}
