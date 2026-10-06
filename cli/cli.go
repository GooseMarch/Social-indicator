package cli

import (
	"flag"
	"fmt"
	"os"

	"lab2_bubble/domain"
)

const version = "1.0"

func RunCLI(args []string) {
	if len(args) == 0 {
		PrintHelp()
		return
	}

	command := args[0]
	switch command {
	case "analyze":
		AnalyzeCLI(args[1:])
	case "help":
		PrintHelp()
	case "version":
		fmt.Printf("Current version: %s\n", version)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		PrintHelp()
		os.Exit(1)
	}
}

func AnalyzeCLI(args []string) {
	Flagset := flag.NewFlagSet("analyze", flag.ExitOnError)
	count := Flagset.Int("count", 0, "Number of values (5..100)")
	country := Flagset.String("country", "", "Country name")
	seed := Flagset.Int64("seed", -1, "Random seed (optional)")

	err := Flagset.Parse(args)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	options := domain.Options{
		Count:   *count,
		Country: *country,
	}

	if *seed != -1 {
		options.Seed = seed
	}

	result, err := domain.Analyze(options)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	PrintResult(result)
}

func PrintResult(result domain.Result) {
	fmt.Printf("\n\n\033[31m===  program for analyzing demographic indicators of country  ===\033[0m\n")
	fmt.Printf("\n\033[36m===                     generation results                    ===\033[0m\n\n")
	fmt.Printf("\033[37mmin: \033[0m\033[34m %.2f years\033[0m, \n\033[37mmax: \033[0m\033[94m %.2f years\033[0m, \n\033[37maverage: \033[0m\033[33m%.2f years\033[0m \n", result.Min, result.Max, result.Average)

	fmt.Print("\n\033[35mRepeated values (map)\033[0m\n")
	for k, v := range result.Duplicate {
		fmt.Printf("\033[37mValue \033[0m\033[32m%.2f years\033[0m\033[37m is found \033[0m\033[93m%d \033[0m\033[37mtimes\033[0m\n", k, v)
	}
}

func PrintHelp() {
	fmt.Println("Usage: go run main.go [command] [flags]")
	fmt.Println("If run without commands, it will launch in TUI mode.")
	fmt.Println("\nCommands:")
	fmt.Println("  analyze  Analyze demographic indicators for a country")
	fmt.Println("  help     Show this help message")
	fmt.Println("  version  Show the current version")
	fmt.Println("\nFlags for 'analyze':")
	fmt.Println("  -count   Number of values (5..100)")
	fmt.Println("  -country Name of the country")
	fmt.Println("  -seed    Random generator seed (optional)")
}