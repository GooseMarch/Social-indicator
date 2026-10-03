package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/socialindicator"
)

type model struct {
	Whatselected int
	input        []textinput.Model
	mode         []string
	wherecursor  int
	count        int
	country      string
	selectedmode string
	errar        error
	submit       bool
	isLoading    bool
	values       []float64
}

type generateModel struct {
	genvalues []float64
	generrar  error
	values    []float64
}

func initialModel() model {
	input := make([]textinput.Model, 2)

	input[0] = textinput.New()
	input[0].Placeholder = "Enter a number between 5 and 100"
	input[0].Focus()

	input[1] = textinput.New()
	input[1].Placeholder = "Enter name of country"

	return model{
		Whatselected: 0,
		input:        input,
		mode:         []string{"default generation", "generation with seed"},
		wherecursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func generatecmd(count int, country string, mode string) tea.Cmd {
	return func() tea.Msg {
		var values []float64
		var errar error

		if mode == "generation with seed" {
			seed := int64(55)
			values, errar = socialindicator.GenerateWithSeed(count, country, seed)
		} else {
			values, errar = socialindicator.Generate(count, country)
		}

		return generateModel{
			genvalues: values,
			generrar:  errar,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit

		case "s", "S", "і", "І":
			m.Whatselected++

			if m.Whatselected > 2 {
				m.Whatselected = 0
			}

			cmds := make([]tea.Cmd, len(m.input))
			for i := 0; i < len(m.input); i++ {
				if i == m.Whatselected {
					cmds[i] = m.input[i].Focus()
				} else {
					m.input[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)

		case "w", "W", "ц", "Ц":
			m.Whatselected--

			if m.Whatselected < 0 {
				m.Whatselected = 2
			}

			cmds := make([]tea.Cmd, len(m.input))
			for i := 0; i < len(m.input); i++ {
				if i == m.Whatselected {
					cmds[i] = m.input[i].Focus()
				} else {
					m.input[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		case "enter":
			country := string(strings.TrimSpace(m.input[1].Value()))
			count, err := strconv.Atoi(strings.TrimSpace(m.input[0].Value()))
			if country == "" {
				m.errar = fmt.Errorf("Country name is required")
				return m, nil
			} else if err != nil || count < 5 || count > 100 {
				m.errar = fmt.Errorf("Count must be a number between 5 and 100")
				return m, nil
			}

			m.country = country
			m.count = count
			m.selectedmode = m.mode[m.wherecursor]
			m.submit = true
			m.isLoading = true

			return m, generatecmd(m.count, m.country, m.selectedmode)

		case "up":
			if m.Whatselected == 2 && m.wherecursor > 0 {
				m.wherecursor--
			}
		case "down":
			if m.Whatselected == 2 && m.wherecursor < len(m.mode)-1 {
				m.wherecursor++
			}

		}

		if m.Whatselected < len(m.input) {
			var cmd tea.Cmd
			m.input[m.Whatselected], cmd = m.input[m.Whatselected].Update(msg)
			return m, cmd
		}

	case generateModel:
		m.isLoading = false
		if msg.generrar != nil {
			m.errar = msg.generrar
			return m, nil
		}
		m.values = msg.genvalues
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	b := strings.Builder{}
	b.WriteString("=== generation settings ===\n\n\n\n")

	// 1. count
	if m.Whatselected == 0 {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("Count (5..100):\n")
	b.WriteString(" ")
	b.WriteString(m.input[0].View())
	b.WriteString("\n\n")

	// 2. country
	if m.Whatselected == 1 {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("country:\n")
	b.WriteString(" ")
	b.WriteString(m.input[1].View())
	b.WriteString("\n\n")

	// 3. generation mode
	if m.Whatselected == 2 {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("Generation Mode:\n")

	for i, mode := range m.mode {
		cursor := " "
		if m.wherecursor == i {
			cursor = "●"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, mode))
	}

	if m.errar != nil {
		b.WriteString(fmt.Sprintf("\n\033[31m%s\033[0m\n", m.errar.Error()))
	}

	b.WriteString("\n[W/S] - change tabs\n")
	b.WriteString("[↑ / ↓] - change generation\n")
	b.WriteString("[Enter] - start all\n")

	return b.String()

}
func main() {
	p := tea.NewProgram(initialModel())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	m, ok := finalModel.(model)
	if !ok || len(m.values) == 0 {
		return
	}

	minval := m.values[0]
	maxval := m.values[0]
	sum := 0.0

	for _, v := range m.values {
		if v < minval {
			minval = v
		}
		if v > maxval {
			maxval = v
		}
		sum += v
	}
	var average float64 = sum / float64(len(m.values))

	fmt.Printf("\n\n\033[31m===  program for analyzing demographic indicators of country  ===\033[0m\n")

	fmt.Printf("\n\033[36m===                     generation results                    ===\033[0m\n\n")
	fmt.Printf("\033[37mmin: \033[0m\033[34m %.2f years\033[0m, \n\033[37mmax: \033[0m\033[94m %.2f years\033[0m, \n\033[37maverage: \033[0m\033[33m%.2f years\033[0m \n", minval, maxval, average)

	// reapeter
	countsMap := make(map[float64]int)
	for _, v := range m.values {
		countsMap[v]++
	}

	fmt.Print("\n\033[35mRepeated values (map)\033[0m\n")
	for k, v := range countsMap {
		fmt.Printf("\033[37mValue \033[0m\033[32m%.2f years\033[0m\033[37m is found \033[0m\033[93m%d \033[0m\033[37mtimes\033[0m\n", k, v)
	}

}
