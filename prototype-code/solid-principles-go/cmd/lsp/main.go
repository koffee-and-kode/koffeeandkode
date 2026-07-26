// LSP — Liskov Substitution Principle.
//
// Every implementation of an interface must honor what callers expect.
// Pretending to deliver while not delivering breaks substitution. The fix
// is usually a sharper abstraction, not a sneakier implementation.
package main

import "fmt"

// The base contract: any way of dispatching a message.
type Notifier interface {
	Notify(to, message string) error
}

// A stricter contract layered on Notifier: also returns a delivery receipt
// for auditing. Only notifiers that genuinely deliver should satisfy this.
type DeliveryNotifier interface {
	Notifier
	Receipt() string
}

type EmailNotifier struct{ receipt string }

func (e *EmailNotifier) Notify(to, message string) error {
	fmt.Printf("[email] to=%s msg=%q\n", to, message)
	e.receipt = "email-msg-abc"
	return nil
}
func (e *EmailNotifier) Receipt() string { return e.receipt }

type SMSNotifier struct{ receipt string }

func (s *SMSNotifier) Notify(to, message string) error {
	fmt.Printf("[sms] to=%s msg=%q\n", to, message)
	s.receipt = "sms-msg-xyz"
	return nil
}
func (s *SMSNotifier) Receipt() string { return s.receipt }

// LogNotifier writes to a local log only — no real delivery, no receipt.
// It honestly satisfies Notifier but NOT DeliveryNotifier — by design.
type LogNotifier struct{}

func (LogNotifier) Notify(to, message string) error {
	fmt.Printf("[log] would-notify=%s msg=%q\n", to, message)
	return nil
}

// AuditedSend takes any DeliveryNotifier; the compiler stops you from
// passing a LogNotifier here, so callers never get lied to about delivery.
func AuditedSend(n DeliveryNotifier, to, message string) {
	_ = n.Notify(to, message)
	fmt.Printf("  audit: receipt=%s\n", n.Receipt())
}

func main() {
	notifiers := []Notifier{&EmailNotifier{}, &SMSNotifier{}, LogNotifier{}}
	for _, n := range notifiers {
		_ = n.Notify("ada@example.com", "hello")
		// Type-assert to the stricter contract before auditing.
		if dn, ok := n.(DeliveryNotifier); ok {
			AuditedSend(dn, "ada@example.com", "ALERT: server down")
		}
	}
}
