package commandbar

import (
	"sequel/main/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CommandBarAction struct {
	Action  string
	Message string
}

type UserResponse struct {
	Id    string
	Value string
}

type Model struct {
	req_id   string
	is_input bool
	message  string
	input    textinput.Model
}

func NewModel() Model {
	return Model{
		input: textinput.New(),
	}
}

func (m Model) Notify(notification string) tea.Model {
	new_model, _ := m.Update(CommandBarAction{
		Action:  "notify",
		Message: notification,
	})

	return new_model
}

func (m Model) Prompt(id string, prompt string) tea.Model {
	m.input.Focus()
	m.req_id = id
	new_model, _ := m.Update(CommandBarAction{
		Action:  "prompt",
		Message: prompt,
	})
	return new_model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case CommandBarAction:
		text_to_display := msg.Message
		switch msg.Action {
		case "notify":
			m.is_input = false
			m.message = text_to_display

		case "prompt":
			m.is_input = true
			m.message = text_to_display
		}
	case tea.KeyMsg:
		key_type := msg.Type

		switch key_type {
		case tea.KeyEnter:
			payload := UserResponse{
				Id:    m.req_id,
				Value: m.input.Value(),
			}
			cmds = append(cmds, utils.MakeCustomCommand(payload))

		case tea.KeyEscape:
			// with no id, it will just lose focus
			payload := UserResponse{
				Id:    "",
				Value: "",
			}

            // resetting the field for render
			m.input.Reset()
			m.message = ""
			m.is_input = false

			cmds = append(cmds, utils.MakeCustomCommand(payload))

		default:
			new_input_model, cmd := m.input.Update(msg)
			m.input = new_input_model
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	view := m.message
	if m.is_input {
		view = m.message + " " + m.input.View()
	}
	return view
}
