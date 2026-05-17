package parse

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
)

// DockerEntry represents a single data point recording Docker container count and startup time.
type DockerEntry struct {
	DockerCount float64 `json:"dockerCount" csv:"n_containers"`
	StartupTime float64 `json:"startupTime" csv:"startup_ms"`
}

// LoadDataFromFolder loads Docker container and startup time data from all CSV files in a directory.
// It returns a slice of DockerEntry records parsed from the CSV files or an error if no CSV files are found or parsing fails.
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

// LoadDataFromFormFile parses Docker container and startup time data from a CSV file provided via HTTP multipart form upload.
// It returns a slice of DockerEntry records or an error if the file format is invalid.
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
