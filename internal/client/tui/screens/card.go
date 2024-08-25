package screens

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Sadere/gophkeeper/internal/client/tui/components"
	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	errNumberEmpty   = "Please enter card number"
	errExpMonthEmpty = "Please enter non-negative expiration month"
	errExpYearEmpty  = "Please enter non-negative expiration year"
	errCVVEmpty      = "Please enter 3 digits for CVV"
)

const (
	cardMetadata = iota
	cardNumber
	cardExpMonth
	cardExpYear
	cardCvv
)

type CardModel struct {
	state      *State
	cardID     uint64
	inputGroup components.InputGroup
	errorMsg   string
}

func validateNumber(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func NewCardModel(state *State, ID uint64) *CardModel {
	inputs := make([]textinput.Model, 5)

	inputs[cardMetadata] = textinput.New()
	inputs[cardMetadata].Placeholder = "Metadata"
	inputs[cardMetadata].Focus()
	inputs[cardMetadata].CharLimit = 100

	inputs[cardNumber] = textinput.New()
	inputs[cardNumber].Placeholder = "4012888888881881"
	inputs[cardNumber].CharLimit = 16
	inputs[cardNumber].Width = 30
	inputs[cardNumber].Validate = validateNumber

	inputs[cardExpMonth] = textinput.New()
	inputs[cardExpMonth].Placeholder = "MM"
	inputs[cardExpMonth].CharLimit = 2
	inputs[cardExpMonth].Width = 5
	inputs[cardExpMonth].Validate = validateNumber

	inputs[cardExpYear] = textinput.New()
	inputs[cardExpYear].Placeholder = "YY"
	inputs[cardExpYear].CharLimit = 2
	inputs[cardExpYear].Width = 5
	inputs[cardExpYear].Validate = validateNumber

	inputs[cardCvv] = textinput.New()
	inputs[cardCvv].Placeholder = "***"
	inputs[cardCvv].CharLimit = 3
	inputs[cardCvv].Width = 5
	inputs[cardCvv].Validate = validateNumber

	m := &CardModel{
		state:      state,
		cardID:     ID,
		inputGroup: components.NewInputGroup(inputs),
	}

	// Load secret if in view/edit mode
	if ID > 0 {
		secret, err := m.state.client.LoadSecret(context.Background(), ID)
		if err != nil {
			m.errorMsg = err.Error()
			return m
		}

		inputs[cardMetadata].SetValue(secret.Metadata)
		inputs[cardNumber].SetValue(secret.Card.Number)
		inputs[cardExpMonth].SetValue(strconv.Itoa(int(secret.Card.ExpMonth)))
		inputs[cardExpYear].SetValue(strconv.Itoa(int(secret.Card.ExpYear)))
		inputs[cardCvv].SetValue(strconv.Itoa(int(secret.Card.Cvv)))
	}

	return m
}

func (m CardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			mainScreen := NewSecretListModel(m.state)
			return NewRootModel(m.state).SwitchScreen(mainScreen)
		case "enter":
			if m.inputGroup.FocusIndex == m.inputGroup.InputNum {
				return m.Submit()
			}
		}
	}

	// Handle input group
	mm, cmd := m.inputGroup.Update(msg)

	m.inputGroup = mm.(components.InputGroup)

	return m, cmd
}

func (m CardModel) Submit() (tea.Model, tea.Cmd) {
	metadata := m.inputGroup.Inputs[cardMetadata].Value()
	number := m.inputGroup.Inputs[cardNumber].Value()
	expMonth, errMonth := strconv.Atoi(m.inputGroup.Inputs[cardExpMonth].Value())
	expYear, errYear := strconv.Atoi(m.inputGroup.Inputs[cardExpYear].Value())
	cvv, errCVV := strconv.Atoi(m.inputGroup.Inputs[cardCvv].Value())

	// Validate inputs
	if len(metadata) == 0 {
		m.errorMsg = errMetadataEmpty
		return m, nil
	}

	if len(number) == 0 {
		m.errorMsg = errNumberEmpty
		return m, nil
	}

	if expMonth <= 0 || errMonth != nil {
		m.errorMsg = errExpMonthEmpty
		return m, nil
	}

	if expYear <= 0 || errYear != nil {
		m.errorMsg = errExpYearEmpty
		return m, nil
	}

	if cvv <= 0 || errCVV != nil {
		m.errorMsg = errCVVEmpty
		return m, nil
	}

	// Save card
	err := m.state.client.SaveCard(
		context.Background(),
		m.cardID,
		metadata,
		number,
		uint32(expMonth),
		uint32(expYear),
		uint32(cvv),
	)
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	// Return to list
	mainScreen := NewSecretListModel(m.state)
	return NewRootModel(m.state).SwitchScreen(mainScreen)
}

func (m CardModel) View() string {

	cardForm := fmt.Sprintf(
		`
%s
%s

%s %s %s
%s %s %s
`,
		style.FocusedStyle.Width(30).Render("Card Number"),
		m.inputGroup.Inputs[cardNumber].View(),
		style.FocusedStyle.Width(8).Render("Exp MM"),
		style.FocusedStyle.Width(8).Render("Exp YY"),
		style.FocusedStyle.Width(6).Render("CVV"),
		m.inputGroup.Inputs[cardExpMonth].View(),
		m.inputGroup.Inputs[cardExpYear].View(),
		m.inputGroup.Inputs[cardCvv].View(),
	)

	// Button
	button := style.BlurredStyle.Render("[ Submit ]")
	if m.inputGroup.FocusIndex == m.inputGroup.InputNum {
		button = style.FocusedStyle.Render("[ Submit ]")
	}

	// Error
	errorBox := ""

	if len(m.errorMsg) > 0 {
		errorBox = style.ErrorStyle.Render(m.errorMsg)
	}

	body := lipgloss.JoinVertical(
		lipgloss.Top,
		"Please enter card details",
		m.inputGroup.Inputs[cardMetadata].View(),
		cardForm,
		button,
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		style.RenderBox(body),
		errorBox,
	)
}
