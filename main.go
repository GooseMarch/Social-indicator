package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	//"github.com/compilyator/VNTU_GoLang/labs/Lab2Data/socialindicator"
)

const appversion = "1.0"


func main() {
	if len(os.Args) < 2 {
		TUI()
		return
	}
	commands := os.Args[1]

	switch commands {
		case "analyze":
			 analyzeCLI(os.Args[2:])
		case "help":
			printHelp()
		case "version":
			fmt.Printf("Current version: %s\n", appversion)
		default:
			fmt.Printf("Unknown command: %s\n", commands)
			printHelp()
			os.Exit(1)
	}
}


// CLI mode
func analyzeCLI(arguments []string) {
	flagSet := flag.NewFlagSet("analyze", flag.ExitOnError)
	count := flagSet.Int("count", 0, "Number of values to generate (between 5 and 100)")
	country := flagSet.String("country","","Country name")

	var seed int64
	var seedIsSet bool

	Seedcheck := func(str string) error {
		var value int64
		_, errar := fmt.Sscan(str, &value)
		if errar != nil {
			return fmt.Errorf("invalid seed value: %v", errar)
		}
		seed = value
		seedIsSet = true
		return nil
}	

	flagSet.Func("seed", "Random generator seed" , Seedcheck)
		
		errar := flagSet.Parse(arguments)
		if errar != nil {
			os.Exit(1)
		}

		options := Options{
			Count: *count,
			Country: *country,
		}
		if seedIsSet {
			options.Seed = &seed
		}

		result, errar := Analyze(options)
		if errar != nil {
			fmt.Printf("Error: %v\n", errar)
			os.Exit(1)
		}

		printResult(result)

}

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
	seed         bool
	seedinput    textinput.Model
	values       []float64
	finalresult  Result
}

/*type generateModel struct {
	genvalues []float64
	generrar  error
	values    []float64
}*/

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
	seedinput.Placeholder = "Enter seed number (optional - recommended 55)"
	seedinput.CharLimit = 20
	seedinput.Width = 20

	return model{
		Whatselected: 0,
		input:        input,
		//mode:         []string{"default generation", "generation with seed"},
		//wherecursor:  0,
		seedinput: seedinput,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

/*func generatecmd(count int, country string, mode string) tea.Cmd {
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
}*/

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

			//cmds := make([]tea.Cmd, len(m.input))
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
			//return m, tea.Batch(cmds...)

		case "w", "W", "ц", "Ц":
			m.Whatselected--

			if m.Whatselected < 0 {
				m.Whatselected = 4
			}

			//cmds := make([]tea.Cmd, len(m.input))
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
			//return m, tea.Batch(cmds...)
			
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

			options := Options{
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

			result, err := Analyze(options)
			if err != nil {
				m.errar = err
				return m, nil
			}

			

			m.country = country
			m.count = count
			m.finalresult = result
			m.submit = true

			

			return m, tea.Quit

		/*case "up":
			if m.Whatselected == 2 && m.wherecursor > 0 {
				m.wherecursor--
			}
		case "down":
			if m.Whatselected == 2 && m.wherecursor < len(m.mode)-1 {
				m.wherecursor++
			}
			*/
		


	/*case generateModel:
		m.isLoading = false
		if msg.generrar != nil {
			m.errar = msg.generrar
			return m, nil
		}
		m.values = msg.genvalues
		return m, tea.Quit
	}

	return m, nil
}*/

		/*result, err := Analyze(options)
		if err != nil {
			m.errar = err
			return m, nil
		}
		m.finalresult = result
		return m, tea.Quit
		}*/

		

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

	/*//3. generation mode
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
	}*/

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
		b.WriteString("\n\n")
		b.WriteString("Seed value:\n")
		b.WriteString(" ")
		b.WriteString(m.seedinput.View())
	}
	if m.errar != nil {
		b.WriteString(fmt.Sprintf("\n\033[31m %s \033[0m\n", m.errar.Error()))
	}

	b.WriteString("\n[W/S] - change tabs\n")
	//b.WriteString("[↑ / ↓] - change generation\n")
	b.WriteString("[Enter] - start all\n")
	b.WriteString("[Esc] - exiting program\n")

	return b.String()

}

func TUI() {
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
		printResult(m.finalresult)
	}
}


func printResult(result Result) {

	fmt.Printf("\n\n\033[31m===  program for analyzing demographic indicators of country  ===\033[0m\n")

	fmt.Printf("\n\033[36m===                     generation results                    ===\033[0m\n\n")
	fmt.Printf("\033[37mmin: \033[0m\033[34m %.2f years\033[0m, \n\033[37mmax: \033[0m\033[94m %.2f years\033[0m, \n\033[37maverage: \033[0m\033[33m%.2f years\033[0m \n", result.Min, result.Max, result.Average)

	fmt.Print("\n\033[35mRepeated values (map)\033[0m\n")
	for k, v := range result.Duplicate {
		fmt.Printf("\033[37mValue \033[0m\033[32m%.2f years\033[0m\033[37m is found \033[0m\033[93m%d \033[0m\033[37mtimes\033[0m\n", k, v)
	}

}


func printHelp() {
	fmt.Println("Usage: go run main.go [command] [flag country] [flag value] [seed value]")
	fmt.Println("If run without 'flag' it will run in TUI mode")
	fmt.Println("Commands:")
	fmt.Println("  analyze: Analyze demographic indicators for a country\n")
	fmt.Println("  help: Show this help message\n")
	fmt.Println("  version: Show the current version of the program\n")
	fmt.Println("\nFlags for 'analyze' command:")
	fmt.Println("  -count: Number of values to generate (between 5 and 100)\n")
	fmt.Println("  -country: Name of the country\n")
	fmt.Println("  -seed: Random generator seed (optional)\n")


}


/*func main() {
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

}*/
