package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/msayib/todo-fiber-dig/internal/model"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Save(ctx context.Context, todo model.Todo) (model.Todo, error)
	FindByID(ctx context.Context, id uint) (model.Todo, error)
	FindAll(ctx context.Context) ([]model.Todo, error)
	Update(ctx context.Context, todo model.Todo) (model.Todo, error)
	Delete(ctx context.Context, id uint) error
}

type TodoRepositoryGORM struct {
	db *gorm.DB
}

func NewTodoRepositoryGORM(db *gorm.DB) *TodoRepositoryGORM {
	return &TodoRepositoryGORM{db}
}

func (r *TodoRepositoryGORM) Save(ctx context.Context, todo model.Todo) (model.Todo, error) {
	result := r.db.WithContext(ctx).Create(&todo)
	return todo, result.Error
}

func (r *TodoRepositoryGORM) FindByID(ctx context.Context, id uint) (model.Todo, error) {
	var todo model.Todo
	result := r.db.WithContext(ctx).First(&todo, id)
	return todo, result.Error
}

func (r *TodoRepositoryGORM) FindAll(ctx context.Context) ([]model.Todo, error) {
	var todos []model.Todo
	result := r.db.WithContext(ctx).Order("created_at desc").Find(&todos)
	return todos, result.Error
}

func (r *TodoRepositoryGORM) Update(ctx context.Context, todo model.Todo) (model.Todo, error) {
	result := r.db.WithContext(ctx).Save(&todo)
	return todo, result.Error
}

func (r *TodoRepositoryGORM) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Todo{}, id)
	return result.Error
}


// --- Implementasi Decorator Cache (Redis) ---
type todoRepositoryCache struct {
	repo   TodoRepository
	redis  *redis.Client
	expiry time.Duration
}

func NewTodoRepositoryCache(repo *TodoRepositoryGORM, redis *redis.Client) TodoRepository {
	return &todoRepositoryCache{
		repo:   repo,
		redis:  redis,
		expiry: time.Minute * 10,
	}
}

func (r *todoRepositoryCache) getCacheKey(id uint) string {
	return fmt.Sprintf("todo:%d", id)
}

func (r *todoRepositoryCache) FindByID(ctx context.Context, id uint) (model.Todo, error) {
	key := r.getCacheKey(id)
	result, err := r.redis.Get(ctx, key).Result()
	if err == nil {
		var todo model.Todo
		if err := json.Unmarshal([]byte(result), &todo); err == nil {
			return todo, nil
		}
	}

	todo, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return model.Todo{}, err
	}

	data, err := json.Marshal(todo)
	if err == nil {
		r.redis.Set(ctx, key, data, r.expiry)
	}
	return todo, nil
}

func (r *todoRepositoryCache) Update(ctx context.Context, todo model.Todo) (model.Todo, error) {
	updatedTodo, err := r.repo.Update(ctx, todo)
	if err != nil {
		return model.Todo{}, err
	}
	key := r.getCacheKey(updatedTodo.ID)
	r.redis.Del(ctx, key)
	return updatedTodo, nil
}

func (r *todoRepositoryCache) Delete(ctx context.Context, id uint) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}
	key := r.getCacheKey(id)
	r.redis.Del(ctx, key)
	return nil
}

func (r *todoRepositoryCache) Save(ctx context.Context, todo model.Todo) (model.Todo, error) {
	return r.repo.Save(ctx, todo)
}

func (r *todoRepositoryCache) FindAll(ctx context.Context) ([]model.Todo, error) {
	return r.repo.FindAll(ctx)
}