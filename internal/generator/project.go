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
	// Create directory structure
	dirs := []string{
		"cmd/api",
		"internal/config",
		"internal/logger",
		"internal/http",
		"internal/domain",
		"internal/usecase",
		"internal/repository",
		"internal/db",
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
	configPath := filepath.Join("internal/config/config.go")
	if err := WriteFile(configPath, []byte(configTemplate)); err != nil {
		return fmt.Errorf("error creating config: %w", err)
	}

	// Generate logger
	loggerPath := filepath.Join("internal/logger/logger.go")
	if err := WriteFile(loggerPath, []byte(loggerTemplate)); err != nil {
		return fmt.Errorf("error creating logger: %w", err)
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
