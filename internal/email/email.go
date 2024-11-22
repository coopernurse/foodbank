package email

import (
	"context"
	"fmt"
	"log"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

// SendEmail sends an email using Brevo
func SendEmail(ctx context.Context, to, subject, content string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	url := "https://api.brevo.com/v3/smtp/email"

	payload := fmt.Sprintf(`{
		"sender": {
			"name": "James Cooper",
			"email": "james@bitmechanic.com"
		},
		"to": [
			{
				"email": "%s",
				"name": "Recipient"
			}
		],
		"subject": "%s",
		"htmlContent": "%s"
	}`, to, subject, content)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(payload))
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

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("email: failed to send: status=%d body=%s", resp.StatusCode, body)
	}

	log.Printf("Email sent successfully")
	return nil
}
