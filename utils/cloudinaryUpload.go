package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadToCloudinary(fileHeader *multipart.FileHeader) (string, error) {
	// Step 1: Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Step 2: Prepare form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Attach file
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	// Attach upload preset
	if err := writer.WriteField("upload_preset", "ml_default"); err != nil {
		return "", fmt.Errorf("failed to write upload preset: %w", err)
	}
	writer.Close()

	// Step 3: Send POST request to Cloudinary
	cloudinaryURL := "https://api.cloudinary.com/v1_1/dcfoqhrxb/upload"
	req, err := http.NewRequest("POST", cloudinaryURL, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: Parse response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", body)
	}

	var result struct {
		SecureURL string `json:"secure_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.SecureURL == "" {
		return "", fmt.Errorf("upload failed: no secure_url in response")
	}

	return result.SecureURL, nil
}
