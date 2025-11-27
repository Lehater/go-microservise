// models/user.go
package models

import "errors"

// Простейшая доменная модель пользователя.
// Для ДЗ достаточно базовой валидации.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate — простая проверка обязательных полей.
// В реальном проекте валидатор был бы богаче.
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	// Можно добавить простую проверку на '@', но не усложняем.
	return nil
}
