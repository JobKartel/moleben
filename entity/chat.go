package entity

// Signal — единичный сигнал/фактор, который повлиял на оценку «покаяния».
type Signal struct {
	Factor    string `json:"factor"`    // Название фактора
	Direction string `json:"direction"` // Направление влияния (например, up/down/neutral)
	Evidence  string `json:"evidence"`  // Краткое обоснование/доказательство
}

// Repentance — агрегированная оценка «покаяния» от ассистента.
type Repentance struct {
	Score      *int     `json:"score"`      // Итоговый балл (0..100), может отсутствовать при RiskFlag
	Trend      string   `json:"trend"`      // Тренд изменения (например, improving/declining/stable)
	Signals    []Signal `json:"signals"`    // Перечень факторов
	Confidence string   `json:"confidence"` // Уверенность модели (низкая/средняя/высокая и т.п.)
}

// BotPayload — структура JSON-ответа ассистента, сохраняемая в БД и отдаваемая клиенту.
type BotPayload struct {
	Reply      string     `json:"reply"`      // Текстовый ответ ассистента для пользователя
	Repentance Repentance `json:"repentance"` // Детализация оценки «покаяния»
	NextStep   string     `json:"next_step"`  // Рекомендованное следующее действие
	RiskFlag   bool       `json:"risk_flag"`  // Флаг риска: при true tally не меняется
}

// CreateSessionResponse — ответ API при создании сессии.
type CreateSessionResponse struct {
	ID     string `json:"id"`     // ID созданной сессии
	Tally  int    `json:"tally"`  // Начальное значение счётчика
	Status string `json:"status"` // Статус сессии
}

// PostMessageRequest — запрос на отправку пользовательского сообщения.
type PostMessageRequest struct {
	Text string `json:"text"` // Текст пользователя
}

// PostMessageResponse — ответ API на отправку сообщения.
type PostMessageResponse struct {
	Assistant BotPayload `json:"assistant"` // Пейлоад ассистента (JSON)
	Tally     int        `json:"tally"`     // Текущее значение счётчика после обработки
	Done      bool       `json:"done"`      // Признак завершения сессии
	Outcome   *string    `json:"outcome"`   // Итог («вы прощены»/«вы не прощены»), если Done=true
}

// GetSessionResponse — ответ API на запрос состояния сессии.
type GetSessionResponse struct {
	ID      string  `json:"id"`      // ID сессии
	Status  string  `json:"status"`  // Текущий статус
	Tally   int     `json:"tally"`   // Текущее значение счётчика
	Outcome *string `json:"outcome"` // Итог (если сессия закрыта)
}

// MessageDTO — DTO сообщения для выдачи наружу.
type MessageDTO struct {
	ID        string `json:"id"`         // ID сообщения
	Role      string `json:"role"`       // Роль автора (user/assistant)
	Content   string `json:"content"`    // Содержимое (для ассистента — JSON-строка)
	CreatedAt string `json:"created_at"` // Время создания в формате RFC3339
}

// ListMessagesResponse — список сообщений сессии.
type ListMessagesResponse struct {
	SessionID string       `json:"session_id"` // ID сессии
	Items     []MessageDTO `json:"items"`      // Сообщения в хронологическом порядке
}
