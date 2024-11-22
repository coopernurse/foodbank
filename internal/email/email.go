package email

import (
	"log"
	"os"

	"github.com/mailgun/mailgun-go/v4"
)

// SendEmail sends an email using Mailgun
func SendEmail(to, subject, content string) error {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"), "", "")

	message := mg.NewMessage(
		"Your Name <your-email@example.com>",
		subject,
		content,
		to,
	)

	_, id, err := mg.Send(message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	log.Printf("Email sent with ID: %s", id)
	return nil
}
