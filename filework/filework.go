package filework

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"time"
	"strings"
	"lab2_bubble/domain"
)

const maxfilesize = 1048576

func SaveDataset(path string, dataset domain.Dataset) error {
    directory := filepath.Dir(path)
	if directory != "." && directory != "" {
		_, errar := os.Stat(directory)
		if os.IsNotExist(errar) {
			return fmt.Errorf("directory does not exist: %s", directory)
}
	}
	file, errar := os.OpenFile(path, os.O_CREATE| os.O_WRONLY|os.O_EXCL, 0604)
	if errar != nil {
		if os.IsExist(errar) {
			return fmt.Errorf("file already exists: %s", path)
		}
		return fmt.Errorf("failed to create file: %v", errar)
	}
	defer file.Close()

	data, errar := json.MarshalIndent(dataset, "", "  ")
	if errar != nil {
		return errar
	}
	_, errarwrite := file.Write(data)
	errarclose := file.Close()
	if errarwrite != nil {
		_ = os.Remove(path)
		return fmt.Errorf("failed to write to file: %v, close error: %v", errarwrite, errarclose)
	}

	return nil
}


func LoadDataset(path string) (domain.Dataset, error) {
	file, errar := os.Open(path)
	if errar != nil {
		return domain.Dataset{}, fmt.Errorf("failed to open file: %v\n", errar)
	}
	defer file.Close()

	limitread := io.LimitReader(file, maxfilesize)
	data, errar := io.ReadAll(limitread)
	if errar != nil {
		return domain.Dataset{}, fmt.Errorf("failed to read file: %v\n", errar)
	}

	if int64(len(data)) >= maxfilesize {
		return domain.Dataset{}, errors.New("file size exceeds maximum limit")
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	var dataset domain.Dataset
	decoder.DisallowUnknownFields()

	errar = decoder.Decode(&dataset)
	if errar != nil {
		return domain.Dataset{}, fmt.Errorf("failed to decode JSON: %v\n", errar)
	}
	if decoder.More() {
		return domain.Dataset{}, errors.New("extra data found after JSON files (possibly multiple JSON files)")
	}
	
	if dataset.Version != "1.0" {
		return domain.Dataset{}, fmt.Errorf("unsupported dataset version: %s\n", dataset.Version)
	}
	if dataset.Country == "" {
		return domain.Dataset{}, errors.New("missing country in dataset")
	}
	if dataset.Count < 5 || dataset.Count > 100 {
		return domain.Dataset{}, fmt.Errorf("invalid count in dataset: %d\n", dataset.Count)
	}
	if len(dataset.Values) != dataset.Count {
		return domain.Dataset{}, errors.New("values in dataset not detected")
	}
	for _, v := range dataset.Values{
		if v != v {
			return domain.Dataset{}, errors.New("value is not a number")
		} 
	}
	_ , errar = time.Parse(time.RFC3339,dataset.Time)
	if errar != nil {
		return domain.Dataset{}, fmt.Errorf("invalid time RFC3339: %v\n", errar)
	}

	return dataset, nil
}

func SaveReport(path string, html string) error{
	directory := filepath.Dir(path)
	if directory != "." && directory != "" {
		_, errar := os.Stat(directory)
		if os.IsNotExist(errar) {
			return fmt.Errorf("directory does not exist: %s\n", directory)
}
	

file, errar := os.OpenFile(path, os.O_CREATE| os.O_WRONLY|os.O_EXCL, 0604)
	if errar != nil {
		if os.IsExist(errar) {
			return fmt.Errorf("file already exists: %s\n", path)
		}
		return fmt.Errorf("failed to create file: %v\n", errar)
	}
	defer file.Close()

	_, errarwrite := file.Write([]byte(html))
	if errarwrite != nil {
		_ = os.Remove(path)
		return fmt.Errorf("failed to write to file: %v\n", errarwrite)
	}
}
	return nil
}