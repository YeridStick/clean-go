package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

// GenerateUsecase generates a new use case
func GenerateUsecase(name string) error {
	// Ensure we're in a Go project
	if !FileExists("go.mod") {
		return fmt.Errorf("no se encontró go.mod. Asegúrate de estar en la raíz del proyecto")
	}

	// Ensure usecase directory exists
	usecaseDir := "domain/usecases"
	if err := EnsureDir(usecaseDir); err != nil {
		return fmt.Errorf("error creando directorio usecases: %w", err)
	}

	// Prepare template data
	data := map[string]string{
		"Name":      ToPascalCase(name),
		"LowerName": ToCamelCase(name),
	}

	// Parse template
	tmpl, err := template.New("usecase").Parse(usecaseTemplate)
	if err != nil {
		return err
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// Write file
	filename := filepath.Join(usecaseDir, ToSnakeCase(name)+".go")
	if FileExists(filename) {
		return fmt.Errorf("el archivo %s ya existe", filename)
	}

	return WriteFile(filename, buf.Bytes())
}

// GenerateAdapter generates a new adapter/repository
func GenerateAdapter(name string) error {
	if !FileExists("go.mod") {
		return fmt.Errorf("no se encontró go.mod. Asegúrate de estar en la raíz del proyecto")
	}

	repoDir := "infrastructure/adapters/database"
	if err := EnsureDir(repoDir); err != nil {
		return fmt.Errorf("error creando directorio database: %w", err)
	}

	data := map[string]string{
		"Name":      ToPascalCase(name),
		"LowerName": ToCamelCase(name),
	}

	tmpl, err := template.New("adapter").Parse(adapterTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	filename := filepath.Join(repoDir, ToSnakeCase(name)+".go")
	if FileExists(filename) {
		return fmt.Errorf("el archivo %s ya existe", filename)
	}

	return WriteFile(filename, buf.Bytes())
}

// GenerateModel generates a new domain model
func GenerateModel(name string) error {
	if !FileExists("go.mod") {
		return fmt.Errorf("no se encontró go.mod. Asegúrate de estar en la raíz del proyecto")
	}

	domainDir := "domain/models"
	if err := EnsureDir(domainDir); err != nil {
		return fmt.Errorf("error creando directorio models: %w", err)
	}

	data := map[string]string{
		"Name": ToPascalCase(name),
	}

	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	filename := filepath.Join(domainDir, ToSnakeCase(name)+".go")
	if FileExists(filename) {
		return fmt.Errorf("el archivo %s ya existe", filename)
	}

	return WriteFile(filename, buf.Bytes())
}

// GenerateHandler generates a new HTTP handler
func GenerateHandler(name string) error {
	if !FileExists("go.mod") {
		return fmt.Errorf("no se encontró go.mod. Asegúrate de estar en la raíz del proyecto")
	}

	httpDir := "infrastructure/entrypoints/http"
	if err := EnsureDir(httpDir); err != nil {
		return fmt.Errorf("error creando directorio http: %w", err)
	}

	data := map[string]string{
		"Name": ToPascalCase(name),
	}

	tmpl, err := template.New("handler").Parse(handlerTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	filename := filepath.Join(httpDir, ToSnakeCase(name)+"_handler.go")
	if FileExists(filename) {
		return fmt.Errorf("el archivo %s ya existe", filename)
	}

	return WriteFile(filename, buf.Bytes())
}
