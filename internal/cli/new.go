package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YeridStick/cleango/internal/generator"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	modulePath string
	framework  string
	database   string
	useRedis   bool
	useKafka   bool
	nonInteractive bool
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Crea un nuevo proyecto Go con Clean Architecture",
	Long: `Crea un nuevo proyecto Go con estructura de Clean Architecture.

Ejemplos:
  cleango new my-service
  cleango new my-service --module github.com/user/my-service
  cleango new my-service --framework chi --database postgres
  cleango new my-service -m github.com/user/my-service -f gin -d postgres --redis --kafka`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNew,
}

func init() {
	newCmd.Flags().StringVarP(&modulePath, "module", "m", "", "Ruta del módulo Go (ej: github.com/user/project)")
	newCmd.Flags().StringVarP(&framework, "framework", "f", "", "Framework HTTP: nethttp, chi, gin, fiber")
	newCmd.Flags().StringVarP(&database, "database", "d", "", "Base de datos: none, postgres, mysql, mongodb, oracle")
	newCmd.Flags().BoolVar(&useRedis, "redis", false, "Incluir Redis")
	newCmd.Flags().BoolVar(&useKafka, "kafka", false, "Incluir Kafka")
	newCmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Modo no interactivo (usa valores por defecto)")
}

func runNew(cmd *cobra.Command, args []string) error {
	var projectName string

	if len(args) > 0 {
		projectName = args[0]
	} else if !nonInteractive {
		prompt := promptui.Prompt{
			Label: "Nombre del proyecto",
		}
		result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("operación cancelada")
		}
		projectName = result
	} else {
		return fmt.Errorf("nombre del proyecto requerido en modo no interactivo")
	}

	// Obtener module path
	if modulePath == "" && !nonInteractive {
		prompt := promptui.Prompt{
			Label:   "Módulo Go",
			Default: fmt.Sprintf("github.com/user/%s", projectName),
		}
		result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("operación cancelada")
		}
		modulePath = result
	} else if modulePath == "" {
		modulePath = fmt.Sprintf("github.com/user/%s", projectName)
	}

	// Obtener framework si no se especificó
	if framework == "" && !nonInteractive {
		prompt := promptui.Select{
			Label: "Selecciona framework HTTP",
			Items: []string{"nethttp", "chi", "gin", "fiber"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("operación cancelada")
		}
		framework = result
	} else if framework == "" {
		framework = "nethttp"
	}

	// Obtener base de datos si no se especificó
	if database == "" && !nonInteractive {
		prompt := promptui.Select{
			Label: "Selecciona base de datos",
			Items: []string{"none", "postgres", "mysql", "mongodb", "oracle"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("operación cancelada")
		}
		database = result
	} else if database == "" {
		database = "none"
	}

	// Preguntar por Redis y Kafka en modo interactivo
	if !nonInteractive && !cmd.Flags().Changed("redis") {
		prompt := promptui.Prompt{
			Label:     "¿Agregar Redis?",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		useRedis = (err == nil)
	}

	if !nonInteractive && !cmd.Flags().Changed("kafka") {
		prompt := promptui.Prompt{
			Label:     "¿Agregar Kafka?",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		useKafka = (err == nil)
	}

	// Crear configuración del proyecto
	config := generator.ProjectConfig{
		Name:       projectName,
		ModulePath: modulePath,
		Framework:  framework,
		Database:   database,
		UseRedis:   useRedis,
		UseKafka:   useKafka,
	}

	// Mostrar resumen
	fmt.Println("\n=== Resumen del proyecto ===")
	fmt.Printf("Nombre:     %s\n", config.Name)
	fmt.Printf("Módulo:     %s\n", config.ModulePath)
	fmt.Printf("Framework:  %s\n", config.Framework)
	fmt.Printf("Database:   %s\n", config.Database)
	fmt.Printf("Redis:      %v\n", config.UseRedis)
	fmt.Printf("Kafka:      %v\n", config.UseKafka)
	fmt.Println()

	// Confirmar en modo interactivo
	if !nonInteractive {
		prompt := promptui.Prompt{
			Label:     "¿Crear proyecto?",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("operación cancelada")
		}
	}

	// Obtener directorio actual
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error obteniendo directorio actual: %w", err)
	}

	targetDir := filepath.Join(cwd, projectName)

	// Verificar si la carpeta existe
	if _, err := os.Stat(targetDir); err == nil {
		if !nonInteractive {
			prompt := promptui.Select{
				Label: fmt.Sprintf("La carpeta '%s' ya existe. ¿Qué deseas hacer?", projectName),
				Items: []string{"Usar carpeta existente", "Cancelar"},
			}
			idx, _, err := prompt.Run()
			if err != nil || idx != 0 {
				return fmt.Errorf("operación cancelada")
			}
		} else {
			fmt.Printf("⚠️  La carpeta '%s' ya existe. Usando carpeta existente.\n", projectName)
		}
	} else {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("error creando directorio: %w", err)
		}
	}

	// Generar proyecto
	fmt.Println("\n🚀 Generando proyecto...")
	if err := generator.GenerateProject(targetDir, config); err != nil {
		return fmt.Errorf("error generando proyecto: %w", err)
	}

	// Mensaje final
	fmt.Println("\n✅ Proyecto creado exitosamente!")
	fmt.Println("\nPróximos pasos:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  go run ./cmd/api")
	if envs := databaseEnvVars(config.Database); len(envs) > 0 {
		fmt.Printf("\nConfiguración base de %s generada en infrastructure/adapters/database/%s.go\n", config.Database, config.Database)
		fmt.Printf("Completa tus variables de conexión en .env (ejemplos en .env.example): %s\n", strings.Join(envs, ", "))
	}
	fmt.Println("\nPara agregar componentes:")
	fmt.Println("  cleango add usecase <nombre>")
	fmt.Println("  cleango add adapter <nombre>")
	fmt.Println("  cleango add model <nombre>")
	fmt.Println("  cleango add handler <nombre>")

	return nil
}

func databaseEnvVars(db string) []string {
	switch db {
	case "postgres":
		return []string{"DB_POSTGRES_URL"}
	case "mysql":
		return []string{"DB_MYSQL_DSN"}
	case "mongodb":
		return []string{"DB_MONGO_URI", "DB_MONGO_DATABASE"}
	case "oracle":
		return []string{"DB_ORACLE_DSN"}
	default:
		return nil
	}
}
