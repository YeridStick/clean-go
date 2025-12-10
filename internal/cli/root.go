package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cleango",
	Short: "Scaffold CLI para proyectos Go con Clean Architecture",
	Long: `cleango es un CLI que facilita la creación y gestión de proyectos Go
siguiendo los principios de Clean Architecture.

Características:
  • Generación rápida de proyectos con estructura predefinida
  • Múltiples frameworks HTTP (net/http, chi, gin, fiber)
  • Soporte para múltiples bases de datos (Postgres, MySQL, MongoDB, Oracle)
  • Generación de componentes (usecases, adapters, models, handlers)
  • Configuración centralizada y logger estructurado`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(addCmd)
}
