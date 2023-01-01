package eform

import (
	"encoding/json"
	"fmt"
)

// Modal represents a modal form. These forms have a body with text and two buttons at the end, typically one
// for Yes and one for No. These buttons may have custom text, but can, unlike with a Menu form, not have
// images next to them.
type Modal struct {
	title, body string
	btn1, btn2  buttonData
	onClose     Handler
}

// NewModal creates a new Modal form. The
// title passed is formatted following the fmt.Sprintln rules.
// Default buttons are 'yes' and 'no'
func NewModal(title ...any) Modal {
	m := Modal{
		title: format(title),
		btn1:  buttonData{YesButton(), nil},
		btn2:  buttonData{NoButton(), nil},
	}
	return m
}

// YesButton returns a Button which may be used as a default 'yes' button for a modal form.
func YesButton() Button {
	return Button{Text: "gui.yes"}
}

// NoButton returns a Button which may be used as a default 'no' button for a modal form.
func NoButton() Button {
	return Button{Text: "gui.no"}
}

// MarshalJSON ...
func (m Modal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "modal",
		"title":   m.title,
		"content": m.body,
		"button1": m.Buttons()[0].Text,
		"button2": m.Buttons()[1].Text,
	})
}

// WithBody creates a copy of the Modal form and changes its body to the body passed, after which the new Modal
// form is returned. The text is formatted following the rules of fmt.Sprintln.
func (m Modal) WithBody(body ...any) Modal {
	m.body = format(body)
	return m
}

// WithButton1 creates a copy of the Modal form and change the button 1 to the button passed.
func (m Modal) WithButton1(btn Button, onClick Handler) Modal {
	m.btn1 = buttonData{btn, onClick}
	return m
}

// WithButton2 creates a copy of the Modal form and change the button 2 to the button passed.
func (m Modal) WithButton2(btn Button, onClick Handler) Modal {
	m.btn2 = buttonData{btn, onClick}
	return m
}

// Title returns the formatted title passed to the menu upon construction using NewModal().
func (m Modal) Title() string {
	return m.title
}

// Body returns the formatted text in the body passed to the menu using WithBody().
func (m Modal) Body() string {
	return m.body
}

// SubmitJSON submits a JSON byte slice to the modal form. This byte slice contains a JSON encoded bool in it,
// which is used to determine which button was clicked.
func (m Modal) SubmitJSON(b []byte, submitter Submitter) error {
	if b == nil {
		m.onClose.Call(submitter)
		return nil
	}

	var value bool
	if err := json.Unmarshal(b, &value); err != nil {
		return fmt.Errorf("error parsing JSON as bool: %w", err)
	}
	if value {
		m.btn1.onClick.Call(submitter)
		return nil
	}
	m.btn2.onClick.Call(submitter)
	return nil
}

// Buttons returns a list of all buttons of the Modal form, which will always be a total of two buttons.
func (m Modal) Buttons() []Button {
	return []Button{m.btn1.btn, m.btn2.btn}
}

func (Modal) __() {}
