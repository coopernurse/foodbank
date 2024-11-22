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
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", os.Getenv("BREVO_API_KEY"))
	cfg.AddDefaultHeader("partner-key", os.Getenv("BREVO_API_KEY"))

	client := brevo.NewAPIClient(cfg)

	email := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "James Cooper",
			Email: "james@bitmechanic.com",
		},
		To:          []brevo.SendSmtpEmailTo{{Email: "james+100@bitmechanic.com"}},
		Subject:     "test of the email system",
		TextContent: "<html><head></head><body><p>Hello,</p>This is my first transactional email sent from Brevo.</p></body></html>",
	}

	log.Printf("Sending email: %+v", email) // Add logging here

	resp, httpResp, err := client.TransactionalEmailsApi.SendTransacEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("email: failed to send: resp=%v resp=%d %s err=%w", resp, httpResp.StatusCode, httpResp.Status, err)
	}
	log.Printf("Email sent with ID: %s", resp.MessageId)
	log.Printf("HTTP Response: %v", httpResp)
	return nil
}
