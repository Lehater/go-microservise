// storage/user_repo_memory.go
package storage

import (
	"errors"
	"sync"

	"go-microservice/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id int) (models.User, error)
	Update(id int, user models.User) (models.User, error)
	Delete(id int) error
}

// Потокобезопасный in-memory репозиторий.
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	data   map[int]models.User
	nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		data:   make(map[int]models.User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepository) Create(user models.User) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.data[user.ID] = user

	return user, nil
}

func (r *InMemoryUserRepository) GetAll() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]models.User, 0, len(r.data))
	for _, u := range r.data {
		users = append(users, u)
	}
	return users, nil
}

func (r *InMemoryUserRepository) GetByID(id int) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.data[id]
	if !ok {
		return models.User{}, ErrUserNotFound
	}
	return u, nil
}

func (r *InMemoryUserRepository) Update(id int, user models.User) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return models.User{}, ErrUserNotFound
	}
	user.ID = id
	r.data[id] = user
	return user, nil
}

func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrUserNotFound
	}
	delete(r.data, id)
	return nil
}
