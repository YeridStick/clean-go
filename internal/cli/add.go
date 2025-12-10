package cli

import (
	"fmt"

	"github.com/YeridStick/cleango/internal/generator"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Agrega componentes al proyecto",
	Long: `Agrega nuevos componentes a tu proyecto siguiendo Clean Architecture.

Componentes disponibles:
  • usecase  - Crea un nuevo caso de uso en domain/usecases
  • adapter  - Crea un nuevo adaptador en infrastructure/adapters/database
  • model    - Crea un nuevo modelo en domain/models
  • handler  - Crea un nuevo handler HTTP en infrastructure/entrypoints/http`,
}

var addUsecaseCmd = &cobra.Command{
	Use:   "usecase [nombre]",
	Short: "Crea un nuevo caso de uso",
	Long: `Crea un nuevo caso de uso en domain/usecases/.

El caso de uso seguirá el patrón de Clean Architecture con:
  • Interface del caso de uso (puerto)
  • Implementación concreta
  • Struct de entrada (Input)
  • Struct de salida (Output)

Ejemplo:
  cleango add usecase GetUser
  cleango add usecase CreateOrder`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("🔧 Generando caso de uso '%s'...\n", name)

		if err := generator.GenerateUsecase(name); err != nil {
			return fmt.Errorf("error generando caso de uso: %w", err)
		}

		fmt.Printf("✅ Caso de uso '%s' creado exitosamente!\n", name)
		fmt.Printf("   Archivo: domain/usecases/%s.go\n", generator.ToSnakeCase(name))
		return nil
	},
}

var addAdapterCmd = &cobra.Command{
	Use:   "adapter [nombre]",
	Short: "Crea un nuevo adaptador/repositorio",
	Long: `Crea un nuevo adaptador en infrastructure/adapters/database/.

El adaptador seguirá el patrón Repository con:
  • Interface del repositorio
  • Implementación concreta
  • Métodos CRUD básicos

Ejemplo:
  cleango add adapter UserRepository
  cleango add adapter ProductRepository`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("🔧 Generando adaptador '%s'...\n", name)

		if err := generator.GenerateAdapter(name); err != nil {
			return fmt.Errorf("error generando adaptador: %w", err)
		}

		fmt.Printf("✅ Adaptador '%s' creado exitosamente!\n", name)
		fmt.Printf("   Archivo: infrastructure/adapters/database/%s.go\n", generator.ToSnakeCase(name))
		return nil
	},
}

var addModelCmd = &cobra.Command{
	Use:   "model [nombre]",
	Short: "Crea un nuevo modelo de dominio",
	Long: `Crea un nuevo modelo en domain/models/.

El modelo será una estructura básica que representa una entidad de dominio.
Entidades de dominio son objetos puros de negocio sin dependencias externas.

Ejemplo:
  cleango add model User
  cleango add model Product`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("🔧 Generando modelo '%s'...\n", name)

		if err := generator.GenerateModel(name); err != nil {
			return fmt.Errorf("error generando modelo: %w", err)
		}

		fmt.Printf("✅ Modelo '%s' creado exitosamente!\n", name)
		fmt.Printf("   Archivo: domain/models/%s.go\n", generator.ToSnakeCase(name))
		return nil
	},
}

var addHandlerCmd = &cobra.Command{
	Use:   "handler [nombre]",
	Short: "Crea un nuevo handler HTTP",
	Long: `Crea un nuevo handler HTTP en infrastructure/entrypoints/http/.

El handler incluirá:
  • Estructura del handler
  • Métodos HTTP básicos (Get, Post, Put, Delete)
  • Validación de entrada
  • Manejo de errores

Ejemplo:
  cleango add handler User
  cleango add handler Product`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		fmt.Printf("🔧 Generando handler '%s'...\n", name)

		if err := generator.GenerateHandler(name); err != nil {
			return fmt.Errorf("error generando handler: %w", err)
		}

		fmt.Printf("✅ Handler '%s' creado exitosamente!\n", name)
		fmt.Printf("   Archivo: infrastructure/entrypoints/http/%s_handler.go\n", generator.ToSnakeCase(name))
		return nil
	},
}

func init() {
	addCmd.AddCommand(addUsecaseCmd)
	addCmd.AddCommand(addAdapterCmd)
	addCmd.AddCommand(addModelCmd)
	addCmd.AddCommand(addHandlerCmd)
}
