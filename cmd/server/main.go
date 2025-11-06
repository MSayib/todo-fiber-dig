package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/msayib/todo-fiber-dig/internal/config"
	"github.com/msayib/todo-fiber-dig/internal/handler"
	"github.com/msayib/todo-fiber-dig/internal/repository"
	"github.com/msayib/todo-fiber-dig/internal/service"
	"github.com/msayib/todo-fiber-dig/pkg/cache"
	"github.com/msayib/todo-fiber-dig/pkg/database"
	"github.com/redis/go-redis/v9"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func runMigrations(cfg config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	m, err := migrate.New("file://internal/migrations", dsn)
	if err != nil {
		log.Fatalf("Gagal membuat instance migrasi: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Gagal menjalankan migrasi 'up': %v", err)
	}

	log.Println("Migrasi database berhasil dijalankan.")
}

func buildContainer() (*dig.Container, error) {
	container := dig.New()

	if err := container.Provide(config.LoadConfig); err != nil {
		return nil, err
	}

	container.Invoke(func(cfg config.Config) {
		runMigrations(cfg)
	})

	if err := container.Provide(database.NewGormDB); err != nil {
		return nil, err
	}
	if err := container.Provide(cache.NewRedisClient); err != nil {
		return nil, err
	}

	// Menyediakan implementasi GORM konkret
	container.Provide(repository.NewTodoRepositoryGORM)

	// Menyediakan implementasi Cache (Decorator) yang membungkus GORM
	// Diubah agar menggunakan tipe *repository.TodoRepositoryGORM yang sudah exported
	container.Provide(func(gormRepo *repository.TodoRepositoryGORM, redisClient *redis.Client) repository.TodoRepository {
		return repository.NewTodoRepositoryCache(gormRepo, redisClient)
	})

	container.Provide(service.NewTodoService)
	container.Provide(handler.NewTodoHandler)

	return container, nil
}

func runServer(container *dig.Container) error {
	return container.Invoke(func(cfg config.Config, db *gorm.DB, rdb *redis.Client, todoHandler *handler.TodoHandler) error {
		app := fiber.New()
		app.Use(logger.New())

		// Setup Rute
		api := app.Group("/api/v1")
		api.Get("/health", todoHandler.HealthCheck)

		todos := api.Group("/todos")
		todos.Post("/", todoHandler.CreateTodo)
		todos.Get("/", todoHandler.GetAllTodos)
		todos.Get("/:id", todoHandler.GetTodoByID)
		todos.Put("/:id", todoHandler.UpdateTodo)
		todos.Delete("/:id", todoHandler.DeleteTodo)

		log.Printf("Server berjalan di port %s", cfg.AppPort)
		return app.Listen(cfg.AppPort)
	})
}

func main() {
	container, err := buildContainer()
	if err != nil {
		log.Fatalf("Gagal membangun container DI: %v", err)
	}
	if err := runServer(container); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
