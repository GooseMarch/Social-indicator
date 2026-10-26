package cli

import (
	"flag"
	"fmt"
	"os"

	"lab2_bubble/domain"
	"lab2_bubble/filework"
	"lab2_bubble/htmlreport"
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
	extra := Flagset.Bool("extra", false , "Extra to do")
	mode := Flagset.String("mode", "online", "Choose a mode: online or offline")
	save := Flagset.String("save", "", "Insert path to save JSON")
	input := Flagset.String("input", "", "Insert path to input JSON (with offline mode)")
	reportPath := Flagset.String("report","","Insert path to save html report")

	err := Flagset.Parse(args)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	if Flagset.NArg() > 0{
		fmt.Println("Error: check if you doesnt have a unnecessary arguments")
		os.Exit(2)
	}

	if *mode != "online" && *mode != "offline"{
		fmt.Printf( "Error: unknown mode %s. Must to be online or offline", *mode)
		os.Exit(1)
	}

	SeedExist := false
	Flagset.Visit(func(f *flag.Flag){
		if f.Name == "seed"{
			SeedExist = true
		}
	})

	var dataset domain.Dataset

	if *mode == "online"{
		if *count == 0 || *country == ""{
			fmt.Println( "Error: count or/and country is not writed in parametrs in online mode")
			os.Exit(2)
		}
		if *input != ""{
			fmt.Println( "Error: input doesn`t need in online mode")
			os.Exit(2)
		}
	
		options := domain.Options{
			Count:   *count,
			Country: *country,
			Extra: *extra,
			Mode: *mode,
			Save: *save,
			Report: *reportPath,
		}

		if SeedExist == true{
			seedVal := *seed
			options.Seed = &seedVal
		}

		data, errar := domain.OnlineGenerate(options)
		if errar != nil{
				fmt.Printf("Error generation is failing: %s", errar)
				os.Exit(3)
		}
		dataset = data

		if *save != "" {
			if *save == *reportPath {
				fmt.Println("Error: save and report cant have the same path")
				os.Exit(3)
			}
			errar := filework.SaveDataset(*save, dataset); 
			if errar != nil {
				fmt.Printf("Error saving JSON: %v\n", errar)
				os.Exit(4)
			}
		}


	} else{
		if *save != "" || *count != 0 || *country != "" || SeedExist == true{
			fmt.Println("Error: count country save and seed don`t use in offline mode")
			os.Exit(5)
		}

		inputPath := *input
		if inputPath == "" {
			inputPath = "dataset.json"
		}
		data, errar := filework.LoadDataset(inputPath)
		if errar != nil{
			fmt.Printf("Error to load dataset: %v", errar)
			os.Exit(6)
		}
		dataset = data
		dataset.Mode = "offline"
	}
	result, errar := domain.Analyze(dataset.Values, dataset.Country, dataset.Time, dataset.Mode, dataset.Generator, dataset.Seed, dataset.Whatusing, dataset.SeedValue, dataset.Count, *extra)
	if errar != nil{
		fmt.Printf("Error cant analyze data try again: %v", errar)
		os.Exit(7)
	}
	if *reportPath != ""{
		html, errar := htmlreport.Generatehtml(result)
		if errar != nil {
			fmt.Printf("Error generating html can`t be done: %v", errar)
			os.Exit(8)
		}
		errar = filework.SaveReport(*reportPath, html)
		if errar != nil{
			if *save != ""{
				fmt.Printf("Error: JSON dataset save in %s, but html report failed: %v", *save, errar)
			}else{
				fmt.Printf("Error saving html: %v", errar)
			}
		}
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
	if result.Seedresult == true{
	fmt.Println("\n\nNumber of values above average:", result.Extra)
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
	fmt.Println("  -mode    Execution mode('online' or 'offline'. using 'online' by default)")
	fmt.Println("               Online mode                      ")
	fmt.Println("  -count   Number of values (5..100)")
	fmt.Println("  -country Name of the country")
	fmt.Println("  -seed    Random generator seed (optional)")
	fmt.Println("  -extra   Count values above average (optional)")
	fmt.Println("  -save    Path to save new JSON dataset file. (for saving file you need in end of path add filename.json)")
	fmt.Println("  -report  Path to save generated HTML report")
	fmt.Println("               Offline mode                     ")
	fmt.Println("  -input   Path to input JSON dataset file (offline mode, default: dataset.json)")
	fmt.Println("  -report  Path to save generated HTML report")
	
}