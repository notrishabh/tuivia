package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/notrishabh/tuivia/quiz"
)

type model struct {
	currentques int
	form        *huh.Form
	questions   []quiz.QuizQuestion
}

func createGroups(questions []quiz.QuizQuestion) []*huh.Group {
	var groups []*huh.Group

	for _, q := range questions {
		group := huh.NewGroup(
			huh.NewSelect[string]().
				Key(string(q.Id)).
				Options(huh.NewOptions(q.AnswersArray...)...).
				Title(q.Question),
		)
		groups = append(groups, group)
	}
	return groups
}

func initialModel() model {
	questions, err := quiz.Quiz()
	if err != nil {
		log.Fatal(err)
	}
	return model{
		form:      huh.NewForm(createGroups(questions)...),
		questions: questions,
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
		for i, v := range m.questions {
			ans := m.form.GetString(string(v.Id))
			s += fmt.Sprintf("%d: %s\n", i+1, ans)
		}
		return s
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
