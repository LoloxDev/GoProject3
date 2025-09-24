package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"loganalyzer/internal/analyzer"
)

func Export(path string, results []analyzer.Result) error {
	if path == "" {
		return fmt.Errorf("chemin de sortie vide")
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser le rapport: %w", err)
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("impossible de créer le dossier cible: %w", err)
		}
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("écriture du rapport: %w", err)
	}

	return nil
}
