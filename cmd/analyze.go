package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"loganalyzer/internal/analyzer"
	"loganalyzer/internal/config"
	"loganalyzer/internal/reporter"
)

var (
	configPath string
	outputPath string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse les logs décrits dans un fichier JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := config.Load(configPath)
		if err != nil {
			var parseErr *config.ParseError
			switch {
			case errors.Is(err, config.ErrConfigFileAccess):
				return fmt.Errorf("impossible de lire le fichier de configuration %q: %w", configPath, err)
			case errors.As(err, &parseErr):
				return fmt.Errorf("fichier de configuration invalide: %w", err)
			default:
				return err
			}
		}

		results := analyzer.Analyze(entries)

		for _, res := range results {
			fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s (%s) -> %s\n", res.Status, res.LogID, res.FilePath, res.Message)
			if res.Err != nil {
				var (
					logErr   *analyzer.LogFileError
					parseErr *analyzer.ParseError
				)

				switch {
				case errors.As(res.Err, &logErr):
					fmt.Fprintf(cmd.OutOrStdout(), "    problème d'accès au fichier: %s\n", logErr.Unwrap())
				case errors.As(res.Err, &parseErr):
					fmt.Fprintf(cmd.OutOrStdout(), "    erreur d'analyse: %s\n", parseErr.Unwrap())
				default:
					fmt.Fprintf(cmd.OutOrStdout(), "    erreur: %s\n", res.Err)
				}
			}
		}

		if outputPath != "" {
			if err := reporter.Export(outputPath, results); err != nil {
				return fmt.Errorf("export du rapport JSON: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Rapport enregistré dans %s\n", outputPath)
		}

		return nil
	},
}

func init() {
	analyzeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Chemin du fichier de configuration JSON")
	analyzeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Chemin du rapport JSON à générer")
	analyzeCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(analyzeCmd)
}
