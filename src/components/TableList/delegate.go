package tablelist

import (
	"sequel/main/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type TableAction struct {
	Action     string
	Table_name string
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		// this is type assertion, so i will be our item (if the type is of item and therefore it exists)
		i, ok := m.SelectedItem().(item)
		if !ok {
			return nil
		}

		title = i.Title()
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(statusMessageStyle("Table " + title + " selected!"))

			case key.Matches(msg, keys.drop):

				new_cmd := TableAction {
                    Action: "drop",
                    Table_name: title,
                }

				return utils.MakeCustomCommand(new_cmd)
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.drop}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	drop   key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.drop,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.drop,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		drop: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "drop table"),
		),
	}
}
