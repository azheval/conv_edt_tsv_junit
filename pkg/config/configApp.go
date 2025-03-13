package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	InputFileFolder  string   `json:"input_file_folder"`
	OutputFileFolder string   `json:"output_file_folder"`
	SkipCategories   []string `json:"skip_categories"`
	SkipObjects      []string `json:"skip_objects"`
	SkipSignificanceCcategories     []string `json:"skip_significance_categories"`
	SkipErrorText     []string `json:"skip_error_text"`
	SkipErrorsFile   string   `json:"skip_errors_file"`
}

func (c *AppConfig) Load(filePath string) {
	workspace, _ := os.Getwd()
	configData, err := os.ReadFile(filepath.Join(workspace, filePath))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configData, c)
	if err != nil {
		panic(err)
	}
}
