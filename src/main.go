package main

import (
	"fmt"
	"os"
	commandbar "sequel/main/components/CommandBar"
	tablelist "sequel/main/components/TableList"
	"sequel/main/services"
	// "sequel/main/utils"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
)

// var appStyle = lipgloss.NewStyle().Padding(1, 2)

var db services.Connection

// this is the go equivalent of an enum apparently
type ComponentId int

const (
	component_none ComponentId = iota
	component_tablelist
	component_textinput
)

type model struct {
	tablelist tablelist.Model
	command_bar            commandbar.Model
	focused_component      ComponentId
	prev_focused_component ComponentId
}

func newModel(db_name string, table_names []string) model {
	return model{
		tablelist: tablelist.NewModel(db_name, table_names),
		command_bar:       commandbar.NewModel(),
		focused_component: component_tablelist,
	}
}

func (m *model) setFocus(new_focus ComponentId) {
	m.prev_focused_component = m.focused_component
	m.focused_component = new_focus
	return
}

func (m *model) returnFocus() {
	m.focused_component = m.prev_focused_component
	m.prev_focused_component = component_none
	return
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.focused_component {
	case component_textinput:
		switch msg := msg.(type) {
		case commandbar.UserResponse:
			switch msg.Id {
			case "drop table":
				status_msg := "No Changes Made"
				if val := strings.ToUpper(msg.Value); val == "Y" || val == "YES" {
					db.DropTable(m.tablelist.Selected_table.FilterValue())
					status_msg = "Table Dropped Successfully"
				}
				m.command_bar = m.command_bar.Notify(status_msg).(commandbar.Model)
				break
			}
			m.returnFocus()

		default:
			new_input_model, cmd := m.command_bar.Update(msg)
			m.command_bar = new_input_model.(commandbar.Model)
			cmds = append(cmds, cmd)

		}

	case component_tablelist:
		switch msg := msg.(type) {

		case tablelist.TableAction:
			switch msg.Action {
			case "drop":
				m.command_bar = m.command_bar.Prompt(
					"drop table",
					"Are you sure you want to drop table "+m.tablelist.Selected_table.FilterValue()+"? Y/N",
				).(commandbar.Model)
				m.setFocus(component_textinput)
			}

		default:
			newListModel, cmd := m.tablelist.Update(msg)
			if updated_model, ok := newListModel.(tablelist.Model); ok {
				m.tablelist = updated_model
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.tablelist.View() + m.command_bar.View()
}

func main() {
	db_name := "test"
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
