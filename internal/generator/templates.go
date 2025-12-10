package generator

// gitignoreTemplate is the template for .gitignore
const gitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# Go build outputs
bin/
build/
dist/

# Coverage
coverage.out
*.coverprofile

# Go modules vendored
vendor/

# Logs
*.log

# OS / editor files
.DS_Store
Thumbs.db

# IDEs
.vscode/
.idea/
*.iml

# Temporary files
*.swp
*~
`

// configTemplate is the template for internal/config/config.go
const configTemplate = `package config

import "os"

type Config struct {
	Env      string
	HTTPPort string

	PostgresURL string
	MySQLDSN    string
	MongoURI    string
	OracleDSN   string

	RedisAddr    string
	KafkaBrokers string
}

func Load() Config {
	return Config{
		Env:          getEnv("APP_ENV", "dev"),
		HTTPPort:     getEnv("APP_PORT", "8080"),
		PostgresURL:  getEnv("DB_POSTGRES_URL", ""),
		MySQLDSN:     getEnv("DB_MYSQL_DSN", ""),
		MongoURI:     getEnv("DB_MONGO_URI", ""),
		OracleDSN:    getEnv("DB_ORACLE_DSN", ""),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
`

// loggerTemplate is the template for internal/logger/logger.go
const loggerTemplate = `package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	z *zap.SugaredLogger
}

func New(env string) *Logger {
	cfg := zap.NewProductionConfig()
	if env == "dev" {
		cfg = zap.NewDevelopmentConfig()
	}
	l, err := cfg.Build()
	if err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}
	return &Logger{z: l.Sugar()}
}

func (l *Logger) Sync() {
	_ = l.z.Sync()
}

func (l *Logger) Info(msg string, kv ...interface{}) {
	l.z.Infow(msg, kv...)
}

func (l *Logger) Error(msg string, kv ...interface{}) {
	l.z.Errorw(msg, kv...)
}

func (l *Logger) Warn(msg string, kv ...interface{}) {
	l.z.Warnw(msg, kv...)
}

func (l *Logger) Debug(msg string, kv ...interface{}) {
	l.z.Debugw(msg, kv...)
}
`

// mainNetHTTPTemplate is the template for cmd/api/main.go using net/http
const mainNetHTTPTemplate = `package main

import (
	"fmt"
	"net/http"

	"{{.ModulePath}}/config"
	"{{.ModulePath}}/infrastructure/adapters/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env)
	defer log.Sync()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	addr := ":" + cfg.HTTPPort
	log.Info("starting server", "addr", addr, "env", cfg.Env)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("server error", "error", err)
	}
}
`

// mainChiTemplate is the template for cmd/api/main.go using chi
const mainChiTemplate = `package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"{{.ModulePath}}/config"
	"{{.ModulePath}}/infrastructure/adapters/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env)
	defer log.Sync()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	addr := ":" + cfg.HTTPPort
	log.Info("starting server", "addr", addr, "env", cfg.Env)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Error("server error", "error", err)
	}
}
`

// mainGinTemplate is the template for cmd/api/main.go using gin
const mainGinTemplate = `package main

import (
	"{{.ModulePath}}/config"
	"{{.ModulePath}}/infrastructure/adapters/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env)
	defer log.Sync()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	addr := ":" + cfg.HTTPPort
	log.Info("starting server", "addr", addr, "env", cfg.Env)
	if err := r.Run(addr); err != nil {
		log.Error("server error", "error", err)
	}
}
`

// mainFiberTemplate is the template for cmd/api/main.go using fiber
const mainFiberTemplate = `package main

import (
	"{{.ModulePath}}/config"
	"{{.ModulePath}}/infrastructure/adapters/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env)
	defer log.Sync()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	addr := ":" + cfg.HTTPPort
	log.Info("starting server", "addr", addr, "env", cfg.Env)
	if err := app.Listen(addr); err != nil {
		log.Error("server error", "error", err)
	}
}
`

// usecaseTemplate is the template for use cases
const usecaseTemplate = `package usecases

import (
	"context"
)

// {{.Name}}Input represents the input for {{.Name}}
type {{.Name}}Input struct {
	// Add input fields here
}

// {{.Name}}Output represents the output for {{.Name}}
type {{.Name}}Output struct {
	// Add output fields here
}

// {{.Name}}UseCase defines the use case interface
type {{.Name}}UseCase interface {
	Execute(ctx context.Context, input {{.Name}}Input) (*{{.Name}}Output, error)
}

// {{.LowerName}}UseCase is the implementation of {{.Name}}UseCase
type {{.LowerName}}UseCase struct {
	// Add dependencies here (repositories, etc.)
}

// New{{.Name}}UseCase creates a new instance of {{.Name}}UseCase
func New{{.Name}}UseCase() {{.Name}}UseCase {
	return &{{.LowerName}}UseCase{}
}

// Execute executes the {{.Name}} use case
func (uc *{{.LowerName}}UseCase) Execute(ctx context.Context, input {{.Name}}Input) (*{{.Name}}Output, error) {
	// TODO: Implement use case logic
	return &{{.Name}}Output{}, nil
}
`

// adapterTemplate is the template for adapters/repositories
const adapterTemplate = `package database

import (
	"context"
)

// {{.Name}} defines the repository interface
type {{.Name}} interface {
	// Add repository methods here
	FindByID(ctx context.Context, id string) (interface{}, error)
	Save(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id string) error
}

// {{.LowerName}} is the implementation of {{.Name}}
type {{.LowerName}} struct {
	// Add dependencies here (db connection, etc.)
}

// New{{.Name}} creates a new instance of {{.Name}}
func New{{.Name}}() {{.Name}} {
	return &{{.LowerName}}{}
}

// FindByID finds an entity by ID
func (r *{{.LowerName}}) FindByID(ctx context.Context, id string) (interface{}, error) {
	// TODO: Implement FindByID
	return nil, nil
}

// Save saves a new entity
func (r *{{.LowerName}}) Save(ctx context.Context, entity interface{}) error {
	// TODO: Implement Save
	return nil
}

// Update updates an existing entity
func (r *{{.LowerName}}) Update(ctx context.Context, entity interface{}) error {
	// TODO: Implement Update
	return nil
}

// Delete deletes an entity by ID
func (r *{{.LowerName}}) Delete(ctx context.Context, id string) error {
	// TODO: Implement Delete
	return nil
}
`

// modelTemplate is the template for domain models
const modelTemplate = `package models

import "time"

// {{.Name}} represents a {{.Name}} entity
type {{.Name}} struct {
	ID        string    ` + "`json:\"id\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	// Add more fields here
}

// New{{.Name}} creates a new {{.Name}} instance
func New{{.Name}}() *{{.Name}} {
	now := time.Now()
	return &{{.Name}}{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate validates the {{.Name}} entity
func (m *{{.Name}}) Validate() error {
	// TODO: Add validation logic
	return nil
}
`

// handlerTemplate is the template for HTTP handlers
const handlerTemplate = `package http

import (
	"encoding/json"
	"net/http"
)

// {{.Name}}Handler handles HTTP requests for {{.Name}}
type {{.Name}}Handler struct {
	// Add dependencies here (use cases, logger, etc.)
}

// New{{.Name}}Handler creates a new {{.Name}}Handler
func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{}
}

// Get handles GET requests
func (h *{{.Name}}Handler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement GET logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Post handles POST requests
func (h *{{.Name}}Handler) Post(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement POST logic
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

// Put handles PUT requests
func (h *{{.Name}}Handler) Put(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement PUT logic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// Delete handles DELETE requests
func (h *{{.Name}}Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement DELETE logic
	w.WriteHeader(http.StatusNoContent)
}
`
