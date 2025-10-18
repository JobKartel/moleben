package service

import (
	"context"
	"encoding/json"
	"errors"
	"moleben/domain"
	"moleben/entity"
	"moleben/repository"
	"moleben/repository/jsonutil"
	"time"

	"github.com/google/uuid"
)

type ChatService interface {
	CreateSession(ctx context.Context) (*domain.Session, error)
	PostMessage(ctx context.Context, sessionID, userText string) (*entity.PostMessageResponse, error)
	GetSession(ctx context.Context, sessionID string) (*domain.Session, error)
	ListMessages(ctx context.Context, sessionID string, limit, offset int) ([]domain.Message, error)
}

type chatService struct {
	repo         repository.SessionRepository
	llm          *repository.Client
	systemPrompt string
}

func NewChatService(repo repository.SessionRepository, client *repository.Client, systemPrompt string) ChatService {
	return &chatService{repo: repo, llm: client, systemPrompt: systemPrompt}
}

func (s *chatService) CreateSession(ctx context.Context) (*domain.Session, error) {
	now := time.Now()
	sess := &domain.Session{ID: uuid.New().String(), CreatedAt: now, UpdatedAt: now, Status: domain.SessionActive, Tally: 50}
	if err := s.repo.Create(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *chatService) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	return s.repo.Get(ctx, sessionID)
}

var ErrSessionClosed = errors.New("session already closed")

func (s *chatService) PostMessage(ctx context.Context, sessionID, userText string) (*entity.PostMessageResponse, error) {
	sess, err := s.repo.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if sess.Status == domain.SessionClosed {
		return nil, ErrSessionClosed
	}

	um := &domain.Message{ID: uuid.New().String(), SessionID: sessionID, Role: domain.RoleUser, Content: userText, CreatedAt: time.Now()}
	if err = s.repo.InsertMessage(ctx, um); err != nil {

		return nil, err
	}

	assistantRaw, err := s.llm.Chat(ctx, s.systemPrompt, userText)
	if err != nil {
		return nil, err
	}

	jsonStr, ok := jsonutil.ExtractJSONObject(assistantRaw)
	if !ok {
		return nil, errors.New("assistant did not return JSON block")
	}

	var payload entity.BotPayload
	if err = json.Unmarshal([]byte(jsonStr), &payload); err != nil {
		return nil, err
	}

	am := &domain.Message{ID: uuid.New().String(), SessionID: sessionID, Role: domain.RoleAssistant, Content: jsonStr, CreatedAt: time.Now()}
	if err = s.repo.InsertMessage(ctx, am); err != nil {
		return nil, err
	}
	if !payload.RiskFlag && payload.Repentance.Score != nil {
		delta := *payload.Repentance.Score - 50
		sess.Tally += delta
		if sess.Tally < 0 {
			sess.Tally = 0
		}
		if sess.Tally > 100 {
			sess.Tally = 100
		}
	}

	resp := &entity.PostMessageResponse{Assistant: payload, Tally: sess.Tally}
	if sess.Tally == 0 || sess.Tally == 100 {
		sess.Status = domain.SessionClosed
		out := "вы не прощены"
		if sess.Tally == 100 {
			out = "вы прощены"
		}
		sess.FinalMessage = &out
		resp.Done, resp.Outcome = true, &out
	}
	if err = s.repo.Update(ctx, sess); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *chatService) ListMessages(ctx context.Context, sessionID string, limit, offset int) ([]domain.Message, error) {
	return s.repo.ListMessages(ctx, sessionID, limit, offset)
}
