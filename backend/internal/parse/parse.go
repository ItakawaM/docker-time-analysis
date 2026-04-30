package parse

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
)

type DockerEntry struct {
	DockerCount float64 `json:"docker_count" csv:"n_containers"`
	StartupTime float64 `json:"startup_time" csv:"startup_ms"`
}

func LoadDataFromFolder(path string) ([]*DockerEntry, error) {
	matches, err := filepath.Glob(filepath.Join(path, "*.csv"))
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no csv files in the provided directory: %s", path)
	}

	dockerData := []*DockerEntry{}
	for _, match := range matches {
		file, err := os.Open(match)
		if err != nil {
			return nil, err
		}

		fileData := []*DockerEntry{}
		if err := gocsv.UnmarshalFile(file, &fileData); err != nil {
			return nil, err
		}

		dockerData = append(dockerData, fileData...)
	}

	return dockerData, nil
}
