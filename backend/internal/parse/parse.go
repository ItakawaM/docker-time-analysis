package parse

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
)

type DockerEntry struct {
	DockerCount float64 `json:"dockerCount" csv:"n_containers"`
	StartupTime float64 `json:"startupTime" csv:"startup_ms"`
}

func LoadDataFromFolder(path string) ([]*DockerEntry, error) {
	matches, err := filepath.Glob(filepath.Join(path, "*.csv"))
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no csv files in the provided directory: %s", path)
	}

	data := []*DockerEntry{}
	for _, match := range matches {
		file, err := os.Open(match)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if err := validateHeaders(file); err != nil {
			return nil, err
		}

		fileData := []*DockerEntry{}
		if err := gocsv.UnmarshalFile(file, &fileData); err != nil {
			return nil, err
		}

		data = append(data, fileData...)
	}

	return data, nil
}

func LoadDataFromFormFile(file multipart.File) ([]*DockerEntry, error) {
	if err := validateHeaders(file); err != nil {
		return nil, err
	}

	data := []*DockerEntry{}
	if err := gocsv.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func validateHeaders(file multipart.File) error {
	r := csv.NewReader(file)
	headers, err := r.Read()
	if err != nil {
		return err
	}

	required := map[string]bool{"n_containers": false, "startup_ms": false}
	for _, h := range headers {
		if _, ok := required[h]; ok {
			required[h] = true
		}
	}

	for column, found := range required {
		if !found {
			return fmt.Errorf("missing required column: %s", column)
		}
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	return nil
}
