package mes

import "errors"

const (
	TypeRegister = "Register"
	TypeText     = "Text"
	TypeMeta     = "Meta"
	TypeError    = "Error"
)

type Message struct {
	Type    string
	Payload any
}

type TextMessage struct {
	Sender string
	Text   string
}

func TextMessageFromMap(m map[string]interface{}) (*TextMessage, error) {
	tm := &TextMessage{}

	if sender, ok := m["Sender"].(string); ok {
		tm.Sender = sender
	} else {
		return nil, errors.New("InvalidSender")
	}

	if text, ok := m["Text"].(string); ok {
		tm.Text = text
	} else {
		return nil, errors.New("InvalidText")
	}

	return tm, nil
}

type RegisterMessage struct {
	Username string
	Room     string
}
