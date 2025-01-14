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
	form      *huh.Form
	questions []quiz.QuizQuestion
	end       bool
}

var selectedCategory string

func createGroups(questions []quiz.QuizQuestion) []*huh.Group {
	var groups []*huh.Group

	for _, q := range questions {
		group := huh.NewGroup(
			huh.NewSelect[string]().
				Description(fmt.Sprintf("[%s] %s [%s]", q.Category, q.Description, q.Difficulty)).
				Key(string(q.Id)).
				Options(huh.NewOptions(q.AnswersArray...)...).
				Title(q.Question).Validate(func(s string) error {
				if q.CorrectAnswer != s {
					return fmt.Errorf("Wrong answer. Select correct ans: %s", q.CorrectAnswer)
				}
				return nil
			}),
		)
		groups = append(groups, group)
	}
	return groups
}

func createCategoryGroup() *huh.Group {
	categories, err := quiz.GetCategories()
	if err != nil {
		log.Fatal(err)
	}
	var categoriesSlice []string
	for _, c := range categories {
		categoriesSlice = append(categoriesSlice, c.Name)
	}

	categoryGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Select Category").
			Options(huh.NewOptions(categoriesSlice...)...).Value(&selectedCategory),
	)

	return categoryGroup
}

func initialModel() model {
	questions, err := quiz.Quiz("all")
	if err != nil {
		log.Fatal(err)
	}
	return model{
		form:      huh.NewForm(createCategoryGroup()),
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
		if m.end {
			cmds = append(cmds, tea.Quit)
		} else {
			questions, err := quiz.Quiz(selectedCategory)
			if err != nil {
				log.Fatal(err)
			}
			m.questions = questions
			m.form = huh.NewForm(createGroups(questions)...)
			m.form.PrevGroup()
			m.end = true
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := "\nA simple tech quiz\n\n"

	if m.form.State == huh.StateCompleted {
		for i, v := range m.questions {
			ans := m.form.GetString(string(v.Id))
			s += fmt.Sprintf("Q%d: %s\n", i+1, v.Question)
			s += fmt.Sprintf("A: %s\n", ans)
			s += fmt.Sprintf("Explanation: %s\n\n", v.Explanation)
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
