package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func generateRefreshToken(email string) string {
	// Generate a unique token source using email, current date, and a random UUID
	tokenSource := fmt.Sprintf("%s-%s-%s", email, time.Now().Format("20060102150405"), uuid.New().String())

	// Hash the token source using a strong hash function (SHA-512)
	hasher := sha512.New()
	hasher.Write([]byte(tokenSource))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}
func invalidateRefreshToken(refreshToken string) error {
	//   invalidate the refresh token

	return nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		http.Error(w, "Refresh token missing", http.StatusBadRequest)
		return
	}

	// Check if the refresh token exists in Directus
	refreshTokenExists, err := checkRefreshTokenExistsInDirectus(refreshToken)
	if err != nil {
		http.Error(w, "Error checking refresh token", http.StatusInternalServerError)
		return
	}

	if !refreshTokenExists {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse := map[string]interface{}{
			"error":   "Unauthorized",
			"message": "Refresh token not found",
		}
		json.NewEncoder(w).Encode(jsonResponse)
		return

	}

	// Invalidate the refresh token
	err = invalidateRefreshToken(refreshToken)
	if err != nil {
		http.Error(w, "Error invalidating refresh token", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Logged out successfully",
	})
}

func checkRefreshTokenExistsInDirectus(refreshToken string) (bool, error) {
	// Build the API request to check if the refresh token exists in Directus
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/items/users?filter[refresh_token]=%s", directusURL, refreshToken), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+directusApiKey)

	// Call API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Handle error status codes
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	// Decode response
	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	// Check if user data with matching refresh_token exists
	users, ok := res["data"].([]interface{})
	if !ok {
		return false, errors.New("Invalid response format")
	}

	return len(users) > 0, nil
}
