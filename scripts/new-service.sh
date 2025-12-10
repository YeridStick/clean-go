#!/usr/bin/env bash
set -e

# ==========================
# 0. Datos básicos (servicio y módulo)
# ==========================

SERVICE_NAME=$1
MODULE_PATH=$2 # ej: github.com/yerid/turismo-back-go

if [ -z "$SERVICE_NAME" ]; then
  read -rp "Ingresa el nombre del servicio (carpeta raíz del proyecto): " SERVICE_NAME
fi

if [ -z "$MODULE_PATH" ]; then
  read -rp "Ingresa el nombre del módulo Go (por ejemplo github.com/yerid/turismo-back-go): " MODULE_PATH
fi

if [ -z "$SERVICE_NAME" ] || [ -z "$MODULE_PATH" ]; then
  echo "❌ Nombre de servicio y módulo Go son obligatorios."
  exit 1
fi

echo "=== Scaffold Go ==="
echo "Servicio:   $SERVICE_NAME"
echo "Módulo Go:  $MODULE_PATH"
echo "Directorio actual: $(pwd)"
echo "El proyecto se creará en: $(pwd)/$SERVICE_NAME"
echo

read -rp "¿Crear/usar esta carpeta para el proyecto? [y/N]: " DIR_CONFIRM
if [[ "$DIR_CONFIRM" != "y" && "$DIR_CONFIRM" != "Y" ]]; then
  echo "❌ Creación de proyecto cancelada por el usuario."
  exit 0
fi

# Si la carpeta existe, preguntar qué hacer
if [ -d "$SERVICE_NAME" ]; then
  echo
  echo "⚠️ La carpeta '$SERVICE_NAME' ya existe."
  echo "  1) Usar la carpeta tal cual (puede mezclar archivos)."
  echo "  2) Borrar la carpeta y crearla de nuevo (rm -rf)."
  echo "  3) Cancelar."
  read -rp "Elige una opción [1-3]: " EXIST_CHOICE

  case "$EXIST_CHOICE" in
    1)
      echo "➡ Usando carpeta existente '$SERVICE_NAME'."
      ;;
    2)
      echo "❗ Se eliminará la carpeta '$SERVICE_NAME' y su contenido."
      read -rp "¿Seguro? Esta acción no se puede deshacer. [escribe 'BORRAR' para confirmar]: " CONFIRM_DELETE
      if [ "$CONFIRM_DELETE" != "BORRAR" ]; then
        echo "❌ No se confirmó el borrado. Cancelando."
        exit 0
      fi
      rm -rf "$SERVICE_NAME"
      mkdir -p "$SERVICE_NAME"
      ;;
    3|*)
      echo "❌ Operación cancelada."
      exit 0
      ;;
  esac
else
  # Si no existe, crearla normalmente
  mkdir -p "$SERVICE_NAME"
fi

cd "$SERVICE_NAME"

# ==========================
# 1. Elegir framework HTTP
# ==========================

echo
echo "Selecciona framework HTTP:"
echo "  1) net/http (sin framework)"
echo "  2) chi"
echo "  3) gin"
echo "  4) fiber"
read -rp "Opción [1-4]: " FW_CHOICE

FRAMEWORK="nethttp"
case "$FW_CHOICE" in
  1|"") FRAMEWORK="nethttp" ;;
  2) FRAMEWORK="chi" ;;
  3) FRAMEWORK="gin" ;;
  4) FRAMEWORK="fiber" ;;
  *) echo "Opción inválida"; exit 1 ;;
esac

# ==========================
# 2. Elegir base de datos principal
# ==========================

echo
echo "Selecciona base de datos principal:"
echo "  1) Ninguna"
echo "  2) Postgres"
echo "  3) MySQL"
echo "  4) MongoDB"
echo "  5) Oracle (avanzado)"
read -rp "Opción [1-5]: " DB_CHOICE

DB_MAIN="none"
case "$DB_CHOICE" in
  1|"") DB_MAIN="none" ;;
  2) DB_MAIN="postgres" ;;
  3) DB_MAIN="mysql" ;;
  4) DB_MAIN="mongodb" ;;
  5) DB_MAIN="oracle" ;;
  *) echo "Opción inválida"; exit 1 ;;
esac

# ==========================
# 3. Extras: Redis / Kafka
# ==========================

echo
read -rp "¿Agregar Redis? [y/N]: " REDIS_CHOICE
read -rp "¿Agregar Kafka? [y/N]: " KAFKA_CHOICE

USE_REDIS="no"
USE_KAFKA="no"
[[ "$REDIS_CHOICE" == "y" || "$REDIS_CHOICE" == "Y" ]] && USE_REDIS="yes"
[[ "$KAFKA_CHOICE" == "y" || "$KAFKA_CHOICE" == "Y" ]] && USE_KAFKA="yes"

# ==========================
# 4. Resumen de dependencias
# ==========================

echo
echo "=== Resumen de configuración ==="
echo "Framework HTTP: $FRAMEWORK"
echo "Base de datos:  $DB_MAIN"
echo "Redis:          $USE_REDIS"
echo "Kafka:          $USE_KAFKA"
echo
echo "Se instalarán módulos Go como mínimo:"

echo "- go.uber.org/zap (logger)"

case "$FRAMEWORK" in
  "chi")   echo "- github.com/go-chi/chi/v5" ;;
  "gin")   echo "- github.com/gin-gonic/gin" ;;
  "fiber") echo "- github.com/gofiber/fiber/v2" ;;
esac

case "$DB_MAIN" in
  "postgres") echo "- github.com/jackc/pgx/v5/stdlib" ;;
  "mysql")    echo "- github.com/go-sql-driver/mysql" ;;
  "mongodb")  echo "- go.mongodb.org/mongo-driver/mongo" ;;
  "oracle")   echo "- github.com/godror/godror (requiere cliente Oracle)" ;;
esac

if [ "$USE_REDIS" == "yes" ]; then
  echo "- github.com/redis/go-redis/v9"
fi

if [ "$USE_KAFKA" == "yes" ]; then
  echo "- github.com/segmentio/kafka-go"
fi

echo
read -rp "¿Continuar con la creación del proyecto? [y/N]: " CONFIRM
if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
  echo "❌ Cancelado."
  exit 0
fi

# ==========================
# 5. Inicializar módulo Go
# ==========================

# Si ya hay go.mod, preguntar antes de sobrescribir
if [ -f "go.mod" ]; then
  echo "⚠️ Ya existe go.mod en esta carpeta."
  read -rp "¿Quieres reutilizarlo? [Y/n]: " REUSE_MOD
  if [[ "$REUSE_MOD" == "n" || "$REUSE_MOD" == "N" ]]; then
    echo "❌ No se reutilizó go.mod. Cancelando para evitar conflictos."
    exit 0
  fi
else
  go mod init "$MODULE_PATH"
fi

# ==========================
# 6. Instalar dependencias
# ==========================

go get go.uber.org/zap

case "$FRAMEWORK" in
  "chi")   go get github.com/go-chi/chi/v5 ;;
  "gin")   go get github.com/gin-gonic/gin ;;
  "fiber") go get github.com/gofiber/fiber/v2 ;;
esac

case "$DB_MAIN" in
  "postgres") go get github.com/jackc/pgx/v5/stdlib ;;
  "mysql")    go get github.com/go-sql-driver/mysql ;;
  "mongodb")  go get go.mongodb.org/mongo-driver/mongo ;;
  "oracle")   go get github.com/godror/godror ;;
esac

if [ "$USE_REDIS" == "yes" ]; then
  go get github.com/redis/go-redis/v9
fi

if [ "$USE_KAFKA" == "yes" ]; then
  go get github.com/segmentio/kafka-go
fi

# ==========================
# 7. Carpetas Clean Architecture
# ==========================

mkdir -p cmd/api
mkdir -p internal/config
mkdir -p internal/logger
mkdir -p internal/http
mkdir -p internal/domain
mkdir -p internal/usecase
mkdir -p internal/repository
mkdir -p internal/db
mkdir -p migrations

# ==========================
# 7.1 .gitignore básico para proyectos Go
# ==========================

cat > .gitignore << 'EOF'
# Binaries
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

# Shell scripts (por si copias generadores dentro del servicio)
*.sh
EOF

# ==========================
# 8. Archivo de config genérico
# ==========================

cat > internal/config/config.go << 'EOF'
package config

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
		Env:         getEnv("APP_ENV", "dev"),
		HTTPPort:    getEnv("APP_PORT", "8080"),
		PostgresURL: getEnv("DB_POSTGRES_URL", ""),
		MySQLDSN:    getEnv("DB_MYSQL_DSN", ""),
		MongoURI:    getEnv("DB_MONGO_URI", ""),
		OracleDSN:   getEnv("DB_ORACLE_DSN", ""),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		KafkaBrokers:getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
EOF

# ==========================
# 9. Logger con zap
# ==========================

cat > internal/logger/logger.go << 'EOF'
package logger

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
EOF

# ==========================
# 10. main.go según framework elegido
# ==========================

MAIN_FILE=cmd/api/main.go

if [ "$FRAMEWORK" = "nethttp" ]; then
cat > "$MAIN_FILE" << EOF
package main

import (
	"fmt"
	"net/http"

	"$MODULE_PATH/internal/config"
	"$MODULE_PATH/internal/logger"
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
EOF
fi

if [ "$FRAMEWORK" = "chi" ]; then
cat > "$MAIN_FILE" << EOF
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"$MODULE_PATH/internal/config"
	"$MODULE_PATH/internal/logger"
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
EOF
fi

if [ "$FRAMEWORK" = "gin" ]; then
cat > "$MAIN_FILE" << EOF
package main

import (
	"$MODULE_PATH/internal/config"
	"$MODULE_PATH/internal/logger"

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
EOF
fi

if [ "$FRAMEWORK" = "fiber" ]; then
cat > "$MAIN_FILE" << EOF
package main

import (
	"$MODULE_PATH/internal/config"
	"$MODULE_PATH/internal/logger"

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
EOF
fi

echo
echo "✔ Proyecto $SERVICE_NAME creado."
echo "   cd $SERVICE_NAME"
echo "   APP_PORT=8080 go run ./cmd/api"
