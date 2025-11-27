// services/user_service.go
package services

import (
	"go-microservice/models"
	"go-microservice/storage"
)

// Интерфейс для audit-логера (реализуется в utils.AsyncAuditLogger).
type AuditLogger interface {
	LogAction(action string, user models.User)
}

// Интерфейс для уведомлений (реализация — StubNotificationSender).
type NotificationSender interface {
	SendUserNotification(action string, user models.User)
}

// UserService инкапсулирует работу с репозиторием и инфраструктурой.
type UserService struct {
	repo   storage.UserRepository
	audit  AuditLogger
	notify NotificationSender
}

func NewUserService(
	repo storage.UserRepository,
	audit AuditLogger,
	notify NotificationSender,
) *UserService {
	return &UserService{
		repo:   repo,
		audit:  audit,
		notify: notify,
	}
}

func (s *UserService) Create(user models.User) (models.User, error) {
	if err := user.Validate(); err != nil {
		return models.User{}, err
	}

	created, err := s.repo.Create(user)
	if err != nil {
		return models.User{}, err
	}

	// Асинхронные действия запускаем в отдельной goroutine.
	go s.audit.LogAction("CREATE", created)
	go s.notify.SendUserNotification("CREATE", created)

	return created, nil
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetByID(id int) (models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(id int, user models.User) (models.User, error) {
	if err := user.Validate(); err != nil {
		return models.User{}, err
	}

	updated, err := s.repo.Update(id, user)
	if err != nil {
		return models.User{}, err
	}

	go s.audit.LogAction("UPDATE", updated)
	go s.notify.SendUserNotification("UPDATE", updated)

	return updated, nil
}

func (s *UserService) Delete(id int) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	go s.audit.LogAction("DELETE", user)
	go s.notify.SendUserNotification("DELETE", user)

	return nil
}
