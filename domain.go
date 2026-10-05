package main

import("fmt"; "errors"; "github.com/compilyator/VNTU_GoLang/labs/Lab2Data/socialindicator")

var ErrarInvalidCount = errors.New("invalid count: must be between 5 and 100")
var ErrarinvalidCountry = errors.New("invalid country: must be something in this line.")

type Options struct {
	Count int
	Country string
	Seed *int64
}

func (o Options) Validate() error {
	if o.Count < 5 || o.Count > 100 {
		return ErrarInvalidCount
	}
	if o.Country == "" {
		return ErrarinvalidCountry
	}
	return nil
}

type Result struct {
	Country string
	Count int
	Min float64
	Max float64
	Average float64
	Duplicate map[float64]int
}

func Analyze(o Options) (Result, error) {
	err := o.Validate()
	if err != nil {
		return Result{}, fmt.Errorf("validation error: %v", err)
	}

	var values []float64
	var errar error

	if o.Seed != nil {
		values, errar = socialindicator.GenerateWithSeed(o.Count, o.Country, *o.Seed)
	} else {
		values, errar = socialindicator.Generate(o.Count, o.Country)
	}

	if errar != nil {
		return Result{}, fmt.Errorf("generation error: %v", errar)
	}

	if len(values) == 0 {
		return Result{}, errors.New("no generation results")
	}
	minVal := values[0]
	maxVal := values[0]
	sum := 0.0
	CountsMap := make(map[float64]int)

	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sum += v
		CountsMap[v]++
	}
	var average float64 = sum / float64(len(values))

	return Result{
		Country: o.Country,
		Count: o.Count,
		Min: minVal,
		Max: maxVal,
		Average: average,
		Duplicate: CountsMap,
	}, nil
	}