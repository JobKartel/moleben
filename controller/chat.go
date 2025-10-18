package controller

import (
	"moleben/entity"
	"moleben/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatController struct{ svc service.ChatService }

func NewChatController(svc service.ChatService) *ChatController { return &ChatController{svc: svc} }

// CreateSession godoc
// @Summary Создать новую чат-сессию
// @Description Запускает новую «исповедальную» сессию с начальным значением счётчика (tally) 50.
// @Tags sessions
// @Accept json
// @Produce json
// @Success 201 {object} entity.CreateSessionResponse
// @Failure 500 {object} map[string]string
// @Router /api/sessions [post]
func (c *ChatController) CreateSession(ctx *gin.Context) {
	sess, err := c.svc.CreateSession(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, entity.CreateSessionResponse{ID: sess.ID, Tally: sess.Tally, Status: string(sess.Status)})
}

// PostMessage godoc
// @Summary Отправить сообщение пользователя в сессию
// @Description Принимает текст пользователя, отправляет его в LLM, возвращает JSON-ответ ассистента и обновляет счётчик (tally).
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "ID сессии"
// @Param body body entity.PostMessageRequest true "Сообщение пользователя"
// @Success 200 {object} entity.PostMessageResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/sessions/{id}/messages [post]
func (c *ChatController) PostMessage(ctx *gin.Context) {
	id := ctx.Param("id")
	var req entity.PostMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Text == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
		return
	}
	resp, err := c.svc.PostMessage(ctx.Request.Context(), id, req.Text)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "session not found" {
			status = http.StatusNotFound
		}
		if err.Error() == "session already closed" {
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}
	// Заголовки для кэширования и даты ответа.
	ctx.Header("Cache-Control", "no-store")
	ctx.Header("Date", time.Now().UTC().Format(http.TimeFormat))
	ctx.JSON(http.StatusOK, resp)
}

// GetSession godoc
// @Summary Получить состояние сессии
// @Description Возвращает текущий статус, значение счётчика и финальный результат (если сессия завершена).
// @Tags sessions
// @Produce json
// @Param id path string true "ID сессии"
// @Success 200 {object} entity.GetSessionResponse
// @Failure 404 {object} map[string]string
// @Router /api/sessions/{id} [get]
func (c *ChatController) GetSession(ctx *gin.Context) {
	id := ctx.Param("id")
	sess, err := c.svc.GetSession(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, entity.GetSessionResponse{ID: sess.ID, Status: string(sess.Status), Tally: sess.Tally, Outcome: sess.FinalMessage})
}

// ListMessages godoc
// @Summary Получить сообщения сессии
// @Description Возвращает хронологическую историю сообщений (первые 200).
// @Tags messages
// @Produce json
// @Param id path string true "ID сессии"
// @Success 200 {object} entity.ListMessagesResponse
// @Failure 500 {object} map[string]string
// @Router /api/sessions/{id}/messages [get]
func (c *ChatController) ListMessages(ctx *gin.Context) {
	id := ctx.Param("id")
	items, err := c.svc.ListMessages(ctx.Request.Context(), id, 200, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]entity.MessageDTO, 0, len(items))
	for _, m := range items {
		out = append(out, entity.MessageDTO{ID: m.ID, Role: string(m.Role), Content: m.Content, CreatedAt: m.CreatedAt.Format(time.RFC3339)})
	}
	ctx.JSON(http.StatusOK, entity.ListMessagesResponse{SessionID: id, Items: out})
}
