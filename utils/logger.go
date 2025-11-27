// utils/logger.go
package utils

import (
	"log"
	"os"

	"go-microservice/models"
)

// AuditEvent — событие для audit-лога.
type AuditEvent struct {
	Action string
	User   models.User
}

// AsyncAuditLogger — один канал + одна goroutine-писатель.
type AsyncAuditLogger struct {
	ch chan AuditEvent
}

// NewAsyncAuditLogger создаёт логер и запускает worker.
// buffer — размер буфера канала (для пиков нагрузки).
func NewAsyncAuditLogger(buffer int) *AsyncAuditLogger {
	l := &AsyncAuditLogger{
		ch: make(chan AuditEvent, buffer),
	}

	go func() {
		logger := log.New(os.Stdout, "[AUDIT] ", log.LstdFlags)
		for e := range l.ch {
			logger.Printf("%s user_id=%d email=%s\n", e.Action, e.User.ID, e.User.Email)
		}
	}()

	return l
}

// LogAction — неблокирующая отправка в канал.
func (l *AsyncAuditLogger) LogAction(action string, user models.User) {
	select {
	case l.ch <- AuditEvent{Action: action, User: user}:
	default:
		// Буфер заполнен — дропаем событие.
	}
}

// StubNotificationSender — заглушка для уведомлений.
type StubNotificationSender struct {
	logger *log.Logger
}

func NewStubNotificationSender() *StubNotificationSender {
	return &StubNotificationSender{
		logger: log.New(os.Stdout, "[NOTIFY] ", log.LstdFlags),
	}
}

func (s *StubNotificationSender) SendUserNotification(action string, user models.User) {
	// Можно вынести в отдельную goroutine, но на уровне UserService
	// мы уже работаем асинхронно, поэтому здесь достаточно sync-логирования.
	s.logger.Printf("%s user_id=%d email=%s\n", action, user.ID, user.Email)
}
