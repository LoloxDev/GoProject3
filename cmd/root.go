package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganalyzer",
	Short: "Analyse rapidement des fichiers de logs",
	Long:  "loganalyzer est un outil pédagogique qui illustre la gestion concurrente et la lecture d'un fichier de configuration JSON.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
