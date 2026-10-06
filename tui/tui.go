package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"lab2_bubble/domain"
	"lab2_bubble/cli"


	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Whatselected int
	input        []textinput.Model
	wherecursor  int
	count        int
	country      string
	errar        error
	submit       bool
	seed         bool
	seedinput    textinput.Model
	finalresult  domain.Result
}

func initialModel() model {
	input := make([]textinput.Model, 2)

	input[0] = textinput.New()
	input[0].Placeholder = "Enter a number between 5 and 100"
	input[0].CharLimit = 3
	input[0].Width = 20
	input[0].Focus()

	input[1] = textinput.New()
	input[1].Placeholder = "Enter name of country"
	input[1].CharLimit = 50
	input[1].Width = 20

	seedinput := textinput.New()
	seedinput.Placeholder = "Enter seed number (recommended 55)"
	seedinput.CharLimit = 20
	seedinput.Width = 20

	return model{
		Whatselected: 0,
		input:        input,
		seedinput:    seedinput,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit

		case "s", "S", "і", "І":
			m.Whatselected++
			if m.Whatselected > 4 {
				m.Whatselected = 0
			}
			for i := 0; i < len(m.input); i++ {
				if i == m.Whatselected {
					m.input[i].Focus()
				} else {
					m.input[i].Blur()
				}
				if m.Whatselected == 2 && m.seed {
					return m, m.seedinput.Focus()
				}
			}
			m.seedinput.Blur()
			return m, nil

		case "w", "W", "ц", "Ц":
			m.Whatselected--
			if m.Whatselected < 0 {
				m.Whatselected = 4
			}
			for i := 0; i < len(m.input); i++ {
				if i == m.Whatselected {
					m.input[i].Focus()
				} else {
					m.input[i].Blur()
				}
				if m.Whatselected == 2 && m.seed {
					return m, m.seedinput.Focus()
				}
			}
			m.seedinput.Blur()

		case "p", "P", "з", "З":
			if m.Whatselected == 2 {
				m.seed = !m.seed
			}

		case "enter":
			count, err := strconv.Atoi(strings.TrimSpace(m.input[0].Value()))
			if err != nil || count < 5 || count > 100 {
				m.errar = fmt.Errorf("Count must be a number between 5 and 100")
				return m, nil
			}

			country := strings.TrimSpace(m.input[1].Value())
			if country == "" {
				m.errar = fmt.Errorf("Country name is required")
				return m, nil
			}

			options := domain.Options{
				Count:   count,
				Country: country,
			}

			if m.seed {
				seed, err := strconv.ParseInt(strings.TrimSpace(m.seedinput.Value()), 10, 64)
				if err != nil {
					m.errar = fmt.Errorf("Invalid seed value")
					return m, nil
				}
				options.Seed = &seed
			}

			result, err := domain.Analyze(options)
			if err != nil {
				m.errar = err
				return m, nil
			}

			m.country = country
			m.count = count
			m.finalresult = result
			m.submit = true

			return m, tea.Quit
		}

		var cmd tea.Cmd
		if m.Whatselected < 2 {
			m.input[m.Whatselected], cmd = m.input[m.Whatselected].Update(msg)
		}
		if m.Whatselected == 2 && m.seed {
			m.seedinput, cmd = m.seedinput.Update(msg)
		}
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	b := strings.Builder{}
	b.WriteString("=== generation settings ===\n\n\n\n")

	if m.Whatselected == 0 {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("Count (5..100):\n ")
	b.WriteString(m.input[0].View())
	b.WriteString("\n\n")

	if m.Whatselected == 1 {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("country:\n ")
	b.WriteString(m.input[1].View())
	b.WriteString("\n\n")

	checkbox := ""
	if m.seed {
		checkbox = "x"
	}
	if m.wherecursor == 2 {
		b.WriteString(checkbox + "Enable Fixed Seed (press p)")
	} else {
		b.WriteString("Enable Fixed Seed (press p)")
	}

	if m.seed {
		b.WriteString("\n\nSeed value:\n ")
		b.WriteString(m.seedinput.View())
	}

	if m.errar != nil {
		b.WriteString(fmt.Sprintf("\n\033[31m %s \033[0m\n", m.errar.Error()))
	}

	b.WriteString("\n[W/S] - change tabs\n")
	b.WriteString("[Enter] - start all\n")
	b.WriteString("[Esc] - exiting program\n")

	return b.String()
}

func Run() {
	p := tea.NewProgram(initialModel())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	m, ok := finalModel.(model)
	if !ok {
		return
	}

	if m.submit && m.errar == nil {
		cli.PrintResult(m.finalresult)
	}
}	