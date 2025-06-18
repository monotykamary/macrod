package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/monotykamary/macrod/pkg/models"
)

type Storage struct {
	configPath string
}

func New() *Storage {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "macrod")
	os.MkdirAll(configDir, 0755)
	
	return &Storage{
		configPath: filepath.Join(configDir, "macros.json"),
	}
}

func (s *Storage) LoadMacros() ([]models.Macro, error) {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Macro{}, nil
		}
		return nil, err
	}

	var macros []models.Macro
	if err := json.Unmarshal(data, &macros); err != nil {
		return nil, err
	}

	return macros, nil
}

func (s *Storage) SaveMacros(macros []models.Macro) error {
	data, err := json.MarshalIndent(macros, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.configPath, data, 0644)
}