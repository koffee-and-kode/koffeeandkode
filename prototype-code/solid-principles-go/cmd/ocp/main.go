// OCP — Open/Closed Principle.
//
// SendWelcome is closed for modification: its body doesn't change when a new
// channel is added. The Notifier interface is the open extension point —
// anyone can drop in a new implementation.
package main

import "fmt"

// The extension point. New channels implement Notify; nothing else moves.
type Notifier interface {
	Notify(to, message string) error
}

type EmailNotifier struct{}

func (EmailNotifier) Notify(to, message string) error {
	fmt.Printf("[email] to=%s msg=%q\n", to, message)
	return nil
}

type SMSNotifier struct{}

func (SMSNotifier) Notify(to, message string) error {
	fmt.Printf("[sms] to=%s msg=%q\n", to, message)
	return nil
}

// Added after the fact — SendWelcome below is unchanged.
type PushNotifier struct{}

func (PushNotifier) Notify(to, message string) error {
	fmt.Printf("[push] to=%s msg=%q\n", to, message)
	return nil
}

// Closed: depends only on the Notifier abstraction.
func SendWelcome(n Notifier, to string) error {
	return n.Notify(to, "Welcome to our app!")
}

func main() {
	channels := []Notifier{EmailNotifier{}, SMSNotifier{}, PushNotifier{}}
	for _, c := range channels {
		_ = SendWelcome(c, "ada@example.com")
	}
}
