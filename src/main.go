package main

import (
	"fmt"
	"os"
	"sequel/main/services"
	"sequel/main/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// declaring lipgloss styles to be used
var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

var (
	db      = services.Connection{}
	db_name = "test"
)

type item struct {
	title       string
	description string
}

// defining the methods on the type of item. The thing after the brackets is the reciever
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleHelpMenu key.Binding
}

// will return our map struct with the needed bindings
func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list           list.Model
	input          textinput.Model
	show_input     bool
	selected_table list.Item
	keys           *listKeyMap
	delegateKeys   *delegateKeyMap
}

func newModel(db_name string, table_names []string) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	num_tables := len(table_names)
	items := make([]list.Item, num_tables)
	for i := 0; i < num_tables; i++ {
		items[i] = item{
			title:       table_names[i],
			description: "",
		}
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	table_list := list.New(items, delegate, 0, 0)
	table_list.Title = db_name
	table_list.Styles.Title = titleStyle
	table_list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleHelpMenu,
		}
	}
	table_list.DisableQuitKeybindings()

	return model{
		list:         table_list,
		input:        textinput.New(),
		show_input:   false,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.show_input {
	case true:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			switch {
			case msg.Type == tea.KeyEscape:
				m.show_input = false
				m.input.Reset()

			case msg.Type == tea.KeyEnter:
				m.show_input = false

				f, err := tea.LogToFile("./debug.log", "")
				if err != nil {
					panic("FUCK")
				}

				f.WriteString(m.input.Value())
				f.Close()

				m.input.Reset()
			}

			new_input_model, cmd := m.input.Update(msg)
			m.input = new_input_model
			cmds = append(cmds, cmd)
		}
	case false:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			// Don't match any of the keys below if we're actively filtering.
			if m.list.FilterState() == list.Filtering {
				break
			}

			switch {
			case key.Matches(msg, m.keys.toggleHelpMenu):
				m.list.SetShowHelp(!m.list.ShowHelp())
				return m, nil

			case msg.String() == "q":
				return m, tea.Quit
			}

		case utils.CustomMessage:
			switch msg.Message {
			case "drop table":
				m.show_input = true
				m.input.Focus()
			}
		}

		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
        m.selected_table = m.list.SelectedItem()
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	input := ""
	if m.show_input {
		input = "Drop table" + m.selected_table.FilterValue() + "? Y/N: " + m.input.View()
	}

	return appStyle.Render(m.list.View()) + input
}

func main() {
	db.CreateConnection("postgres", "localhost", 5432, "dev", "dev", db_name)
	tables, err := db.GetTables()
	if err != nil {
		fmt.Println("there was an error getting the tables: " + err.Error())
	}

	if _, err := tea.NewProgram(newModel(db_name, tables), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
