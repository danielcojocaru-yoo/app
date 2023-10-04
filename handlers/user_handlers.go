package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// UpdateTagHandler handles updating the tag_name based on the provided email.
func UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the email and tag_name parameters
	var request struct {
		Email   string `json:"email"`
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	email := request.Email
	tagName := request.TagName

	// Get user ID
	userID, err := getUserIDByEmail(email)
	if err != nil {
		http.Error(w, "Error fetching user ID", http.StatusInternalServerError)
		return
	}

	// Update the tag_name in Directus for the given user
	if err := updateTagName(userID, tagName); err != nil {
		response := map[string]string{"error": "Error updating tag_name"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"message": "Tag_name updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateTagName updates the tag_name in Directus for a given user.
func updateTagName(userID, tagName string) error {
	// Build the API request to update the tag_name
	reqBody := map[string]interface{}{
		"tag_name": tagName,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/items/users/%s", directusURL, userID)
	req, err := http.NewRequest(http.MethodPatch, reqURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+directusApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Handle error
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

func getUserIDByEmail(email string) (string, error) {
	// Build API request
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/items/users?filter[email]=%s", directusURL, email), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+directusApiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Handle response
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	users := data["data"].([]interface{})
	if len(users) == 0 {
		return "", errors.New("User not found")
	}

	user := users[0].(map[string]interface{})
	userID, ok := user["id"].(string)
	if !ok {
		// If it's not a string, try to convert it to a string
		userIDFloat, ok := user["id"].(float64)
		if !ok {
			return "", fmt.Errorf("Unexpected user ID type: %T", user["id"])
		}
		userID = fmt.Sprintf("%.0f", userIDFloat)
	}

	return userID, nil
}

// UpdateAvatarHandler handles updating the avatar_url based on the provided email.
func UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the email and avatar_url parameters
	var request struct {
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	email := request.Email
	avatarURL := request.AvatarURL

	// Get user ID
	userID, err := getUserIDByEmail(email)
	if err != nil {
		http.Error(w, "Error fetching user ID", http.StatusInternalServerError)
		return
	}

	// Update the avatar_url in Directus for the given user
	if err := updateAvatarURL(userID, avatarURL); err != nil {
		response := map[string]string{"error": "Error updating avatar_url"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]string{"message": "Avatar_url updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateAvatarURL updates the avatar_url in Directus for a given user.
func updateAvatarURL(userID, avatarURL string) error {
	// Build the API request to update the avatar_url
	reqBody := map[string]interface{}{
		"avatar_url": avatarURL,
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	reqURL := fmt.Sprintf("%s/items/users/%s", directusURL, userID)
	req, err := http.NewRequest(http.MethodPatch, reqURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+directusApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Handle error
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}
