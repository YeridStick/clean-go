package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// GenerateProject generates a new Go project with Clean Architecture
func GenerateProject(targetDir string, config ProjectConfig) error {
	// Create directory structure following Clean Architecture
	dirs := []string{
		"cmd/api",
		"config",
		"domain/models",
		"domain/usecases",
		"infrastructure/adapters/database",
		"infrastructure/adapters/logger",
		"infrastructure/entrypoints/http",
		"migrations",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(targetDir, dir)
		if err := EnsureDir(dirPath); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	// Change to target directory
	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("error changing to target directory: %w", err)
	}

	// Initialize go module if not exists
	if !FileExists("go.mod") {
		cmd := exec.Command("go", "mod", "init", config.ModulePath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("error initializing go module: %w\nOutput: %s", err, output)
		}
	}

	// Generate .gitignore
	if err := WriteFile(".gitignore", []byte(gitignoreTemplate)); err != nil {
		return fmt.Errorf("error creating .gitignore: %w", err)
	}

	// Generate config
	configPath := filepath.Join("config/config.go")
	if err := WriteFile(configPath, []byte(configTemplate)); err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	// Generate logger
	loggerPath := filepath.Join("infrastructure/adapters/logger/logger.go")
	if err := WriteFile(loggerPath, []byte(loggerTemplate)); err != nil {
		return fmt.Errorf("error creating logger: %w", err)
	}

	// Generate README with structure explanation
	readmePath := filepath.Join("README.md")
	readmeContent := generateReadme(config)
	if err := WriteFile(readmePath, []byte(readmeContent)); err != nil {
		return fmt.Errorf("error creating README.md: %w", err)
	}

	// Generate main.go based on framework
	mainContent, err := generateMainFile(config)
	if err != nil {
		return fmt.Errorf("error generating main.go: %w", err)
	}

	mainPath := filepath.Join("cmd/api/main.go")
	if err := WriteFile(mainPath, mainContent); err != nil {
		return fmt.Errorf("error creating main.go: %w", err)
	}

	// Generate database-specific files
	if err := generateDatabaseFiles(config); err != nil {
		return fmt.Errorf("error generating database files: %w", err)
	}

	// Generate .env.example
	envContent, err := renderEnvExample(config)
	if err != nil {
		return fmt.Errorf("error creating .env.example: %w", err)
	}
	if err := WriteFile(".env.example", envContent); err != nil {
		return fmt.Errorf("error creating .env.example: %w", err)
	}

	// Generate Makefile if PostgreSQL
	if config.Database == "postgres" {
		makefileContent, err := generateMakefile(config)
		if err != nil {
			return fmt.Errorf("error generating Makefile: %w", err)
		}
		if err := WriteFile("Makefile", makefileContent); err != nil {
			return fmt.Errorf("error creating Makefile: %w", err)
		}
	}

	// Install dependencies
	fmt.Println("📦 Instalando dependencias...")
	for _, dep := range config.GetDependencies() {
		fmt.Printf("   - %s\n", dep)
		cmd := exec.Command("go", "get", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("⚠️  Advertencia: No se pudo instalar %s: %v\n", dep, err)
		}
	}

	// Run go mod tidy
	fmt.Println("🧹 Ejecutando go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Advertencia: Error ejecutando go mod tidy: %v\n", err)
	}

	return nil
}

// generateMainFile generates the main.go file based on the framework
func generateMainFile(config ProjectConfig) ([]byte, error) {
	var tmplStr string

	switch config.Framework {
	case "chi":
		tmplStr = mainChiTemplate
	case "gin":
		tmplStr = mainGinTemplate
	case "fiber":
		tmplStr = mainFiberTemplate
	default:
		tmplStr = mainNetHTTPTemplate
	}

	tmpl, err := template.New("main").Parse(tmplStr)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generateDatabaseFiles generates database-specific files based on configuration
func generateDatabaseFiles(config ProjectConfig) error {
	templates := map[string]struct {
		filename string
		content  string
		test     string
	}{
		"postgres": {"postgres", postgresTemplate, postgresTestTemplate},
		"mysql":    {"mysql", mysqlTemplate, mysqlTestTemplate},
		"mongodb":  {"mongodb", mongoTemplate, mongoTestTemplate},
		"oracle":   {"oracle", oracleTemplate, oracleTestTemplate},
	}

	tmpl, ok := templates[config.Database]
	if !ok {
		return nil
	}

	dbPath := filepath.Join("infrastructure/adapters/database", tmpl.filename+".go")
	if err := WriteFile(dbPath, []byte(tmpl.content)); err != nil {
		return fmt.Errorf("error creating %s: %w", dbPath, err)
	}

	if tmpl.test != "" {
		testPath := filepath.Join("infrastructure/adapters/database", tmpl.filename+"_test.go")
		if err := WriteFile(testPath, []byte(tmpl.test)); err != nil {
			return fmt.Errorf("error creating %s: %w", testPath, err)
		}
	}

	return nil
}

// generateMakefile generates a Makefile based on configuration
func generateMakefile(config ProjectConfig) ([]byte, error) {
	tmpl, err := template.New("makefile").Parse(makefileTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generateReadme generates a README with the project structure
func generateReadme(config ProjectConfig) string {
	readme := "# " + config.Name + "\n\n"
	readme += "Proyecto Go con Clean Architecture\n\n"
	readme += "## Estructura del Proyecto\n\n"
	readme += "Este proyecto sigue los principios de Clean Architecture:\n\n"
	readme += "```\n"
	readme += config.Name + "/\n"
	readme += "├── cmd/api/                          # Punto de entrada de la aplicación\n"
	readme += "│   └── main.go\n"
	readme += "├── config/                           # Configuraciones\n"
	readme += "│   └── config.go\n"
	readme += "├── domain/                           # Capa de Dominio (Reglas de Negocio)\n"
	readme += "│   ├── models/                       # Entidades de dominio\n"
	readme += "│   └── usecases/                     # Casos de uso (puertos)\n"
	readme += "├── infrastructure/                   # Capa de Infraestructura\n"
	readme += "│   ├── adapters/                     # Adaptadores (implementaciones)\n"
	readme += "│   │   ├── database/                 # Repositorios de base de datos\n"
	readme += "│   │   └── logger/                   # Sistema de logging\n"
	readme += "│   └── entrypoints/                  # Puntos de entrada\n"
	readme += "│       └── http/                     # Handlers HTTP\n"
	readme += "├── migrations/                       # Migraciones de base de datos\n"
	readme += "├── .env.example                      # Variables de entorno ejemplo\n"
	readme += "├── .gitignore\n"
	readme += "├── go.mod\n"
	if config.Database == "postgres" {
		readme += "├── Makefile                          # Comandos útiles\n"
	}
	readme += "└── README.md\n"
	readme += "```\n\n"
	readme += "## Capas de Clean Architecture\n\n"
	readme += "### Domain (Dominio)\n"
	readme += "- **models/**: Entidades de dominio, objetos de negocio puros sin dependencias externas\n"
	readme += "- **usecases/**: Lógica de negocio, casos de uso de la aplicación (interfaces/puertos)\n\n"
	readme += "### Infrastructure (Infraestructura)\n"
	readme += "- **adapters/**: Implementaciones concretas de los puertos definidos en domain\n"
	readme += "  - **database/**: Repositorios que acceden a la base de datos\n"
	readme += "  - **logger/**: Sistema de logging\n"
	readme += "- **entrypoints/**: Puntos de entrada a la aplicación\n"
	readme += "  - **http/**: Handlers HTTP, controladores REST\n\n"
	readme += "## Configuración\n\n"
	readme += fmt.Sprintf("- **Framework**: %s\n", config.Framework)
	readme += fmt.Sprintf("- **Base de datos**: %s\n", config.Database)
	readme += fmt.Sprintf("- **Redis**: %v\n", config.UseRedis)
	readme += fmt.Sprintf("- **Kafka**: %v\n\n", config.UseKafka)

	// Add database-specific quick start
	if config.Database == "postgres" {
		readme += "## Inicio Rápido con PostgreSQL\n\n"
		readme += "### 1. Configurar variables de entorno\n\n"
		readme += "Copia el archivo `.env.example` a `.env` y configura tu URL de PostgreSQL:\n\n"
		readme += "```bash\n"
		readme += "cp .env.example .env\n"
		readme += "# Edita .env con tus credenciales\n"
		readme += "```\n\n"
		readme += "### 2. Iniciar base de datos (opcional con Docker)\n\n"
		readme += "```bash\n"
		readme += "make db-up\n"
		readme += "```\n\n"
		readme += "### 3. Ejecutar tests\n\n"
		readme += "```bash\n"
		readme += "# Solo tests unitarios (sin base de datos)\n"
		readme += "make test-short\n\n"
		readme += "# Tests de integración (requiere base de datos)\n"
		readme += "make test-integration\n\n"
		readme += "# Todos los tests\n"
		readme += "make test\n"
		readme += "```\n\n"
		readme += "### 4. Ejecutar la aplicación\n\n"
		readme += "```bash\n"
		readme += "make dev\n"
		readme += "# o\n"
		readme += "go run ./cmd/api\n"
		readme += "```\n\n"
	} else {
		readme += "## Ejecutar la aplicación\n\n"
		readme += "```bash\n"
		readme += "go run ./cmd/api\n"
		readme += "```\n\n"
	}

	readme += "## Agregar componentes\n\n"
	readme += "```bash\n"
	readme += "# Agregar un nuevo modelo de dominio\n"
	readme += "cleango add model User\n\n"
	readme += "# Agregar un nuevo caso de uso\n"
	readme += "cleango add usecase CreateUser\n\n"
	readme += "# Agregar un nuevo adaptador/repositorio\n"
	readme += "cleango add adapter UserRepository\n\n"
	readme += "# Agregar un nuevo handler HTTP\n"
	readme += "cleango add handler User\n"
	readme += "```\n\n"
	readme += "## Principios de Clean Architecture\n\n"
	readme += "1. **Independencia de frameworks**: El dominio no depende de frameworks\n"
	readme += "2. **Testeable**: La lógica de negocio puede testearse sin UI, BD, etc.\n"
	readme += "3. **Independencia de la UI**: La UI puede cambiar sin afectar el dominio\n"
	readme += "4. **Independencia de la BD**: Puedes cambiar de BD sin afectar el negocio\n"
	readme += "5. **Independencia de agentes externos**: El dominio no conoce nada del mundo exterior\n\n"
	readme += "## Flujo de dependencias\n\n"
	readme += "```\n"
	readme += "entrypoints -> usecases -> models\n"
	readme += "              ↑\n"
	readme += "           adapters\n"
	readme += "```\n\n"
	readme += "Las dependencias siempre apuntan hacia adentro (hacia el dominio).\n"
	return readme
}

func renderEnvExample(config ProjectConfig) ([]byte, error) {
	tmpl, err := template.New("env").Parse(envExampleTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
