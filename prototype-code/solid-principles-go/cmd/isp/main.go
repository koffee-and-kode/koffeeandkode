// ISP — Interface Segregation Principle.
//
// Each function depends only on the methods it actually uses. A fat
// "MessageStore" interface would force every implementation (including
// test fakes) to satisfy methods nobody calls.
package main

import "fmt"

type Message struct {
	ID   int
	To   string
	Body string
}

// Used by the sender path: only needs to record outgoing messages.
type MessageSaver interface {
	Save(m Message) error
}

// Used by the inbox/audit path: only needs to read them back.
type MessageReader interface {
	Get(id int) (Message, error)
}

// RecordSent depends only on Save — never sees Get.
func RecordSent(s MessageSaver, m Message) error {
	return s.Save(m)
}

// ViewSent depends only on Get — never sees Save.
func ViewSent(r MessageReader, id int) (Message, error) {
	return r.Get(id)
}

// One concrete type can satisfy multiple small interfaces — but each caller
// only sees the slice it actually needs.
type memStore struct {
	messages map[int]Message
}

func (m *memStore) Save(msg Message) error {
	m.messages[msg.ID] = msg
	return nil
}

func (m *memStore) Get(id int) (Message, error) {
	return m.messages[id], nil
}

func main() {
	store := &memStore{messages: map[int]Message{}}
	_ = RecordSent(store, Message{ID: 1, To: "ada@example.com", Body: "Welcome!"})
	m, _ := ViewSent(store, 1)
	fmt.Printf("retrieved: %+v\n", m)
}
