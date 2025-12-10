package generator

// ProjectConfig holds the configuration for a new project
type ProjectConfig struct {
	Name       string
	ModulePath string
	Framework  string
	Database   string
	UseRedis   bool
	UseKafka   bool
}

// GetDependencies returns the list of Go dependencies to install
func (c *ProjectConfig) GetDependencies() []string {
	deps := []string{
		"go.uber.org/zap", // logger
	}

	// Add framework dependencies
	switch c.Framework {
	case "chi":
		deps = append(deps, "github.com/go-chi/chi/v5")
	case "gin":
		deps = append(deps, "github.com/gin-gonic/gin")
	case "fiber":
		deps = append(deps, "github.com/gofiber/fiber/v2")
	}

	// Add database dependencies
	switch c.Database {
	case "postgres":
		deps = append(deps, "github.com/jackc/pgx/v5/stdlib")
	case "mysql":
		deps = append(deps, "github.com/go-sql-driver/mysql")
	case "mongodb":
		deps = append(deps, "go.mongodb.org/mongo-driver/mongo")
	case "oracle":
		deps = append(deps, "github.com/godror/godror")
	}

	// Add optional dependencies
	if c.UseRedis {
		deps = append(deps, "github.com/redis/go-redis/v9")
	}

	if c.UseKafka {
		deps = append(deps, "github.com/segmentio/kafka-go")
	}

	return deps
}
