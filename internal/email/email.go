package email

import (
	"context"
)

type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, content string) error
}

type RealEmailSender struct{}

func (s *RealEmailSender) SendEmail(ctx context.Context, to, subject, content string) error {
	// Real email sending logic
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
