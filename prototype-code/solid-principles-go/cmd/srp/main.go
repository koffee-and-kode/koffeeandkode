// SRP — Single Responsibility Principle.
//
// WelcomeService has exactly one reason to change: the rule for sending a
// welcome notification. How users are stored and how mail is delivered are
// concerns of other types behind interfaces.
package main

import "fmt"

// Where users come from. The service doesn't care if it's Postgres, Redis,
// or an in-memory map.
type UserRepo interface {
	Email(userID int) (string, error)
}

// How a message is delivered. Could be SMTP, SES, or a fake — service unaware.
type Mailer interface {
	Send(to, subject, body string) error
}

// One responsibility: fetch a user's email, then send them a welcome.
type WelcomeService struct {
	users  UserRepo
	mailer Mailer
}

func (s *WelcomeService) SendWelcome(userID int) error {
	email, err := s.users.Email(userID)
	if err != nil {
		return err
	}
	return s.mailer.Send(email, "Welcome!", "Glad you're here.")
}

// --- bare-bones implementations for the demo ---

type memUsers map[int]string

func (m memUsers) Email(id int) (string, error) { return m[id], nil }

type stdoutMailer struct{}

func (stdoutMailer) Send(to, subject, body string) error {
	fmt.Printf("[mail] to=%s subj=%q body=%q\n", to, subject, body)
	return nil
}

func main() {
	// Wire concrete implementations to the service at the edge.
	svc := &WelcomeService{
		users:  memUsers{1: "ada@example.com"},
		mailer: stdoutMailer{},
	}
	_ = svc.SendWelcome(1)
}
