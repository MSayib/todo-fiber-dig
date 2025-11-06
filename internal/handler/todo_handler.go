package handler

import (
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/msayib/todo-fiber-dig/internal/model"
	"github.com/msayib/todo-fiber-dig/internal/service"
	"gorm.io/gorm"
)

type TodoHandler struct {
	service  service.TodoService
	validate *validator.Validate
	db       *gorm.DB
	redis    *redis.Client
}

func NewTodoHandler(s service.TodoService, db *gorm.DB, redis *redis.Client) *TodoHandler {
	return &TodoHandler{
		service:  s,
		validate: validator.New(),
		db:       db,
		redis:    redis,
	}
}

func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var req model.CreateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return sendJSONResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.validate.Struct(req); err != nil {
		return sendJSONResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	todo, err := h.service.Create(c.Context(), req)
	if err != nil {
		return sendJSONResponse(c, fiber.StatusInternalServerError, "Failed to create todo", nil)
	}

	return sendJSONResponse(c, fiber.StatusCreated, "Todo created successfully", todo)
}

func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	// Diubah dari Atoi ke ParseUint
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return sendJSONResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	todo, err := h.service.GetByID(c.Context(), uint(id)) // Casting ke uint
	if err != nil {
		log.Printf("Error getting todo by id %d: %v", id, err)
		return sendJSONResponse(c, fiber.StatusNotFound, "Todo not found", nil)
	}

	return sendJSONResponse(c, fiber.StatusOK, "Todo found", todo)
}

func (h *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
	todos, err := h.service.GetAll(c.Context())
	if err != nil {
		return sendJSONResponse(c, fiber.StatusInternalServerError, "Failed to retrieve todos", nil)
	}
	return sendJSONResponse(c, fiber.StatusOK, "Todos retrieved successfully", todos)
}

func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	// Diubah dari Atoi ke ParseUint
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return sendJSONResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var req model.UpdateTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return sendJSONResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	todo, err := h.service.Update(c.Context(), uint(id), req) // Casting ke uint
	if err != nil {
		return sendJSONResponse(c, fiber.StatusInternalServerError, "Failed to update todo", nil)
	}
	return sendJSONResponse(c, fiber.StatusOK, "Todo updated successfully", todo)
}

func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	// Diubah dari Atoi ke ParseUint
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return sendJSONResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	if err := h.service.Delete(c.Context(), uint(id)); err != nil { // Casting ke uint
		return sendJSONResponse(c, fiber.StatusInternalServerError, "Failed to delete todo", nil)
	}

	return sendJSONResponse(c, fiber.StatusOK, "Todo deleted successfully", nil)
}

func (h *TodoHandler) HealthCheck(c *fiber.Ctx) error {
	dbStatus := "connected"
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "disconnected"
	}

	redisStatus := "connected"
	if err := h.redis.Ping(c.Context()).Err(); err != nil {
		redisStatus = "disconnected"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"app_status": "running",
		"timestamp":  time.Now(),
		"dependencies": fiber.Map{
			"postgres": dbStatus,
			"redis":    redisStatus,
		},
	})
}

// Fungsi helper untuk response JSON
func sendJSONResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  status,
		"message": message,
		"data":    data,
	})
}