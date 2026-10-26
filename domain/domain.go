package domain

import("fmt"; "errors"; "time"; "github.com/compilyator/VNTU_GoLang/labs/Lab2Data/socialindicator")

var ErrarInvalidCount = errors.New("invalid count: must be between 5 and 100")
var ErrarinvalidCountry = errors.New("invalid country: must be something in this line.")
var ErrarInvalidjson = errors.New("invalid json format")

type Dataset struct {
	Country string `json:"country"`
	Count int `json:"count"`
	Version string `json:"version"`
	Values []float64 `json:"values"`
	Generator string `json:"generator"`
	Time string `json:"time"`
	Seed bool `json:"seed"`
	Whatusing string `json:"whatusing"`
	SeedValue int64 `json:"seed_value"`
	Mode string `json:"mode"`
}

type Options struct {
	Count int
	Country string
	Seed *int64
	Extra bool
	Mode string
	Save string
	Input string
	Report string
}

type valuewithpointer struct {
	Value float64
	Pointer int
}
type Result struct {
	Country string
	Count int
	Min float64
	Max float64
	Sum float64
	Average float64
	Duplicate map[float64]int
	Extra int64
	Values []valuewithpointer
	Time string
	Reporttime string
	Mode string
	Generator string
	Seedresult bool
	Seed bool
	Whatusing string
	Seedvalue int64
}

	func Analyze(values []float64, country string, dataTime string, mode string, generator string, seed bool, whatUsing string, seedValue int64, count int, extra bool) (Result, error) {
		if len(values) == 0 {
			return Result{}, errors.New("no generation results. Please write any values in json file")
		}
	minVal := values[0]
	maxVal := values[0]
	sum := 0.0
	countsMap := make(map[float64]int)
	pointerValues := make([]valuewithpointer, len(values))

	for i, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sum += v
		countsMap[v]++
		pointerValues[i] = valuewithpointer{Value: v, Pointer: i}
	}
	var average float64 = sum / float64(len(values))

	var extracount int64 = 0

	
	if extra {
			for _, v := range values{
		if (average < v){
			extracount++
		}
	}
	}

	return Result{
		Country: country,
		Count: count,
		Min: minVal,
		Max: maxVal,
		Average: average,
		Duplicate: countsMap,
		Extra: extracount,
		Values: pointerValues,
		Time: dataTime,
		Mode: mode,
		Generator: generator,
		Seedvalue: seedValue,
		Whatusing: whatUsing,
		Seed: seed,
		Sum: sum,
		Seedresult: extra,
		Reporttime: time.Now().Format(time.RFC3339),
	}, nil
	}
	func OnlineGenerate(o Options) (Dataset, error) {
		if o.Count < 5 || o.Count > 100 {
		return Dataset{}, ErrarInvalidCount
	}
	if o.Country == "" {
		return Dataset{}, ErrarinvalidCountry
	}
	
	var values []float64
	var errar error
	var seedVal int64
	seedb := false

	if o.Seed != nil {
		seedVal = *o.Seed
		seedb = true
		values, errar = socialindicator.GenerateWithSeed(o.Count, o.Country, seedVal)
	} else {
		values, errar = socialindicator.Generate(o.Count, o.Country)
	}
	

	if errar != nil {
		return Dataset{}, fmt.Errorf("generation error: %v", errar)
	}

	return Dataset{
		Version: "1.0",
		Country: o.Country,
		Count: o.Count,
		Values: values,
		SeedValue: seedVal,
		Seed: seedb,
		Time: time.Now().Format(time.RFC3339),
		Generator: "socialindicator.Generate",
		Whatusing: "years",
		Mode: "online",
	}, nil
}
