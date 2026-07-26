// DIP — Dependency Inversion Principle.
//
// High-level policy (NotificationService) depends on an abstraction
// (EmailSender). Low-level transports (SMTP, SES) also depend on the
// abstraction — not the other way around. main wires concrete to abstract.
package main

import "fmt"

// The abstraction the high-level service depends on.
type EmailSender interface {
	SendEmail(to, subject, body string) error
}

// High-level policy — doesn't import any transport package.
type NotificationService struct {
	sender EmailSender
}

func NewNotificationService(s EmailSender) *NotificationService {
	return &NotificationService{sender: s}
}

func (s *NotificationService) SendWelcome(to string) error {
	return s.sender.SendEmail(to, "Welcome!", "Glad you're here.")
}

// --- low-level transports: depend on the abstraction, not vice versa ---

type SMTPSender struct{}

func (SMTPSender) SendEmail(to, subject, body string) error {
	fmt.Printf("[smtp] to=%s subj=%q\n", to, subject)
	return nil
}

type SESender struct{}

func (SESender) SendEmail(to, subject, body string) error {
	fmt.Printf("[ses] to=%s subj=%q\n", to, subject)
	return nil
}

func main() {
	// Pick the transport at the edge; the service body never changes.
	svc := NewNotificationService(SMTPSender{})
	_ = svc.SendWelcome("ada@example.com")

	svc = NewNotificationService(SESender{})
	_ = svc.SendWelcome("ada@example.com")
}
