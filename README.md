# clean-go 🚀

**CLI profesional para crear proyectos Go con Clean Architecture**

`cleango` es una herramienta de línea de comandos que facilita la creación y gestión de proyectos Go siguiendo los principios de Clean Architecture. Similar a Spring CLI o Rails, pero diseñado específicamente para Go.

## 🎯 Características

- ✨ Generación rápida de proyectos con estructura predefinida
- 🎨 Múltiples frameworks HTTP: `net/http`, `chi`, `gin`, `fiber`
- 💾 Soporte para múltiples bases de datos: Postgres, MySQL, MongoDB, Oracle
- 📦 Instalación automática de dependencias
- 🔧 Generación de componentes: usecases, adapters, models, handlers
- ⚙️ Configuración centralizada y logger estructurado
- 🎯 Modo interactivo y no interactivo

---

## 📋 Requisitos

- **Go 1.20 o superior**
- Git (opcional, para clonar el repositorio)

Verifica tu versión de Go:

```bash
go version
```

---

## 🔧 Instalación

### Opción 1: Instalación directa con `go install`

```bash
go install github.com/YeridStick/cleango/cmd/cleango@latest
```

Asegúrate de que `$GOPATH/bin` esté en tu `PATH`:

```bash
# Linux/macOS
export PATH=$PATH:$(go env GOPATH)/bin

# Windows (PowerShell)
$env:Path += ";$(go env GOPATH)\bin"
```

### Opción 2: Instalación desde el código fuente

```bash
# Clonar el repositorio
git clone https://github.com/YeridStick/cleango.git
cd cleango

# Instalar dependencias
go mod download

# Instalar el CLI
go install ./cmd/cleango
```

### Verificar instalación

```bash
cleango --version
```

---

## 🚀 Uso Rápido

### Crear un nuevo proyecto

#### Modo interactivo (recomendado)

```bash
cleango new
```

El CLI te guiará paso a paso preguntando:
- Nombre del proyecto
- Módulo Go
- Framework HTTP
- Base de datos
- Extras (Redis, Kafka)

#### Modo no interactivo

```bash
cleango new my-service \
  --module github.com/user/my-service \
  --framework chi \
  --database postgres \
  --redis \
  --kafka
```

Flags disponibles:
- `-m, --module`: Ruta del módulo Go
- `-f, --framework`: Framework HTTP (`nethttp`, `chi`, `gin`, `fiber`)
- `-d, --database`: Base de datos (`none`, `postgres`, `mysql`, `mongodb`, `oracle`)
- `--redis`: Incluir Redis
- `--kafka`: Incluir Kafka
- `--non-interactive`: Modo no interactivo (usa valores por defecto)

---

## 🔨 Generación de Componentes

Después de crear tu proyecto, puedes agregar componentes fácilmente:

### Crear un caso de uso

```bash
cd my-service
cleango add usecase GetUser
```

Esto genera: `internal/usecase/get_user.go` con:
- Interface del caso de uso
- Implementación concreta
- Structs de Input/Output

### Crear un adaptador/repositorio

```bash
cleango add adapter UserRepository
```

Genera: `internal/repository/user_repository.go` con:
- Interface del repositorio
- Implementación con métodos CRUD
- Métodos: FindByID, Save, Update, Delete

### Crear un modelo de dominio

```bash
cleango add model User
```

Genera: `internal/domain/user.go` con:
- Estructura del modelo
- Campos base (ID, CreatedAt, UpdatedAt)
- Métodos de validación

### Crear un handler HTTP

```bash
cleango add handler User
```

Genera: `internal/http/user_handler.go` con:
- Estructura del handler
- Métodos HTTP (Get, Post, Put, Delete)
- Manejo básico de requests/responses

---

## 📁 Estructura del Proyecto Generado

```
my-service/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go            # Configuración centralizada
│   ├── logger/
│   │   └── logger.go            # Logger estructurado (zap)
│   ├── http/
│   │   └── *_handler.go         # Handlers HTTP
│   ├── domain/
│   │   └── *.go                 # Modelos de dominio
│   ├── usecase/
│   │   └── *.go                 # Casos de uso (lógica de negocio)
│   ├── repository/
│   │   └── *.go                 # Adaptadores/Repositorios
│   └── db/
│       └── ...                  # Conexiones a DB
├── migrations/                  # Migraciones de base de datos
├── go.mod
├── go.sum
└── .gitignore
```

Esta estructura sigue los principios de **Clean Architecture**:
- 🎯 **Domain**: Entidades de negocio (independientes)
- 💼 **Usecase**: Lógica de negocio (independiente de frameworks)
- 🔌 **Adapters**: Conexión con el mundo externo (DB, HTTP, etc.)
- 🌐 **HTTP**: Capa de presentación

---

## 🎨 Frameworks Soportados

### HTTP Frameworks

| Framework | Uso | Descripción |
|-----------|-----|-------------|
| `nethttp` | Estándar | net/http nativo de Go |
| `chi` | Ligero | Router minimalista y rápido |
| `gin` | Popular | Framework completo y performante |
| `fiber` | Express-like | Inspirado en Express.js |

### Bases de Datos

| Base de Datos | Driver |
|---------------|--------|
| `postgres` | `github.com/jackc/pgx/v5` |
| `mysql` | `github.com/go-sql-driver/mysql` |
| `mongodb` | `go.mongodb.org/mongo-driver` |
| `oracle` | `github.com/godror/godror` |

### Extras

- **Redis**: `github.com/redis/go-redis/v9`
- **Kafka**: `github.com/segmentio/kafka-go`

---

## 📖 Ejemplos Completos

### Ejemplo 1: API REST con Chi y Postgres

```bash
# Crear proyecto
cleango new user-api \
  -m github.com/myorg/user-api \
  -f chi \
  -d postgres

cd user-api

# Agregar componentes
cleango add model User
cleango add adapter UserRepository
cleango add usecase CreateUser
cleango add usecase GetUser
cleango add handler User

# Ejecutar
go mod tidy
APP_PORT=8080 go run ./cmd/api
```

### Ejemplo 2: Microservicio con Gin, Redis y Kafka

```bash
cleango new notification-service \
  -m github.com/myorg/notification-service \
  -f gin \
  -d mongodb \
  --redis \
  --kafka \
  --non-interactive

cd notification-service
go mod tidy
go run ./cmd/api
```

---

## ⚙️ Configuración

El proyecto generado usa variables de entorno para configuración:

```bash
# Configuración general
APP_ENV=dev          # Entorno: dev, prod
APP_PORT=8080        # Puerto HTTP

# Base de datos
DB_POSTGRES_URL=postgresql://user:pass@localhost:5432/dbname
DB_MYSQL_DSN=user:pass@tcp(localhost:3306)/dbname
DB_MONGO_URI=mongodb://localhost:27017
DB_ORACLE_DSN=user/pass@localhost:1521/ORCL

# Extras
REDIS_ADDR=localhost:6379
KAFKA_BROKERS=localhost:9092
```

Puedes crear un archivo `.env` (no incluido en git) o exportar las variables:

```bash
export APP_ENV=dev
export APP_PORT=8080
export DB_POSTGRES_URL="postgresql://..."

go run ./cmd/api
```

---

## 🔍 Comandos Disponibles

```bash
# Ver ayuda general
cleango --help

# Ver ayuda de un comando específico
cleango new --help
cleango add --help

# Crear nuevo proyecto
cleango new [nombre] [flags]

# Agregar componentes
cleango add usecase [nombre]
cleango add adapter [nombre]
cleango add model [nombre]
cleango add handler [nombre]

# Ver versión
cleango --version
```

---

## 🛠️ Desarrollo del CLI

Si quieres contribuir o modificar el CLI:

```bash
# Clonar repositorio
git clone https://github.com/YeridStick/cleango.git
cd cleango

# Instalar dependencias
go mod download

# Ejecutar sin instalar
go run ./cmd/cleango new test-project

# Compilar
go build -o cleango ./cmd/cleango

# Ejecutar tests (cuando estén implementados)
go test ./...
```

---

## 📝 Roadmap

- [ ] Tests unitarios completos
- [ ] Comando `cleango migrate` para migraciones
- [ ] Templates personalizables
- [ ] Soporte para gRPC
- [ ] Generación de Dockerfiles
- [ ] Generación de CI/CD configs
- [ ] Comando `cleango deploy`

---

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

---

## 📄 Licencia

Este proyecto está bajo la licencia MIT. Ver archivo `LICENSE` para más detalles.

---

## 🙏 Agradecimientos

Inspirado por herramientas como:
- [Spring Initializr](https://start.spring.io/)
- [Rails Generators](https://guides.rubyonrails.org/generators.html)
- [Scaffold Clean Architecture (Bancolombia)](https://github.com/bancolombia/scaffold-clean-architecture)

---

## 📧 Contacto

**Yerid Stick**
- GitHub: [@YeridStick](https://github.com/YeridStick)

---

## ⭐ Si te gusta el proyecto, ¡dale una estrella!

```bash
# Instalación rápida
go install github.com/YeridStick/cleango/cmd/cleango@latest

# Primer proyecto
cleango new my-awesome-api

# ¡A programar! 🚀
```
