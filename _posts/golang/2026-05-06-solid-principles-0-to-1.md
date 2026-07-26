---
title: "SOLID Principles 0 to 1"
date: 2026-05-06
tags: [design, go, solid, oop]
go_categories: [idioms-and-design]
excerpt: "What each letter in SOLID actually buys you — written as five bare-bones Go programs, one per principle."
read_time: 8
repo_url: https://github.com/koffee-and-kode/koffeeandkode/tree/main/prototype-code/solid-principles-go
---

SOLID is five rules for designing code that doesn't fight you when it grows. They were written for class-heavy OOP, but the ideas translate cleanly to Go: packages instead of classes, interfaces instead of inheritance, composition instead of subclass hierarchies.

This post walks through each letter — what it means, what violating it looks like, and a tiny Go program showing the principle in action. The runnable code lives at [`prototype-code/solid-principles-go/`](https://github.com/koffee-and-kode/koffeeandkode/tree/main/prototype-code/solid-principles-go), one `cmd/<letter>/main.go` per principle.

Each example is intentionally bare-bones. No frameworks, no DBs, no real I/O — just enough wiring to show what the principle actually changes about the code.

## S — Single Responsibility Principle

**A module should have one reason to change.**

The smell: a `Service` or `Manager` struct that fetches from the DB, formats output, sends email, *and* writes audit logs. Five reasons to change, all colliding in one type.

```go
// Before: one struct, four reasons to change.
func (n *NotificationService) SendWelcomeEmail(userID int) error {
    var email string
    n.db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
    sendSMTP(email, "Welcome!", "Hi, thanks for signing up.")
    log.Printf("sent welcome email to user=%d", userID)
    return nil
}
```

The fix is to split each concern behind an interface and have the service depend on those interfaces — not the concrete DB or SMTP client:

```go
type UserRepo interface {
    Email(userID int) (string, error)
}
type Mailer interface {
    Send(to, subject, body string) error
}

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
```

`WelcomeService` now owns exactly one rule: *when a user signs up, send them a welcome.* The DB schema and the SMTP library can change without touching this file — and adding new collaborators (a logger, a metrics emitter) is just another interface field, not a rewrite.

SRP isn't about counting methods. It's about isolating which parts change for which reasons.

## O — Open/Closed Principle

**Open for extension, closed for modification.**

The smell is a big `switch` over a type tag:

```go
func SendNotification(ch Channel, to, message string) error {
    switch ch {
    case Email: return sendEmail(to, message)
    case SMS:   return sendSMS(to, message)
    }
    return fmt.Errorf("unsupported channel: %s", ch)
}
```

Every new channel forces an edit to this function and a recompile of everything that imports it. The function isn't *closed* — it's the central point all changes pass through.

The fix is polymorphism via an interface:

```go
type Notifier interface {
    Notify(to, message string) error
}

type EmailNotifier struct{}
func (EmailNotifier) Notify(to, message string) error { /* ... */ }

type SMSNotifier struct{}
func (SMSNotifier) Notify(to, message string) error { /* ... */ }

func SendWelcome(n Notifier, to string) error {
    return n.Notify(to, "Welcome to our app!")
}
```

Adding push notifications is a new file with a new type — `SendWelcome` stays untouched:

```go
type PushNotifier struct{}
func (PushNotifier) Notify(to, message string) error { /* ... */ }
```

The win isn't elegance. It's that the radius of any future change is one new file, not edits scattered through the call sites of a switch.

## L — Liskov Substitution Principle

**Any implementation of an interface should be usable wherever the interface is expected — without surprising the caller.**

Staying in the notification domain, here's the classic violation:

```go
type Notifier interface {
    Notify(to, message string) error
}

type LogNotifier struct{}
func (LogNotifier) Notify(to, message string) error {
    // "Fake" implementation: pretends to send but only writes a log line.
    return nil
}
```

`LogNotifier` *technically* satisfies `Notifier`, but anything calling `Notify` and expecting a real delivery just got lied to. The interface failed to capture that delivery is what callers care about.

Since Go's interface satisfaction is implicit, this kind of mismatch is easy to introduce. The fix is to sharpen the abstraction so the contract reflects what callers actually need:

```go
type Notifier interface {
    Notify(to, message string) error
}

type DeliveryNotifier interface {
    Notifier
    Receipt() string  // returns a delivery receipt id
}

type EmailNotifier struct{ receipt string }
func (e *EmailNotifier) Notify(to, msg string) error { /* ... */ }
func (e *EmailNotifier) Receipt() string             { return e.receipt }

type LogNotifier struct{}
func (LogNotifier) Notify(to, msg string) error { /* writes a log line */ }
// LogNotifier satisfies Notifier, NOT DeliveryNotifier — by design.

func AuditedSend(n DeliveryNotifier, to, msg string) {
    n.Notify(to, msg)
    log.Printf("audit: receipt=%s", n.Receipt())
}
```

`AuditedSend` takes a `DeliveryNotifier`, and every `DeliveryNotifier` actually delivers. The compiler stops you from passing `LogNotifier`.

LSP violations almost always point at a *bad abstraction*, not a bad implementation. If types need to lie about what they do to fit the interface, the interface is wrong.

## I — Interface Segregation Principle

**Clients should not depend on methods they don't use.**

Same notification system. A fat `MessageStore` interface forces every implementation — including the test fake the sender uses — to satisfy methods that have nothing to do with sending:

```go
type MessageStore interface {
    Save(m Message) error
    Get(id int) (Message, error)
    Delete(id int) error
    List() ([]Message, error)
    ExportCSV(w io.Writer) error
}
```

Now a test that only exercises the send path still has to stub `ExportCSV`. The coupling has nothing to do with what the test is testing.

Split by what callers actually need:

```go
type MessageSaver interface {
    Save(m Message) error
}

type MessageReader interface {
    Get(id int) (Message, error)
}

func RecordSent(s MessageSaver, m Message) error {
    return s.Save(m)
}

func ViewSent(r MessageReader, id int) (Message, error) {
    return r.Get(id)
}
```

A concrete type can still implement all of them — but each function only sees the slice it needs. This is why the Go stdlib has `io.Reader`, `io.Writer`, `io.ReadWriter` instead of one `IO` interface. Most idiomatic Go interfaces are 1–3 methods on purpose.

## D — Dependency Inversion Principle

**High-level modules shouldn't depend on low-level modules. Both should depend on abstractions.**

This one sounds abstract until you see the violation: the *business logic* file imports the *SMTP library* directly.

```go
type NotificationService struct {
    smtpHost string
    smtpPort int
}

func (s *NotificationService) SendWelcome(to string) error {
    // low-level SMTP code right here
}
```

Now the welcome-email rule and the SMTP transport are welded together. Want to switch to SES for one environment? Want to test without a real mail server? You can't, without rewriting the service.

Invert it: the service depends on an interface, and the transport implements it.

```go
type EmailSender interface {
    SendEmail(to, subject, body string) error
}

type NotificationService struct {
    sender EmailSender
}

func NewNotificationService(s EmailSender) *NotificationService {
    return &NotificationService{sender: s}
}

func (s *NotificationService) SendWelcome(to string) error {
    return s.sender.SendEmail(to, "Welcome!", "Glad you're here.")
}

type SMTPSender struct{}
func (SMTPSender) SendEmail(to, subject, body string) error { /* real SMTP */ }

type SESender struct{}
func (SESender) SendEmail(to, subject, body string) error { /* AWS SES */ }
```

The arrow now points the other way: `SMTPSender` depends on `EmailSender`, not the reverse. `main` (or a wire-style setup) picks the implementation at the edge:

```go
svc := NewNotificationService(SMTPSender{})
_ = svc.SendWelcome("ada@example.com")
```

DIP is what makes the other four principles cash out. SRP gives you small types; OCP gives you interfaces; ISP keeps them focused; LSP keeps them honest. DIP is the wiring rule that says: at every layer boundary, depend on the abstraction, and let `main` glue concrete things together.

## Running the examples

Each principle has its own bare-bones runnable in the prototype:

```bash
cd prototype-code/solid-principles-go
go run ./cmd/srp
go run ./cmd/ocp
go run ./cmd/lsp
go run ./cmd/isp
go run ./cmd/dip
```

Every program is under 60 lines and exercises exactly one idea.

## A checklist for the next time you design a Go service

- **SRP** — Can I describe this struct's job in one sentence without using "and"?
- **OCP** — If I add a new variant, do I edit central logic or drop in a new file?
- **LSP** — Does every implementation of my interface actually behave the way callers assume?
- **ISP** — Is any caller forced to depend on methods it doesn't use?
- **DIP** — Does my business logic import a concrete DB / HTTP / SMTP client, or an interface?

SOLID isn't a set of rules to follow ceremoniously. It's friction-reduction: when adding a feature touches many files, breaks unrelated tests, or needs a new branch in an old switch — that's the design telling you which letter you skipped.
