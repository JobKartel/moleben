package domain

import "time"

// SessionStatus — статус сессии.
type SessionStatus string

const (
	// SessionActive — сессия активна, можно отправлять сообщения.
	SessionActive SessionStatus = "active"
	// SessionClosed — сессия закрыта (достигнут итог).
	SessionClosed SessionStatus = "closed"
)

// Session — доменная модель сессии чата.
type Session struct {
	ID           string        // ID сессии (UUID)
	CreatedAt    time.Time     // Время создания
	UpdatedAt    time.Time     // Время последнего обновления
	Status       SessionStatus // Статус сессии: active/closed
	Tally        int           // Текущий «счётчик покаяния» в диапазоне [0..100]
	FinalMessage *string       // Итоговое сообщение («вы прощены» / «вы не прощены»), если закрыта
}

// MessageRole — роль автора сообщения.
type MessageRole string

const (
	// RoleUser — сообщение от пользователя.
	RoleUser MessageRole = "user"
	// RoleAssistant — сообщение от ассистента (LLM).
	RoleAssistant MessageRole = "assistant"
)

// Message — доменная модель сообщения в рамках сессии.
type Message struct {
	ID        string      // ID сообщения (UUID)
	SessionID string      // Связь с сессией
	Role      MessageRole // Роль автора
	Content   string      // Содержимое (для ассистента — сырой JSON-пейлоад)
	CreatedAt time.Time   // Время создания
}
