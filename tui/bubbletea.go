package tui

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/notrishabh/tuivia/quiz"
)

const maxWidth = 80

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	Highlight,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

type model struct {
	form      *huh.Form
	questions []quiz.QuizQuestion
	end       bool
	width     int
	styles    *Styles
	lg        *lipgloss.Renderer
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
	m := model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)
	m.form = huh.NewForm(createCategoryGroup()).WithShowHelp(false).WithShowErrors(false)

	return m
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
			m.form = huh.NewForm(createGroups(questions)...).WithShowHelp(false)
			m.form.PrevGroup()
			m.end = true
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	st := m.styles

	v := strings.TrimSuffix(m.form.View(), "\n\n")
	form := m.lg.NewStyle().Margin(1, 0).Render(v)

	header := m.appBoundaryView(fmt.Sprintf("Huh? you know %s?", selectedCategory))
	body := lipgloss.JoinHorizontal(lipgloss.Top, form)
	footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))

	if m.form.State == huh.StateCompleted {
		s := ""
		for i, v := range m.questions {
			ans := m.form.GetString(string(v.Id))
			s += fmt.Sprintf("Q%d: %s\n", i+1, st.Highlight.Render(v.Question))
			s += fmt.Sprintf("A: %s\n\n", ans)
			if v.Explanation != "" {
				s += fmt.Sprintf("Explanation: %s\n\n\n", v.Explanation)
			}
		}
		return st.Status.Margin(0, 1).Padding(1, 2).Width(80).Render(s) + "\n\n"
	}

	return st.Base.Render(header + "\n" + body + "\n\n" + footer)
}

func (m model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func RunTui() {
	_, err := tea.NewProgram(initialModel()).Run()

	if err != nil {
		fmt.Printf("Error encountered: %v", err)
		os.Exit(1)
	}
}
