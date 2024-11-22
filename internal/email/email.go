package email

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// SendEmail sends an email using Brevo
func SendEmail(ctx context.Context, to, subject, content string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	url := "https://api.brevo.com/v3/smtp/email"

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "James Cooper",
			"email": "james@bitmechanic.com",
		},
		"to": []map[string]string{
			{
				"email": to,
				"name":  "Recipient",
			},
		},
		"subject":     subject,
		"htmlContent": content,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("email: failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("email: failed to create request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", apiKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("email: failed to send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("email: failed to send: status=%d body=%s", resp.StatusCode, body)
	}

	log.Printf("Email sent successfully")
	return nil
}
