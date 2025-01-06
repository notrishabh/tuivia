package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	currentques int
	form        *huh.Form
}

func initialModel() model {
	return model{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().Key("name").Options(huh.NewOptions("yushu", "ippi", "smol")...).Title("Choose your name"),

				huh.NewSelect[int]().Key("level").Options(huh.NewOptions(1, 2, 999)...).Title("Choose your level").Description("This will determine your level"),
			),
		),
	}
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := "\nThis is a simple tech quiz\n\n"

	if m.form.State == huh.StateCompleted {
		name := m.form.GetString("name")
		level := m.form.GetInt("level")
		return fmt.Sprintf("You selected: %s, Lvl. %d", name, level)
	}
	q := "\n\nPress q to quit.\n"

	return s + m.form.View() + q
}

func RunTui() {
	_, err := tea.NewProgram(initialModel()).Run()

	if err != nil {
		fmt.Printf("Error boi: %v", err)
		os.Exit(1)
	}
}
