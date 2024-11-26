package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, content string) error
}

type RealEmailSender struct{}

func (s *RealEmailSender) SendEmail(ctx context.Context, to, subject, content string) error {
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

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonPayload))
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
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("email: failed to send: status=%d body=%s", resp.StatusCode, body)
	}

	log.Printf("email: sent password reset email to: %s", to)
	return nil
}

type MockEmailSender struct {
	SentEmails []SentEmail
}

type SentEmail struct {
	To      string
	Subject string
	Content string
}

func (m *MockEmailSender) SendEmail(ctx context.Context, to, subject, content string) error {
	m.SentEmails = append(m.SentEmails, SentEmail{To: to, Subject: subject, Content: content})
	return nil
}
